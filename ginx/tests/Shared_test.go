package tests

import (
    "net/http"

    "github.com/adeynack/gin-extensions/ginx"
    "github.com/gin-gonic/gin"
)

type PersonIn struct {
    FirstName string `json:"first_name" yaml:"first_name"`
    LastName  string `json:"last_name" yaml:"last_name"`
    BirthDate string `json:"birth_date" yaml:"birth_date"`
}
type PersonOut struct {
    Id        int64  `json:"id" yaml:"id"`
    FirstName string `json:"first_name" yaml:"first_name"`
    LastName  string `json:"last_name" yaml:"last_name"`
    BirthDate string `json:"birth_date" yaml:"birth_date"`
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
