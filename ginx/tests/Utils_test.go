package tests

import (
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func assertBody(t *testing.T, r *httptest.ResponseRecorder, expected string) {
    if body := r.Body.String(); body != expected {
        if len(expected) == 0 {
            t.Errorf("Expected empty body, got %v", body)
        } else {
            t.Errorf("Expected body %v, got %v", expected, body)
        }
    }
}

func assertStatus(t *testing.T, c *gin.Context, expected int) {
    if status := c.Writer.Status(); status != expected {
        t.Errorf("Expected status %v, got %v", expected, status)
    }
}
