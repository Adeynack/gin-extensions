package ginx

import (
    "github.com/gin-gonic/gin"
    "github.com/go-http-utils/negotiator"
    "github.com/go-http-utils/headers"
    "net/http"
    "fmt"
    "log"
    "strings"
)

type Conclusion struct {
    Status       int
    ResponseBody interface{}
}

// todo: Try `Conclusion` as an interface instead...
/*
type Conclusion interface {

    // Gets any reason why this exchange could not complete (`nil` if successful).
    Error() error

    // Gets the resulting HTTP status code.
    Status() int

    // Gets the object to be serialized as the response's body.
    ResponseBody() interface{}
}
*/

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
type HandlerFunc = func(Exchange) (*Conclusion, error)

// Generates a Gin `HandlerFunc` from a content materializer.
func Handler(materializer ContentMaterializer, handler HandlerFunc) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        Handle(ctx, materializer, handler)
    }
}

// Handles a `gin.Context` the `ginx` way.
//
// ctx              The base `gin.Context` object.
// materializer     The logic of deserializing request and serializing response.
// handler          The function handling the exchange and returning either an
//                  exchange `Conclusion` or an error.
func Handle(ctx *gin.Context, materializer ContentMaterializer, handler HandlerFunc) {
    nego := negotiator.New(ctx.Request.Header)
    exchange := &exchange{
        context:      ctx,
        negotiator:   nego,
        materializer: materializer,
    }
    if ok := ensureRequestContentTypeIsSupported(ctx, exchange, materializer); !ok {
        return
    }
    if ok := ensureAcceptedResponseTypeIsSupported(ctx, exchange, materializer); !ok {
        return
    }
    handleExchangeAndWriteResponse(ctx, exchange, materializer, handler)
}

// Ensure request's content type is supported.
func ensureRequestContentTypeIsSupported(ctx *gin.Context, exchange *exchange, materializer ContentMaterializer) bool {
    acceptedConsumedType := ""
    contentType := ctx.GetHeader(headers.ContentType)
    if contentType == "" {
        exchange.consumedType = materializer.ConsumedTypes()[0]
        return true
    }

    // a `content-type` header is provided, it needs to be supported by the negotiator.
    acceptedConsumedType = findFirstSupportedMediaType(contentType, materializer.ConsumedTypes())
    if acceptedConsumedType != "" {
        exchange.consumedType = acceptedConsumedType
        return true
    }

    answerWithProblem(ctx, Problem{
        Status: http.StatusNotAcceptable,
        Title:  fmt.Sprintf("The request type is not supported (header '%s')", headers.ContentType),
        Detail: fmt.Sprintf("A request of type '%s' is not supported.", contentType),
    })
    return false
}

func findFirstSupportedMediaType(contentType string, accepted []string) string {
    requestedMediaType := parseMediaType(contentType)
    if requestedMediaType == nil {
        return ""
    }
    if requestedMediaType.Type == "*" && requestedMediaType.Subtype == "*" {
        return accepted[0]
    }
    for _, a := range accepted {
        acceptedMediaType := parseMediaType(a)
        typeIsCompatible := (acceptedMediaType.Type == requestedMediaType.Type) ||
            acceptedMediaType.Type == "*" ||
            requestedMediaType.Type == "*"
        subTypeIsCompatible := (acceptedMediaType.Subtype == requestedMediaType.Subtype) ||
            acceptedMediaType.Subtype == "*" ||
            requestedMediaType.Subtype == "*"
        if typeIsCompatible && subTypeIsCompatible {
            return a
        }
    }
    return ""
}

type MediaType struct {
    Type       string
    Subtype    string
    Parameters string
}

func parseMediaType(source string) *MediaType {
    if source == "" {
        return nil
    }
    result := MediaType{}
    sourceParts := strings.SplitN(source, ";", 2)
    if len(sourceParts) > 1 {
        result.Parameters = sourceParts[1]
    }
    mediaTypeParts := strings.Split(sourceParts[0], "/")
    if len(mediaTypeParts) != 2 {
        return nil
    }
    result.Type = mediaTypeParts[0]
    result.Subtype = mediaTypeParts[1]
    return &result
}

// Ensure accepted type (response type) is supported.
func ensureAcceptedResponseTypeIsSupported(ctx *gin.Context, exchange Exchange, materializer ContentMaterializer) bool {
    acceptedProducedType := ""
    accept := ctx.GetHeader(headers.Accept)
    if accept != "" {
        // if a `accept` header is provided, it needs to be supported by the negotiator.
        acceptedProducedType = exchange.Negotiator().Type(materializer.ProducedTypes()...)
        if acceptedProducedType == "" {
            answerWithProblem(ctx, Problem{
                Status: http.StatusUnsupportedMediaType,
                Title:  fmt.Sprintf("The accepted types are not supported (header '%s')", headers.Accept),
                Detail: fmt.Sprintf("A request that accepts type(s) '%s' is not supported.", accept),
            })
            return false
        }
    }
    return true
}

// Handle the exchange (call client code) and write the response.
func handleExchangeAndWriteResponse(ctx *gin.Context, exchange Exchange, materializer ContentMaterializer, handler HandlerFunc) {
    conclusion, err := handler(exchange)
    if err != nil {
        answerWithProblem(ctx, err)
        return
    }
    if conclusion == nil {
        // NOOP. The exchange was managed within the handler.
        return
    }
    // Write the conclusion to the response.
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
