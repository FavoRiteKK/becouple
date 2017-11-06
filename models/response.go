package models

// json fields must have upper-case first letter, to be visible for marshal package to encode/decode

type Data map[string]interface{}

type ServerResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Err     string                 `json:"error,omitempty"`
	ErrCode int                    `json:"ecode,omitempty"`
}

type AuthResponse struct {
	Jwt string `json:"token,omitempty"`
}
