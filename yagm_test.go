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
    respose := &Response{}
    mux.ServeHTTP(respose, request)

    if respose.Status != 404 {
        t.Fatal("Status should be 404.")
    }
}
