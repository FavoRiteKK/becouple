package models

// json fields must have upper-case first letter, to be visible for marshal package to encode/decode

type ServerResponse struct {
	Success bool   `json:"success"`
	Err     string `json:"err,omitempty"`
}

type AuthResponse struct {
	Jwt string `json:"token,omitempty"`
	*ServerResponse
}
