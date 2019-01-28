package materializer

import (
    "io/ioutil"

    "github.com/adeynack/gin-extensions/ginx"
    "github.com/gin-gonic/gin"
    "gopkg.in/yaml.v2"
)

type YamlContentMaterializer struct {
}

var _ ginx.ContentMaterializer = &YamlContentMaterializer{}

func (m YamlContentMaterializer) Handler(handler ginx.HandlerFunc) gin.HandlerFunc {
    return ginx.Handler(m, handler)
}

func (m YamlContentMaterializer) Handle(ctx *gin.Context, handler ginx.HandlerFunc) {
    ginx.Handle(ctx, m, handler)
}

func (m YamlContentMaterializer) ConsumedTypes() []string {
    return yamlConsumedTypes
}

func (m YamlContentMaterializer) ReadRequestBody(exchange ginx.Exchange, requestBody interface{}) (err error) {
    bodyBytes, err := ioutil.ReadAll(exchange.Context().Request.Body)
    if err != nil {
        return err
    }
    err = yaml.Unmarshal(bodyBytes, requestBody)
    return err
}

func (m YamlContentMaterializer) ProducedTypes() []string {
    return yamlProducedTypes
}

func (m YamlContentMaterializer) WriteResponseBody(exchange ginx.Exchange, conclusion *ginx.Conclusion) (err error) {
    var responseBytes []byte
    if conclusion.ResponseBody != nil {
        responseBytes, err = yaml.Marshal(conclusion.ResponseBody)
        if err != nil {
            return err
        }
    } else {
        responseBytes = []byte{}
    }
    exchange.Context().Data(conclusion.Status, yamlProducedTypes[0], responseBytes)
    return nil
}

var yamlConsumedTypes = []string{"*/yaml"}

var yamlProducedTypes = []string{"application/yaml"}
