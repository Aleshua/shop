// internal/proxy/proxy.go

package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewReverseProxyHandler создает Gin.HandlerFunc для проксирования на targetURL.
func NewReverseProxyHandler(targetURL, prefixToStrip string) gin.HandlerFunc {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefixToStrip)

		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
