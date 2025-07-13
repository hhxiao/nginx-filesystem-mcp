package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log"
	"strings"
)

func listFolder(client *Client, ctx context.Context, path string, recursive bool) ([]DirectoryEntry, error) {
	content, err := client.GetContent(ctx, path)
	if err != nil {
		return nil, err
	}
	entries, err := parseDirectoryListing(content.Data)
	if err != nil {
		return nil, err
	}
	if entries == nil {
		return nil, errors.New("404 Not Found")
	}
	if recursive {
		for i := range entries {
			if entries[i].Type == Directory {
				subPath := path + entries[i].Name + "/"
				subEntries, err := listFolder(client, ctx, subPath, recursive)
				if err != nil {
					return nil, err
				}
				entries[i].Entries = subEntries
			}
		}
	}
	return entries, nil
}

func readFile(client *Client, ctx context.Context, path string) (*Content, error) {
	content, err := client.GetContent(ctx, path)
	if err != nil {
		return nil, err
	}
	if content.Length == "" {
		return nil, errors.New("404 Not Found")
	}
	return content, nil
}

func toJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func handle(client *Client, ctx context.Context, path string, recursive bool) (*mcp.CallToolResult, error) {
	folder := path == "" || strings.HasSuffix(path, "/")
	if folder {
		entries, err := listFolder(client, ctx, path, recursive)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(
			toJson(entries),
		), nil
	}
	content, err := readFile(client, ctx, path)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(
		content.Data,
	), nil
}

func RegisterTools(server *server.MCPServer, client *Client) {
	server.AddTool(mcp.NewTool("listFolder",
		mcp.WithDescription("list a folder"),
		mcp.WithString("path",
			mcp.Description("path of the directory to list"),
		),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "list a folder",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("%s - %v\n", request.Params.Name, toJson(request.Params.Arguments))
		path := request.GetString("path", "")
		path = strings.TrimSuffix(path, "/") + "/"
		return handle(client, ctx, path, false)
	})

	server.AddTool(mcp.NewTool("treeFolder",
		mcp.WithDescription("Returns a hierarchical JSON representation of directory structure"),
		mcp.WithString("path",
			mcp.Description("path of the directory to traverse"),
		),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "traverse a folder",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("%s - %v\n", request.Params.Name, toJson(request.Params.Arguments))
		path := request.GetString("path", "")
		path = strings.TrimSuffix(path, "/") + "/"
		return handle(client, ctx, path, true)
	})

	server.AddTool(mcp.NewTool("readPFile",
		mcp.WithDescription("read a file"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("path to the file to read"),
		),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "read a file",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("%s - %v\n", request.Params.Name, toJson(request.Params.Arguments))
		path := request.GetString("path", "")
		path = strings.TrimSuffix(path, "/")
		return handle(client, ctx, path, false)
	})
}

func RegisterResources(server *server.MCPServer, client *Client) {
	server.AddResourceTemplate(mcp.NewResourceTemplate(
		"files://{path}",
		"File",
		mcp.WithTemplateDescription("Return path information for a folder or a file"),
		mcp.WithTemplateMIMEType("application/json"),
	), func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		path := strings.TrimLeft(request.Params.URI, "files://")
		folder := path == "" || strings.HasSuffix(path, "/")
		if folder {
			entries, err := listFolder(client, ctx, path, false)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      request.Params.URI,
					MIMEType: "application/json",
					Text:     toJson(entries),
				},
			}, nil
		} else {
			content, err := readFile(client, ctx, path)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      request.Params.URI,
					MIMEType: content.Type,
					Text:     content.Data,
				},
			}, nil
		}
	})
}
