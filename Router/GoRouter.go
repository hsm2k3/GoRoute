package Router

import (
	"net/http"
	"net/url"
	"strings"
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
	routes []Route
}

// NewRouter creates a new Router.
func NewRouter() *Router {
	return &Router{
		routes: make([]Route, 0),
	}
}

// Add adds a new route to the router.
func (r *Router) Add(urlString, method string, handler http.HandlerFunc) {
	route := parseRouteFromURLString(urlString)
	route.Method = method
	route.Handler = handler
	r.routes = append(r.routes, route)
}

// ServeHTTP matches the request to a route and calls the handler for the route.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.Method == req.Method && route.Path == req.URL.Path {
			route.Handler(w, req)
			return
		}
	}
	http.NotFound(w, req)
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
