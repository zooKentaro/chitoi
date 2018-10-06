package data

// BaseResponse is XXX
type BaseResponse struct {
	Code         int         `json:"code"`
	DebugMessage string      `json:"debug_message"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}
