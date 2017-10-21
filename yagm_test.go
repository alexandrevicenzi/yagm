package yagm

import (
	"net/http"
	"testing"
)

type Response struct {
	Status int
}

func (r *Response) Header() http.Header {
	return http.Header{}
}

func (r *Response) Write(b []byte) (int, error) {
	return 0, nil
}

func (r *Response) WriteHeader(status int) {
	r.Status = status
}

func TestRegexp(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
	}

	mux := New()
	mux.HandleFunc("^/[a-zA-Z]+$", handler)

	request, _ := http.NewRequest("GET", "/foobar", nil)
	mux.ServeHTTP(nil, request)

	if !called {
		t.Fatal("Handler not called.")
	}
}

func TestNamedGroup(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true

		if name, ok := Param(r, "name"); ok {
			if name != "foobar" {
				t.Fatal("Param 'name' is not 'foobar'.")
			}
		} else {
			t.Fatal("Param 'name' empty.")
		}
	}

	mux := New()
	mux.HandleFunc("^/(?P<name>[a-zA-Z]+)$", handler)

	request, _ := http.NewRequest("GET", "/foobar", nil)
	mux.ServeHTTP(nil, request)

	if !called {
		t.Fatal("Handler not called.")
	}
}

func TestRequestCycle(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true

		if name, ok := Param(r, "name"); ok {
			if name != "foobar" {
				t.Fatal("Param 'name' is not 'foobar'.")
			}
		} else {
			t.Fatal("Param 'name' empty.")
		}
	}

	mux := New()
	mux.HandleFunc("^/(?P<name>[a-zA-Z]+)$", handler)

	request, _ := http.NewRequest("GET", "/foobar", nil)
	mux.ServeHTTP(nil, request)

	if !called {
		t.Fatal("Handler not called.")
	}

	if _, ok := Params(request); ok {
		t.Fatal("Request not deleted.")
	}
}

func TestNotFound(t *testing.T) {
	mux := New()
	request, _ := http.NewRequest("GET", "/foobar", nil)
	response := &Response{}
	mux.ServeHTTP(response, request)

	if response.Status != 404 {
		t.Fatal("Status should be 404.")
	}
}


func TestRouteOptimization(t *testing.T){
	mux := New()
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	mux.HandleFunc("^/A[a-zA-Z]+$", handler)
	mux.HandleFunc("^/B[a-zA-Z]+$", handler)

	if mux.routes[0].pattern != "^/A[a-zA-Z]+$" {
		t.Fatal("Expected first item in the array to be route A")
	}
	if mux.routes[0].hitCount != 0 {
		t.Fatal("Expected route /A* to have a hit count of zero ")
	}
	if mux.routes[1].hitCount != 0 {
		t.Fatal("Expected route /B* to have a hit count of zero ")
	}

	requestA, _ := http.NewRequest("GET", "/Aardvark", nil)
	response := &Response{}
	for i:=0; i<100; i++ {
		mux.ServeHTTP(response, requestA)
	}
	if mux.routes[0].pattern != "^/A[a-zA-Z]+$" {
		t.Fatal("Expected first item in the array to be route A")
	}
	if mux.routes[0].hitCount != 100 {
		t.Fatal("Expected route /A* to have a hit count of 100 ", mux.routes[0].hitCount)
	}
	if mux.routes[1].hitCount != 0 {
		t.Fatal("Expected route /B* to have a hit count of zero ", mux.routes[1].hitCount)
	}


	requestB, _ := http.NewRequest("GET", "/BananasBananas", nil)
	for i:=0; i<200; i++ {
		mux.ServeHTTP(response, requestB)
	}

	if mux.routes[0].pattern != "^/A[a-zA-Z]+$" {
		t.Fatal("Expected first item in the array to be route A")
	}
	if mux.routes[0].hitCount != 100 {
		t.Fatal("Expected route /A* to have a hit count of 100", mux.routes[0].hitCount)
	}
	if mux.routes[1].hitCount != 200 {
		t.Fatal("Expected route /B* to have a hit count of 200 ", mux.routes[1].hitCount)
	}

	RouteOptimizeRequestCount = 250
	for i:=0; i<100; i++ {
		mux.ServeHTTP(response, requestB)
	}

	if mux.routes[0].pattern != "^/B[a-zA-Z]+$" {
		t.Fatal("Expected first item in the array to be route B")
	}
	if mux.routes[0].hitCount != 300 {
		t.Fatal("Expected route /A* to have a hit count of 300", mux.routes[0].hitCount)
	}
	if mux.routes[1].hitCount != 100 {
		t.Fatal("Expected route /B* to have a hit count of 100 ", mux.routes[1].hitCount)
	}

}