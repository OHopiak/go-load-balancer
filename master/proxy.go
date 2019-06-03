package master

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// ProxyConfig defines the config for Proxy middleware.
	ProxyConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Balancer defines a load balancing technique.
		// Required.
		Balancer ProxyBalancer

		// Rewrite defines URL path rewrite rules. The values captured in asterisk can be
		// retrieved by index e.g. $1, $2 and so on.
		// Examples:
		// "/old":              "/new",
		// "/api/*":            "/$1",
		// "/js/*":             "/public/javascripts/$1",
		// "/users/*/orders/*": "/user/$1/order/$2",
		Rewrite map[string]string

		// Context key to store selected ProxyTarget into context.
		// Optional. Default value "target".
		ContextKey string

		// To customize the transport to remote.
		// Examples: If custom TLS certificates are required.
		Transport http.RoundTripper

		rewriteRegex map[*regexp.Regexp]string
	}

	// ProxyTarget defines the upstream target.
	ProxyTarget struct {
		Name string
		URL  *url.URL
		Meta echo.Map
	}

	// ProxyBalancer defines an interface to implement a load balancing technique.
	ProxyBalancer interface {
		AddTarget(*ProxyTarget) bool
		RemoveTarget(string) bool
		Next(echo.Context) *ProxyTarget
	}

	commonBalancer struct {
		targets []*ProxyTarget
		mutex   sync.RWMutex
	}

	// RandomBalancer implements a random load balancing technique.
	randomBalancer struct {
		*commonBalancer
		random *rand.Rand
	}

	// RoundRobinBalancer implements a round-robin load balancing technique.
	roundRobinBalancer struct {
		*commonBalancer
		i uint32
	}
)

var (
	// DefaultProxyConfig is the default Proxy middleware config.
	DefaultProxyConfig = ProxyConfig{
		Skipper:    middleware.DefaultSkipper,
		ContextKey: "target",
	}
)

func newCommonBalancer(targets []*ProxyTarget) *commonBalancer {
	b := new(commonBalancer)
	if targets != nil {
		b.targets = targets
	} else {
		b.targets = []*ProxyTarget{}
	}
	return b
}

// NewRandomBalancer returns a random proxy balancer.
func NewRandomBalancer(targets []*ProxyTarget) ProxyBalancer {
	b := &randomBalancer{commonBalancer: newCommonBalancer(targets)}
	return b
}

// NewRoundRobinBalancer returns a round-robin proxy balancer.
func NewRoundRobinBalancer(targets []*ProxyTarget) ProxyBalancer {
	b := &roundRobinBalancer{commonBalancer: newCommonBalancer(targets)}
	return b
}

// AddTarget adds an upstream target to the list.
func (b *commonBalancer) AddTarget(target *ProxyTarget) bool {
	for _, t := range b.targets {
		if t.Name == target.Name {
			return false
		}
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.targets = append(b.targets, target)
	return true
}

// RemoveTarget removes an upstream target from the list.
func (b *commonBalancer) RemoveTarget(name string) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for i, t := range b.targets {
		if t.Name == name {
			b.targets = append(b.targets[:i], b.targets[i+1:]...)
			return true
		}
	}
	return false
}

// Next randomly returns an upstream target.
func (b *randomBalancer) Next(c echo.Context) *ProxyTarget {
	if b.random == nil {
		b.random = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	}
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.targets[b.random.Intn(len(b.targets))]
}

// Next returns an upstream target using round-robin technique.
func (b *roundRobinBalancer) Next(c echo.Context) *ProxyTarget {
	if len(b.targets) == 0 {
		return nil
	}
	b.i = b.i % uint32(len(b.targets))
	t := b.targets[b.i]
	atomic.AddUint32(&b.i, 1)
	return t
}

// Proxy returns a Proxy middleware.
//
// Proxy middleware forwards the request to upstream server using a configured load balancing technique.
func Proxy(balancer ProxyBalancer) echo.MiddlewareFunc {
	c := DefaultProxyConfig
	c.Balancer = balancer
	return ProxyWithConfig(c)
}

// ProxyWithConfig returns a Proxy middleware with config.
// See: `Proxy()`
func ProxyWithConfig(config ProxyConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultProxyConfig.Skipper
	}
	if config.Balancer == nil {
		panic("echo: proxy middleware requires balancer")
	}
	config.rewriteRegex = map[*regexp.Regexp]string{}

	// Initialize rewrite configuration
	for k, v := range config.Rewrite {
		k = strings.Replace(k, "*", "(\\S*)", -1)
		config.rewriteRegex[regexp.MustCompile(k)] = v
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			tgt := config.Balancer.Next(c)
			if tgt == nil {
				return next(c)
			}

			for !ping(tgt.URL.String() + "/status"){
				config.Balancer.RemoveTarget(tgt.URL.String())
				tgt = config.Balancer.Next(c)
				if tgt == nil {
					return next(c)
				}
			}
			c.Set(config.ContextKey, tgt)

			// Rewrite url
			for k, v := range config.rewriteRegex {
				replacer := captureTokens(k, req.URL.Path)
				if replacer != nil {
					req.URL.Path = replacer.Replace(v)
				}
			}

			// Fix header
			if req.Header.Get(echo.HeaderXRealIP) == "" {
				req.Header.Set(echo.HeaderXRealIP, c.RealIP())
			}
			if req.Header.Get(echo.HeaderXForwardedProto) == "" {
				req.Header.Set(echo.HeaderXForwardedProto, c.Scheme())
			}

			errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
				config.Balancer.RemoveTarget(tgt.URL.String())
				//http.Redirect(w, r, c.Request().URL.String(), http.StatusSeeOther)
			}

			proxyHTTP(tgt, errorHandler).ServeHTTP(res, req)

			return nil
		}
	}
}

func proxyHTTP(t *ProxyTarget, errorHandler func(http.ResponseWriter, *http.Request, error)) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(t.URL)
	proxy.ErrorHandler = errorHandler
	return proxy
}

func captureTokens(pattern *regexp.Regexp, input string) *strings.Replacer {
	groups := pattern.FindAllStringSubmatch(input, -1)
	if groups == nil {
		return nil
	}
	values := groups[0][1:]
	replace := make([]string, 2*len(values))
	for i, v := range values {
		j := 2 * i
		replace[j] = "$" + strconv.Itoa(i+1)
		replace[j+1] = v
	}
	return strings.NewReplacer(replace...)
}

func ping(url string) bool {
	client := http.Client{
		Timeout: time.Second,
	}
	response, err := client.Get(url)
	return err == nil && response.StatusCode == http.StatusOK
}