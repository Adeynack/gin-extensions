package ginx

// Responsible for deserializing the request's body and serialize the
// request's response for a specific set of media types.
//
// Example: JSON
//
type ContentMaterializer interface {

    ExchangeHandler

    // Gets the list of media types that can be deserialized from the request.
    ConsumedTypes() []string

    // Read (aka deserialize) the request's body to a struct.
    ReadRequestBody(exchange Exchange, requestBody interface{}) (err error)

    // Gets the list of media types that can be serialized to the response.
    ProducedTypes() []string

    // Writes (aka serialize) the response's body.
    WriteResponseBody(exchange Exchange, conclusion *Conclusion) (err error)
}
