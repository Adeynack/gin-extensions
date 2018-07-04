package ginx

import (
    "github.com/gin-gonic/gin"
    "github.com/go-http-utils/negotiator"
    "log"
)

type Exchange interface {
    Context() *gin.Context
    Negotiator() *negotiator.Negotiator
    ConsumedType() string
    ProducedType() string
    Bind(requestBody interface{}) (err error)
    AutoBind(requestBody interface{}) (handled bool)
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

func (xc *exchange) Bind(requestBody interface{}) (err error) {
    return xc.materializer.ReadRequestBody(xc, requestBody)
}

func (xc *exchange) AutoBind(requestBody interface{}) (handled bool) {
    if err := xc.Bind(requestBody); err != nil {
        log.Printf("Unable to bind request body: %v", err)
        answerWithProblem(xc.Context(), err)
        handled = true
    }
    return
}
