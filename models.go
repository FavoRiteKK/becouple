package main

// json fields must have upper-case first letter, to be visible for marshal package to encode/decode

type ServerResponse struct {
    Success bool    `json:"success"`
    Err string      `json:"err"`
}

type AuthResponse struct {
    Jwt string      `json:"token"`
    ServerResponse
}