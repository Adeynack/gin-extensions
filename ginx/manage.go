package ginx

import (
    "github.com/gin-gonic/gin"
    "github.com/go-http-utils/negotiator"
    "github.com/go-http-utils/headers"
    "net/http"
    "fmt"
    "log"
)

type Conclusion interface {
    Status() int
    ResponseBody() interface{}
}

// Function managing the exchange.
//
// If returns an error, it will be automatically used to answer the exchange.
//
// If error is nil, but a conclusion is returned, it will be used to automate
// the content negotiation and the answer to the user.
//
// If both error and conclusion are nil, it will be assumed that the request
// was manually handled inside of the manager and nothing needs to be performed.
//
type Handler = func(Exchange) (Conclusion, error)

// Handles a `gin.Context` the `ginx` way.
//
// ctx              The base `gin.Context` object.
// materializer     The logic of deserializing request and serializing response.
// handler          The function handling the exchange and returning either an
//                  exchange `Conclusion` or an error.
func Handle(ctx *gin.Context, materializer ContentMaterializer, handler Handler) {
    exchange := &exchange{
        context:    ctx,
        negotiator: negotiator.New(ctx.Request.Header),
    }
    //
    // Ensure request's content type is supported.
    //
    acceptedConsumedType := ""
    contentType := ctx.GetHeader(headers.ContentType)
    if contentType != "" {
        // if a `accept` header is provided, it needs to be supported by the negotiator.
        acceptedConsumedType = exchange.Negotiator().Type(materializer.ProducedTypes()...)
        if acceptedConsumedType == "" {
            answerWithProblem(ctx, Problem{
                Status: http.StatusNotAcceptable,
                Title:  fmt.Sprintf("The request type is not supported (header `%s`)", headers.ContentType),
                Detail: fmt.Sprintf("A request of type `%s` is not supported.", contentType),
            })
            return
        }
    }
    //
    // Ensure accepted type (response type) is supported.
    //
    acceptedProducedType := ""
    accept := ctx.GetHeader(headers.Accept)
    if accept != "" {
        // if a `accept` header is provided, it needs to be supported by the negotiator.
        acceptedProducedType = exchange.Negotiator().Type(materializer.ProducedTypes()...)
        if acceptedProducedType == "" {
            answerWithProblem(ctx, Problem{
                Status: http.StatusNotAcceptable,
                Title:  fmt.Sprintf("The accepted types are not supported (header `%s`)", headers.Accept),
                Detail: fmt.Sprintf("A request that accepts type(s) `%s` is not supported.", accept),
            })
            return
        }
    }
    //
    // Handle the exchange (call client code) and write the response.
    //
    conclusion, err := handler(exchange)
    if err != nil {
        answerWithProblem(ctx, err)
        return
    }
    materializer.WriteResponseBody(exchange, conclusion)
}

func answerWithProblem(ctx *gin.Context, err error) {
    problem, errIsProblem := err.(Problem)
    if !errIsProblem {
        log.Printf("Exchange manager returned a pure error. %v", err)
        problem = Problem{
            Status: http.StatusInternalServerError,
        }
    }
    ctx.Header(headers.ContentType, contentTypeProblemJson)
    ctx.JSON(problem.Status, problem)
}
