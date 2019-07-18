package main

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handler).Methods("GET")
	router.HandleFunc("/nodes", nodeHandler).Methods("GET")
	router.HandleFunc("/cni", cniHandler).Methods("GET")
	return router
}
func TestServer(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "Expected Output")
	request, _ = http.NewRequest("GET", "/nodes", nil)
	response = httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "Expected Output")
	request, _ = http.NewRequest("GET", "/cni", nil)
	response = httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "Expected Output")
	go StartRestServer()
}
