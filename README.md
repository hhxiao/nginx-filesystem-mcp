# nginx-filesystem-mcp

This MCP server enables secure access to files served via nginx through the Model Context Protocol (MCP).

## Configuration Guide

### LangGraph Python Client

To configure the nginx-filesystem-mcp in a LangGraph-based Python client:

```python
import os
from rest_mcp import MultiServerMCPClient

client = MultiServerMCPClient(
    {
        "content": {
            "url": "http://0.0.0.0:9292/mcp",
            "transport": "streamable_http"
        }
    }
)
```

### VS Code (`mcp.json`)

To configure the nginx-filesystem-mcp for use in VS Code:

```json
{
  "servers": {
    "content": {
      "type": "http",
      "url": "http://0.0.0.0:9292/mcp"
    }
  }
}
```

---

## Tools

The following tools are available in nginx-filesystem-mcp:

### `listFolder`

- **Description:** Read a folder.
- **Input Schema:**
    - `path` (string, optional): Path of the directory to list.
- **Read-Only:** Yes

---

### `readFile`

- **Description:** Read a file.
- **Input Schema:**
    - `path` (string, required): Path to the file to read.
- **Read-Only:** Yes

---

### `treeFolder`

- **Description:** Returns a hierarchical JSON representation of a project's directory structure.
- **Input Schema:**
    - `path` (string, optional): Path of the directory to traverse.
- **Read-Only:** Yes


## Resource Templates

The following resource templates are supported by the nginx-filesystem-mcp:

### `files://{path}`

- **Name:** File
- **Description:** Return path information for a folder or a file.
- **MIME Type:** `application/json`  
