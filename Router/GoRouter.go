package Router

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Route represents an HTTP route with path parameters and middleware support
type Route struct {
	Scheme     string
	Subdomain  string
	Domain     string
	Path       string
	Method     string
	Handler    http.HandlerFunc
	Middleware []MiddlewareFunc
	params     map[string]string
}

// MiddlewareFunc defines the contract for middleware
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// Router is an HTTP router with middleware and path parameter support
type Router struct {
	routes     sync.Map
	middleware []MiddlewareFunc
	notFound   http.HandlerFunc
}

// NewRouter creates a new Router instance
func NewRouter() *Router {
	return &Router{
		notFound: http.NotFound,
	}
}

// Use adds middleware to the router
func (r *Router) Use(middleware ...MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware...)
}

// SetNotFound sets a custom handler for 404 responses
func (r *Router) SetNotFound(handler http.HandlerFunc) {
	r.notFound = handler
}

// Add registers a new route with the router
func (r *Router) Add(urlString, method string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	route := parseRouteFromURLString(urlString)
	route.Method = method
	route.Handler = handler
	route.Middleware = middleware
	route.params = make(map[string]string)
	r.addRoute(route)
}

func (r *Router) addRoute(route Route) {
	key := getRouteKey(route.Method, route.Path)
	r.routes.Store(key, route)
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var matchedRoute Route
	var params map[string]string
	found := false

	r.routes.Range(func(key, value interface{}) bool {
		route := value.(Route)
		if matches, routeParams := matchRoute(route.Path, req.URL.Path); matches && route.Method == req.Method {
			matchedRoute = route
			params = routeParams
			found = true
			return false
		}
		return true
	})

	if !found {
		r.notFound(w, req)
		return
	}

	// Apply global middleware
	handler := matchedRoute.Handler
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}

	// Apply route-specific middleware
	for i := len(matchedRoute.Middleware) - 1; i >= 0; i-- {
		handler = matchedRoute.Middleware[i](handler)
	}

	// Add params to request context
	ctx := req.Context()
	for k, v := range params {
		ctx = context.WithValue(ctx, k, v)
	}
	req = req.WithContext(ctx)

	handler(w, req)
}

func parseRouteFromURLString(urlString string) Route {
	u, err := url.Parse(urlString)
	if err != nil {
		return Route{
			Path: urlString,
		}
	}

	return Route{
		Scheme:    u.Scheme,
		Subdomain: getSubdomain(u.Host),
		Domain:    getDomain(u.Host),
		Path:      u.Path,
	}
}

func matchRoute(pattern, path string) (bool, map[string]string) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			// This is a parameter
			paramName := strings.TrimPrefix(part, ":")
			params[paramName] = pathParts[i]
		} else if part != pathParts[i] {
			return false, nil
		}
	}

	return true, params
}

func getSubdomain(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		return strings.Join(parts[:len(parts)-2], ".")
	}
	return ""
}

func getDomain(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) > 1 {
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return host
}

func getRouteKey(method, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}
