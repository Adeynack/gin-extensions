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
    Cause    error  `json:"-"` // for debugging purposes only, never serialize with the JSON output
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
    if p.Cause == nil {
        return p.String()
    } else {
        return fmt.Sprintf("%s caused by: %s", p.String(), p.Error())
    }
}
