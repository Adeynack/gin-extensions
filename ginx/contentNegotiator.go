package ginx

type ContentMaterializer interface {
    ConsumedTypes() []string
    ReadRequestBody(exchange Exchange, requestBody interface{}) (err error)
    ProducedTypes() []string
    WriteResponseBody(exchange Exchange, conclusion Conclusion) (err error)
}
