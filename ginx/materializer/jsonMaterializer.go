package materializer

import "github.com/adeynack/gin-extensions/ginx"

type JsonContentMaterializer struct {
}

var _ ginx.ContentMaterializer = &JsonContentMaterializer{}

func (m *JsonContentMaterializer) ConsumedTypes() []string {
    return consumedTypes
}

func (m *JsonContentMaterializer) ReadRequestBody(exchange ginx.Exchange, requestBody interface{}) (err error) {
    panic("implement me") // todo

}

func (m *JsonContentMaterializer) ProducedTypes() []string {
    return producedTypes
}

func (m *JsonContentMaterializer) WriteResponseBody(exchange ginx.Exchange, conclusion ginx.Conclusion) (err error) {
    panic("implement me") // todo
}

var consumedTypes = []string{"*/json"}

var producedTypes = []string{"application/json"}
