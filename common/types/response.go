package types

type Response struct {
    Data    any    `json:"data"`
    Message string `string:"message"`
    Meta any       `json:"meta"`
}