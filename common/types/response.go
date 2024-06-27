package types

type Response struct {
    Data    any    `json:"data"`
    Message string `json:"message"`
    Meta any       `json:"meta"`
}