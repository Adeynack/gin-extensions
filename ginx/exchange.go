package ginx

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-http-utils/negotiator"
)

type Exchange interface {
    // Gets the Gin context.
    Context() *gin.Context

    // Gets the Negotiator for this exchange.
    Negotiator() *negotiator.Negotiator

    // Gets the list of types accepted in the request.
    ConsumedType() string // todo: Used?

    // Gets the list of types accepted for the response.
    ProducedType() string // todo: Used?

    // Bind the request's body to a struct.
    // When it fails, returns a `Problem` with a 400 BAD REQUEST status code.
    Bind(requestBody interface{}) error
}

type exchange struct {
    context      *gin.Context
    negotiator   *negotiator.Negotiator
    consumedType string
    producedType string
    materializer ContentMaterializer
}

var _ Exchange = &exchange{}

func (xc *exchange) Context() *gin.Context {
    return xc.context
}

func (xc *exchange) Negotiator() *negotiator.Negotiator {
    return xc.negotiator
}

func (xc *exchange) ConsumedType() string {
    return xc.consumedType
}

func (xc *exchange) ProducedType() string {
    return xc.producedType
}

func (xc *exchange) Bind(requestBody interface{}) error {
    if err := xc.materializer.ReadRequestBody(xc, requestBody); err != nil {
        log.Printf("Unable to bind request body: %v", err)
        return &Problem{
            Status: http.StatusBadRequest,
            Title:  "Request body could not be parsed.",
            Cause:  err,
        }
    }
    return nil
}
