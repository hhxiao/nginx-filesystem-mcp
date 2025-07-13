package main

import (
	"flag"
	"fmt"
	"github.com/hhxiao/nginx-filesystem-mcp/pkg"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"log"
)

var version = "0.0.0"

func main() {
	_ = godotenv.Load()

	transport := "stream"
	port := "9292"

	flag.StringVar(&transport, "t", transport, "Transport type (stdio|sse|stream)")
	flag.StringVar(&port, "port", port, "Port for SSE/Streamable HTTP transport type")
	flag.Parse()

	// Create a new MCP handler
	s := server.NewMCPServer(
		"nginx filesystem mcp",
		version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	client := pkg.NewClient()
	pkg.RegisterTools(s, client)
	pkg.RegisterResources(s, client)

	// Start the handler
	switch transport {
	case "sse":
		sseServer := server.NewSSEServer(s,
			server.WithBaseURL(fmt.Sprintf("http://0.0.0.0:%s", port)),
		)
		sseEndpoint, _ := sseServer.CompleteSseEndpoint()
		log.Printf("nginx filesystem mcp server(%s) listening on %s 🚀", version, sseEndpoint)
		if err := sseServer.Start(":" + port); err != nil {
			log.Fatalf("nginx filesystem mcp server error: %v", err)
		}
	case "stream":
		streamServer := server.NewStreamableHTTPServer(s,
			server.WithStateLess(true),
		)
		streamEndpoint := fmt.Sprintf("http://0.0.0.0:%s/mcp", port)
		log.Printf("nginx filesystem mcp server(%s) listening on %s 🚀", version, streamEndpoint)
		if err := streamServer.Start(":" + port); err != nil {
			log.Fatalf("nginx filesystem mcp server error: %v", err)
		}
	case "stdio":
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("nginx filesystem mcp server error: %v\n", err)
		}
	default:
		log.Fatalf("unknown transport type: %s", transport)
	}
}
