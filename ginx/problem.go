package ginx

import (
    "fmt"
    "encoding/json"
)

type Problem struct {
    Status   int    `json:"status"`
    Type     string `json:"type"`
    Title    string `json:"title"`
    Detail   string `json:"detail"`
    Instance string `json:"instance"`
}

var _ fmt.Stringer = &Problem{}
var _ error = &Problem{}

func (p Problem) String() string {
    bytes, err := json.Marshal(p)
    if err != nil {
        panic(err)
    }
    return string(bytes)
}

func (p Problem) Error() string {
    return p.String()
}
