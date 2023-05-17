package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

// RequestID is a middleware that injects a 'RequestID' into the context and header of each request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		xRequestID := c.GetHeader("X-Request-ID")
		if xRequestID == "" {
			xRequestID = uuid.NewString()
		}
		c.Set("requestId", xRequestID)
		c.Header("X-Request-ID", xRequestID)
		fmt.Printf("[GIN-debug] %s [%s] - \"%s %s\"\n", time.Now().Format(time.RFC3339), xRequestID, c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}
