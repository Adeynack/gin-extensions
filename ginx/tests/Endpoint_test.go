package tests

import (
    "testing"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
    "github.com/adeynack/gin-extensions/ginx/materializer"
    "net/http"
    "github.com/adeynack/gin-extensions/ginx"
    "github.com/stretchr/testify/assert"
    "strings"
)

// GET request with a response body to serialize (no request body).
func TestGet(t *testing.T) {
    route := route(&materializer.JsonContentMaterializer{})
    rec := httptest.NewRecorder()
    req, _ := http.NewRequest(http.MethodGet, "/", nil)
    route.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t,
        `{"id":4573098657423896,"first_name":"Lilly","last_name":"Wachowski","birth_date":"1967-12-29"}`,
        rec.Body.String())
}

// POST request with a request body to deserialize and a response body to serialize.
func TestPost(t *testing.T) {
    route := route(&materializer.JsonContentMaterializer{})
    rec := httptest.NewRecorder()
    requestBody := `{"first_name":"Lana","last_name":"Wachowski","birth_date":"1965-06-21"}`
    req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
    route.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusCreated, rec.Code)
    assert.Equal(t,
        `{"id":542857589043,"first_name":"Lana","last_name":"Wachowski","birth_date":"1965-06-21"}`,
        rec.Body.String())
}

type PersonIn struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    BirthDate string `json:"birth_date"`
}
type PersonOut struct {
    Id        int64  `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    BirthDate string `json:"birth_date"`
}

func route(m ginx.ContentMaterializer) *gin.Engine {
    route := gin.Default()

    route.GET("/", m.Handler(func(xc ginx.Exchange) (*ginx.Conclusion, error) {
        result := PersonOut{
            Id:        4573098657423896,
            FirstName: "Lilly",
            LastName:  "Wachowski",
            BirthDate: "1967-12-29",
        }
        return &ginx.Conclusion{
            Status:       http.StatusOK,
            ResponseBody: result,
        }, nil
    }))

    route.POST("/", m.Handler(func(xc ginx.Exchange) (*ginx.Conclusion, error) {
        personIn := PersonIn{}
        if err := xc.Bind(&personIn); err != nil {
            return nil, err
        }
        result := PersonOut{
            Id:        542857589043,
            FirstName: personIn.FirstName,
            LastName:  personIn.LastName,
            BirthDate: personIn.BirthDate,
        }
        return &ginx.Conclusion{
            Status:       http.StatusCreated,
            ResponseBody: result,
        }, nil
    }))

    return route
}
