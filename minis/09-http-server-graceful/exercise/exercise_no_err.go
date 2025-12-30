//go:build !solution
// +build !solution

package exercise

/*
Core HTTP Server Algorithm - Simplified Version
================================================

This file focuses on the core HTTP server patterns without error handling.
Use this to understand the fundamental request handling flow.

ALGORITHM STEPS:

1. Request Routing:
   - Match HTTP method (GET, POST, etc.)
   - Route to appropriate handler
   - Extract parameters (query params, body)

2. Middleware Pattern:
   - Create wrapper function
   - Execute pre-handler logic
   - Call next handler
   - Execute post-handler logic

3. Graceful Shutdown:
   - Receive shutdown signal
   - Stop accepting new connections
   - Wait for active requests to complete
   - Close server

DEBUGGING GUIDE:
- Set breakpoint at handler entry to trace requests
- Watch HTTP method to see routing decisions
- Watch middleware execution order
- Set breakpoint at shutdown to see graceful close
*/

// TODO: Implement SimpleRouteHandler
// Core algorithm for routing HTTP requests by method
// Steps:
// 1. Check r.Method
// 2. If POST: decode body, store data
// 3. If GET: extract params, retrieve data
// 4. Return appropriate response
// Signature: func SimpleRouteHandler(method string, data interface{}) interface{}

// TODO: Implement SimpleMiddleware
// Core algorithm for middleware wrapping
// Steps:
// 1. Accept handler function as input
// 2. Return new function that:
//    - Executes pre-logic (e.g., increment counter)
//    - Calls original handler
//    - Executes post-logic (if needed)
// Signature: func SimpleMiddleware(handler func()) func()

// TODO: Implement SimpleGracefulShutdown
// Core algorithm for graceful server shutdown
// Steps:
// 1. Stop accepting new connections
// 2. Wait for active connections to finish
// 3. Close server resources
// Signature: func SimpleGracefulShutdown(activeConnections int) bool
