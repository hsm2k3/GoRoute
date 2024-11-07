package Router

import (
	"net/http"
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
func (r *Router) Add(method, path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
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
	parts := strings.Split(urlString, "://")
	scheme := parts[0]

	if len(parts) > 1 {
		hostParts := strings.Split(parts[1], "/")
		domainParts := strings.Split(hostParts[0], ".")
		if len(domainParts) > 2 {
			route := Route{
				Scheme:    scheme,
				Subdomain: strings.Join(domainParts[:len(domainParts)-2], "."),
				Domain:    strings.Join(domainParts[len(domainParts)-2:], "."),
				Path:      "/" + strings.Join(hostParts[1:], "/"),
			}
			return route
		}
	}

	return Route{
		Scheme: scheme,
		Path:   urlString,
	}
}
