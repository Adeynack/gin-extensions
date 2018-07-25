package materializer

import (
    "encoding/json"
    "io/ioutil"

    "github.com/adeynack/gin-extensions/ginx"
    "github.com/gin-gonic/gin"
)

type JsonContentMaterializer struct {
}

var _ ginx.ContentMaterializer = &JsonContentMaterializer{}

func (m *JsonContentMaterializer) Handler(handler ginx.HandlerFunc) gin.HandlerFunc {
    return ginx.Handler(m, handler)
}

func (m *JsonContentMaterializer) Handle(ctx *gin.Context, handler ginx.HandlerFunc) {
    ginx.Handle(ctx, m, handler)
}

func (m *JsonContentMaterializer) ConsumedTypes() []string {
    return consumedTypes
}

func (m *JsonContentMaterializer) ReadRequestBody(exchange ginx.Exchange, requestBody interface{}) (err error) {
    bodyBytes, err := ioutil.ReadAll(exchange.Context().Request.Body)
    if err != nil {
        return err
    }
    err = json.Unmarshal(bodyBytes, requestBody)
    return err
}

func (m *JsonContentMaterializer) ProducedTypes() []string {
    return producedTypes
}

func (m *JsonContentMaterializer) WriteResponseBody(exchange ginx.Exchange, conclusion *ginx.Conclusion) (err error) {
    var responseJsonBytes []byte
    if conclusion.ResponseBody != nil {
        responseJsonBytes, err = json.Marshal(conclusion.ResponseBody)
        if err != nil {
            return err
        }
    } else {
        responseJsonBytes = []byte{}
    }
    exchange.Context().Data(conclusion.Status, producedTypes[0], responseJsonBytes)
    return nil
}

var consumedTypes = []string{"*/json"}

var producedTypes = []string{"application/json"}
