package Router

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Route is a simple HTTP route that matches requests based on method and path.
type Route struct {
	Scheme    string
	Subdomain string
	Domain    string
	Path      string
	Method    string
	Handler   http.HandlerFunc
}

// Router is a simple HTTP router that matches requests based on method and path.
type Router struct {
	routes sync.Map
}

// NewRouter creates a new Router.
func NewRouter() *Router {
	return &Router{}
}

// Add adds a new route to the router.
func (r *Router) Add(urlString, method string, handler http.HandlerFunc) {
	route := parseRouteFromURLString(urlString)
	route.Method = method
	route.Handler = handler
	r.addRoute(route)
}

func (r *Router) addRoute(route Route) {
	key := getRouteKey(route.Method, route.Path)
	r.routes.Store(key, route)
}

// ServeHTTP matches the request to a route and calls the handler for the route.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := getRouteKey(req.Method, req.URL.Path)
	if route, ok := r.routes.Load(key); ok {
		route.(Route).Handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func parseRouteFromURLString(urlString string) Route {
	u, err := url.Parse(urlString)
	if err != nil {
		return Route{
			Path: urlString,
		}
	}

	route := Route{
		Scheme:    u.Scheme,
		Subdomain: getSubdomain(u.Host),
		Domain:    getDomain(u.Host),
		Path:      u.Path,
	}

	return route
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
	return method + ":" + path
}
