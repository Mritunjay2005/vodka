package vodka

import "net/http"

func AllowCORS(origins []string) HandlerFunc {
	return func(c *Context) {
		clientOrigin := c.Request.Header.Get("Origin")
		allowThisOrigin := ""

		for _, o := range origins {
			if o == "*" || o == clientOrigin {
				allowThisOrigin = clientOrigin

				if o == "*" {
					allowThisOrigin = "*"
				}
				break
			}
		}

		if allowThisOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowThisOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		}

		if c.Request.Method == http.MethodOptions {
			if allowThisOrigin != "" {
				c.Writer.WriteHeader(http.StatusNoContent) // 204: All good, proceed!
			} else {
				c.Writer.WriteHeader(http.StatusForbidden) // 403: Origin not allowed
			}

			c.Abort()
			return
		}

		c.Next()
	}
}
