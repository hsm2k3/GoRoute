package Router

import (
	"net/http"
)

// Route is a simple HTTP route that matches requests based on method and path.
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
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
