package gremgo

import (
	"fmt"
)

// GremlinNetworkError - This returns decorated error with useful information
type GremlinNetworkError struct {
	Attributes interface{} `json:"attributes" omitempty`
	Code int32 `json:"code" omitempty`
	Message string `json:"message" omitempty`
	ConnStr string `json:"conn_str omitempty`
}

func (e GremlinNetworkError) Error() string {
	return fmt.Sprint("connStr:", e.ConnStr, " ; code: ", e.Code, " ; message:", e.Message)
}
