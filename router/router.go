package router

import (
	"delonakc.com/api/database"
	"delonakc.com/api/util"
	"fmt"
	"net/http"
)

type Router struct {
	// mongodb session
	DB *database.MongoDB

	// routes
	routes []*Route
}

type RouteMatch struct {
	Handler http.Handler
	Vars map[string]string
	MatchError error
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	if r, ok := r.Match(req); ok {
		req = r.SetMatch(req)
		r.Handler().ServeHTTP(w, req)
		return
	}
	PageNotFound(w, req)
}

func (r *Router) HandleFunc(path string, f http.Handler) *Route {
	return r.NewRoute().Path(path).HandleFunc(f)
}

func (r *Router) Match(req *http.Request) (*Route, bool) {
	for _, route := range r.routes {
		if route.Match(req) {
			return route, true
		}
	}
	return nil, false
}

func (r *Router) NewRoute() *Route {
	route := &Route{ DB: r.DB }
	r.routes = append(r.routes, route)
	return route
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "不支持该方法")
}

func PageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, util.ResponseError("page not found."))
}