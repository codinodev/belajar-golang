package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	router := httprouter.New()
	router.GET("/", func(wirter http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(wirter, "Hello World")
	})
	request := httptest.NewRequest("GET", "http://localhost:3000/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Hello World", string(body))
}

func TestRouterParams(t *testing.T) {
	router := httprouter.New()
	router.GET("/product/:id", func(wirter http.ResponseWriter, request *http.Request, params httprouter.Params) {
		text := "Product " + params.ByName("id")
		fmt.Fprint(wirter, text)
	})
	request := httptest.NewRequest("GET", "http://localhost:3000/product/1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Product 1", string(body))
}

func TestRouterPatternNamedParams(t *testing.T) {
	router := httprouter.New()
	router.GET("/products/:id/items/:itemid", func(wirter http.ResponseWriter, request *http.Request, params httprouter.Params) {
		text := "Product " + params.ByName("id") + " Item " + params.ByName("itemid")
		fmt.Fprint(wirter, text)
	})
	request := httptest.NewRequest("GET", "http://localhost:3000/products/1/items/1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Product 1 Item 1", string(body))
}

func TestRouterPatternNamedParamsCatchAll(t *testing.T) {
	router := httprouter.New()
	router.GET("/images/*image", func(wirter http.ResponseWriter, request *http.Request, params httprouter.Params) {
		text := "Image : " + params.ByName("image")
		fmt.Fprint(wirter, text)
	})
	request := httptest.NewRequest("GET", "http://localhost:3000/images/small/profile.png", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Image : /small/profile.png", string(body))
}

//go:embed resources
var resources embed.FS

func TestServerFile(t *testing.T) {
	router := httprouter.New()
	directory, _ := fs.Sub(resources, "resources")
	router.ServeFiles("/files/*filepath", http.FS(directory))
	request := httptest.NewRequest("GET", "http://localhost:3000/files/goodbye.txt", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "GoodBye", string(body))
}

func TestPanicHandler(t *testing.T) {
	router := httprouter.New()
	router.PanicHandler = func(wirter http.ResponseWriter, request *http.Request, error interface{}) {
		fmt.Fprint(wirter, "Panic : ", error)
	}
	router.GET("/", func(wirter http.ResponseWriter, request *http.Request, params httprouter.Params) {
		panic("Ups")
	})
	request := httptest.NewRequest("GET", "http://localhost:3000/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Panic : Ups", string(body))
}

func TestPanicHandlerNotFound(t *testing.T) {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "Gak Ketemu")
	})
	request := httptest.NewRequest("GET", "http://localhost:3000/404", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Gak Ketemu", string(body))
}

func TestMethodNotFoundAllow(t *testing.T) {
	router := httprouter.New()
	router.MethodNotAllowed = http.HandlerFunc(func(wirter http.ResponseWriter, request *http.Request) {
		fmt.Fprint(wirter, "Gak Boleh")
	})
	router.POST("/", func(wirter http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(wirter, "POST")
	})
	request := httptest.NewRequest("GET", "http://localhost:3000/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Gak Boleh", string(body))
}

type LogMiddleware struct {
	http.Handler
}

func (middleware *LogMiddleware) ServeHTTP(wirter http.ResponseWriter, request *http.Request) {
	fmt.Println("Receive request")
	middleware.Handler.ServeHTTP(wirter, request)
}

func TestMiddleware(t *testing.T) {
	router := httprouter.New()
	router.GET("/", func(wirter http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(wirter, "Middleware")
	})
	middleware := LogMiddleware{Handler: router}
	request := httptest.NewRequest("GET", "http://localhost:3000/", nil)
	recorder := httptest.NewRecorder()

	middleware.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Middleware", string(body))
}
