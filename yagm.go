package yagm

import (
	"net/http"
	"regexp"
	"sort"
	"sync"
	"sync/atomic"
)

var RouteOptimizeRequestCount uint64 = 1000

type yagmRoute struct {
	pattern  string
	handler  http.Handler
	re       *regexp.Regexp
	hitCount uint64 // incremented every time the route is used
}

type yagmRequest struct {
	route           *yagmRoute
	params          map[string]string
	paramsProcessed bool
}

// YagmMux is an HTTP request multiplexer.
// It matches the URL of each incoming request
// against a list of registered patterns and calls
// the handler for the pattern that matches the URL.
type YagmMux struct {
	mu           sync.RWMutex
	routes       []*yagmRoute
	requestCount uint64
}

var (
	// Hold request info until request is finished.
	requests = make(map[*http.Request]*yagmRequest)
)

// Does path match pattern?
func (route *yagmRoute) match(path string) bool {
	return route.re.MatchString(path)
}

// Extrac all named groups.
func (route *yagmRoute) extractParams(r *http.Request) map[string]string {
	result := make(map[string]string, 0)
	match := route.re.FindStringSubmatch(r.URL.Path)

	for i, name := range route.re.SubexpNames() {
		result[name] = match[i]
	}

	return result
}

// New allocates and returns a new YagmMux.
func New() *YagmMux {
	return &YagmMux{
		routes: make([]*yagmRoute, 0),
	}
}

// findRoute iterates the routes to find the pattern
func (mux *YagmMux) findRoute(pattern string) (*yagmRoute, bool) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	for _, r := range mux.routes {
		if r.pattern == pattern {
			return r, true
		}
	}
	return nil, false
}

// optimizeRoutes sorts the routes descending based on the route.hitCount, thus making the most used be first in the array
func (mux *YagmMux) optimizeRoutes() {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	sort.Slice(mux.routes, func(i, j int) bool {
		return mux.routes[i].hitCount > mux.routes[j].hitCount
	})
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
// If a pattern isn't a valid regular expression, Handle panics.
func (mux *YagmMux) Handle(pattern string, handler http.Handler) {

	if pattern == "" {
		panic("yagm: empty pattern")
	}

	if handler == nil {
		panic("yagm: nil handler")
	}

	if _, ok := mux.findRoute(pattern); ok {
		panic("yagm: route already registered")
	} else {
		re := regexp.MustCompile(pattern)

		route := &yagmRoute{
			pattern,
			handler,
			re,
			0,
		}
		mux.mu.Lock() // write lock for minimum time.
		mux.routes = append(mux.routes, route)
		mux.mu.Unlock()
	}
}

// HandleFunc registers the handler function for the given pattern.
func (mux *YagmMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.Handle(pattern, http.HandlerFunc(handler))
}

// ServeHTTP dispatches the request to the handler whose
// pattern matches the request URL.
func (mux *YagmMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rc := atomic.AddUint64(&mux.requestCount, 1)
	var handler http.Handler

	for _, route := range mux.routes {
		if route.match(r.URL.Path) {
			atomic.AddUint64(&route.hitCount, 1)
			handler = route.handler

			requests[r] = &yagmRequest{
				route: route,
			}

			break
		}
	}

	// Delete request info after processing request.
	defer func() {
		delete(requests, r)
	}()

	if handler == nil {
		handler = http.NotFoundHandler()
	}

	handler.ServeHTTP(w, r)

	if rc > RouteOptimizeRequestCount {
		mux.optimizeRoutes()
		atomic.StoreUint64(&mux.requestCount, 0)
	}
}

// Param return the value of a route variable.
func Param(r *http.Request, name string) (string, bool) {
	var value string

	params, ok := Params(r)

	if ok {
		value, ok = params[name]
	}

	return value, ok
}

// Params returns the route variables for the current request, if any.
func Params(r *http.Request) (map[string]string, bool) {
	req, ok := requests[r]

	if !ok {
		return nil, false
	}

	// Only process params if requested.
	// This save some request time if not used.
	if !req.paramsProcessed {
		req.params = req.route.extractParams(r)
		req.paramsProcessed = true
	}

	return req.params, ok
}
