package yagm

import (
    "net/http"
    "testing"
)

func BenchmarkYagmMuxRegexp(b *testing.B) {
    handler := func(w http.ResponseWriter, r *http.Request) {}
    mux := New()
    mux.HandleFunc("^/[a-zA-Z]+$", handler)

    request, _ := http.NewRequest("GET", "/foobar", nil)
    for i := 0; i < b.N; i++ {
        mux.ServeHTTP(nil, request)
    }
}

func BenchmarkYagmMuxNamedGroup(b *testing.B) {
    handler := func(w http.ResponseWriter, r *http.Request) {}
    mux := New()
    mux.HandleFunc("^/(?P<name>[a-zA-Z]+)$", handler)

    request, _ := http.NewRequest("GET", "/foobar", nil)
    for i := 0; i < b.N; i++ {
        mux.ServeHTTP(nil, request)
    }
}
