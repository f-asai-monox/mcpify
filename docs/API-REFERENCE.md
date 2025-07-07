# API Reference

## Available Tools

The MCP Bridge server automatically converts REST API endpoints into MCP tools. The available tools depend on your configuration file.

## Default Tools (Example Configuration)

### User Management
- `health` - Health check endpoint
- `get_users` - Get all users  
- `create_user` - Create a new user
- `get_user` - Get a specific user by ID
- `update_user` - Update a user by ID
- `delete_user` - Delete a user by ID

### Product Management
- `get_products` - Get all products
- `get_product` - Get a specific product by ID

## Usage Examples

### Health Check
```javascript
await callTool("health", {});
```

### User Operations
```javascript
// Get all users
await callTool("get_users", {});

// Create new user
await callTool("create_user", {
  name: "John Doe",
  email: "john@example.com"
});

// Get specific user
await callTool("get_user", {
  id: 1
});

// Update user
await callTool("update_user", {
  id: 1,
  name: "Jane Doe",
  email: "jane@example.com"
});

// Delete user
await callTool("delete_user", {
  id: 1
});
```

### Product Operations
```javascript
// Get all products
await callTool("get_products", {});

// Get specific product
await callTool("get_product", {
  id: 1
});
```

## Parameter Types

### Path Parameters
Used in URL paths (e.g., `/users/{id}`)
- Automatically extracted from the path
- Usually numeric IDs

### Body Parameters
Sent in request body for POST/PUT operations
- JSON format
- Can be strings, numbers, booleans, or objects

### Query Parameters
Added to URL as query strings
- Optional parameters for filtering/pagination
- Format: `?param1=value1&param2=value2`

### Header Parameters
Custom headers sent with requests
- Authentication tokens
- Content-Type specifications
- Custom API headers

## Response Format

All tools return JSON responses from the target REST API. The response structure depends on the specific API endpoint.

Common response patterns:
- **Success**: JSON object or array with data
- **Error**: JSON object with error details
- **Empty**: Empty response for DELETE operations

## Resources

The MCP server also provides resources for API documentation:

### Available Resources
- `rest-api://docs` - Complete REST API specification in JSON format

### Resource Usage
```javascript
// Get API documentation
const docs = await getResource("rest-api://docs");
```

## Error Handling

The bridge handles common HTTP errors and converts them to MCP error responses:

- **400 Bad Request**: Invalid parameters
- **401 Unauthorized**: Authentication required
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server error

Error responses include:
- Error code
- Error message
- Original HTTP status code
- Additional context when available

## Tool Discovery

To see what tools are available in your specific configuration:

1. Use the MCP `tools/list` method
2. Check the `rest-api://docs` resource
3. Review your configuration file's `endpoints` section

Each tool includes:
- Tool name
- Description
- Parameter schema
- Required vs optional parameters