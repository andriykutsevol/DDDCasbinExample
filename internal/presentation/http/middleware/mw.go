package middleware

import (
	"fmt"
	"strings"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"

	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http"

	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/response"

	"github.com/gin-gonic/gin"
)

// But how to wait for http response in asynchronous rest api?
// ChatGPT
// In an asynchronous REST API, clients typically do not wait for the HTTP response immediately after making a request.
// Instead, they receive an immediate acknowledgment that the request has been accepted and a mechanism
// to track the status or retrieve the result of the asynchronous operation later.

// Here are some common approaches to handle waiting for HTTP responses in an asynchronous REST API:

// Polling:
// Clients periodically query the API to check the status of the asynchronous operation.
// The API returns the current status (e.g., "pending", "processing", "completed") along with any relevant information.
// Clients continue to poll the API until the operation is completed or reaches a terminal state.

// Webhooks:
// Clients provide a callback URL to the API when initiating the asynchronous operation. Once the operation is completed,
// the API sends a notification to the callback URL, informing the client of the result.
// Webhooks allow for real-time notifications without the need for continuous polling.

// Long polling:
// Similar to polling, but the server holds the request open until there is new information to send back to the client.
// This reduces the latency compared to periodic polling but may consume server resources for long-lived connections.

// WebSocket:
// Clients establish a persistent bidirectional connection with the server using WebSocket protocol.
// The server can push updates or notifications to the client in real-time as the asynchronous operation progresses.

// HTTP streaming:
//	Similar to WebSocket, but using HTTP/2 or HTTP/3 server push to stream updates to the client in real-time.

// Each approach has its advantages and disadvantages, and the choice depends on factors such as the nature of the asynchronous operation,
// the expected latency, the client's capabilities, and the infrastructure constraints.
// It's essential to consider these factors carefully when designing the asynchronous
// API and choose the approach that best fits the requirements of your application.

// ConcurrentHandlerMiddleware is a middleware that wraps the handler to execute it concurrently.
func ConcurrentHandlerMiddleware(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ch := make(chan struct {
		//     StatusCode int
		//     Body       gin.H
		// })

		ch := make(chan response.DemoOk)

		c.Set("channel", ch)

		// Start a new goroutine to handle the request concurrently
		go func() {
			handlerFunc(c)
			// Call the next handler in the chain
			//c.Next()
		}()

		// Wait for the response data from the channel
		fmt.Println("Wait for the response data from the channel")
		resp := <-ch

		//http.ResSuccess(c, resp)

		c.JSON(200, resp.Message)

	}
}

func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		http.ResError(c, errors.ErrMethodNotAllow)
	}
}

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		http.ResError(c, errors.ErrNotFound)
	}
}

type SkipperFunc func(*gin.Context) bool

func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

func AllowMethodAndPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := JoinRouter(c.Request.Method, c.Request.URL.Path)
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}

func SkipHandler(c *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}

func EmptyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
