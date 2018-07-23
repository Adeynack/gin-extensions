package ginx

import "github.com/gin-gonic/gin"

type ExchangeHandler interface {

    // Handle a Gin request through a specific GinX handler function.
    Handle(ctx *gin.Context, handler HandlerFunc)

    // Creates a Gin handler function, wrapping a specific GinX handler function.
    Handler(handler HandlerFunc) gin.HandlerFunc

}
