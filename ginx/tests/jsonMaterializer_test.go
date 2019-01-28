package tests

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/adeynack/gin-extensions/ginx/materializer"
    "github.com/go-http-utils/headers"
    "github.com/stretchr/testify/assert"
)

// GET request with a response body to serialize (no request body), `Accept` header not provided
func TestGetNoAccept(t *testing.T) {
    route := route(&materializer.JsonContentMaterializer{})
    rec := httptest.NewRecorder()
    req, err := http.NewRequest(http.MethodGet, "/", nil)
    assert.Error(t, err)
    route.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, []string{"application/json"}, rec.HeaderMap[headers.ContentType])
    assert.Equal(t,
        `{"id":4573098657423896,"first_name":"Lilly","last_name":"Wachowski","birth_date":"1967-12-29"}`,
        rec.Body.String())
}

// GET request with a response body to serialize (no request body), with `Accept` header = `application/json`
func TestGetAcceptJSON(t *testing.T) {
    route := route(&materializer.JsonContentMaterializer{})
    rec := httptest.NewRecorder()
    req, _ := http.NewRequest(http.MethodGet, "/", nil)
    req.Header.Add(headers.Accept, "application/json")
    route.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, []string{"application/json"}, rec.HeaderMap[headers.ContentType])
    assert.Equal(t,
        `{"id":4573098657423896,"first_name":"Lilly","last_name":"Wachowski","birth_date":"1967-12-29"}`,
        rec.Body.String())
}

// GET request with a response body to serialize (no request body), with `Accept` header = non supported type
func TestGetAcceptNotSupported(t *testing.T) {
    route := route(&materializer.JsonContentMaterializer{})
    rec := httptest.NewRecorder()
    req, _ := http.NewRequest(http.MethodGet, "/", nil)
    req.Header.Add(headers.Accept, "foo/bar")
    route.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
    assert.Equal(t, []string{"application/problem+json"}, rec.HeaderMap[headers.ContentType])
    assert.Equal(t,
        `{"status":415,"title":"The accepted types are not supported (header 'Accept')","detail":"A request that accepts type(s) 'foo/bar' is not supported."}`,
        rec.Body.String())
}

// POST request with a request body to deserialize and a response body to serialize, `content-type` not provided
func TestPostNoContentType(t *testing.T) {
    route := route(&materializer.JsonContentMaterializer{})
    rec := httptest.NewRecorder()
    requestBody := `{"first_name":"Lana","last_name":"Wachowski","birth_date":"1965-06-21"}`
    req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
    route.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusCreated, rec.Code)
    assert.Equal(t, []string{"application/json"}, rec.HeaderMap[headers.ContentType])
    assert.Equal(t,
        `{"id":542857589043,"first_name":"Lana","last_name":"Wachowski","birth_date":"1965-06-21"}`,
        rec.Body.String())
}

// POST request with a request body to deserialize and a response body to serialize, `content-type` header = `application/json`
func TestPostContentTypeJSON(t *testing.T) {
    route := route(&materializer.JsonContentMaterializer{})
    rec := httptest.NewRecorder()
    requestBody := `{"first_name":"Lana","last_name":"Wachowski","birth_date":"1965-06-21"}`
    req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
    req.Header.Add(headers.ContentType, "application/json")
    route.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusCreated, rec.Code)
    assert.Equal(t, []string{"application/json"}, rec.HeaderMap[headers.ContentType])
    assert.Equal(t,
        `{"id":542857589043,"first_name":"Lana","last_name":"Wachowski","birth_date":"1965-06-21"}`,
        rec.Body.String())
}

// POST request with a request body to deserialize and a response body to serialize, `content-type` header = non supported type
func TestPostContentTypeNotSupported(t *testing.T) {
    route := route(&materializer.JsonContentMaterializer{})
    rec := httptest.NewRecorder()
    requestBody := `{"first_name":"Lana","last_name":"Wachowski","birth_date":"1965-06-21"}`
    req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
    req.Header.Add(headers.ContentType, "foo/bar")
    route.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusNotAcceptable, rec.Code)
    assert.Equal(t, []string{"application/problem+json"}, rec.HeaderMap[headers.ContentType])
    assert.Equal(t,
        `{"status":406,"title":"The request type is not supported (header 'Content-Type')","detail":"A request of type 'foo/bar' is not supported."}`,
        rec.Body.String())
}
