package models

// json fields must have upper-case first letter, to be visible for marshal package to encode/decode

type Data map[string]interface{}

type ServerResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Err     string                 `json:"error"`
	ErrCode uint                   `json:"ecode"`
}
