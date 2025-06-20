package models

import "encoding/json"

type KafkaBookRequest struct {
	Method     string          `json:"method"`
	Type       string          `json:"type"`
	RelationID string          `json:"relation_id"`
	Payload    json.RawMessage `json:"payload"`
}

type KafkaError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type KafkaBookResponse struct {
	Method     string          `json:"method"`
	Type       string          `json:"type"`
	RelationID string          `json:"relation_id"`
	Result     json.RawMessage `json:"result"`
	Error      KafkaError      `json:"error"`
}

type GetUserBooksRequest struct {
	UserID uint   `json:"user_id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Limit  int    `json:"limit"`
}

type GetUserBooksResponse struct {
	Books []Book `json:"books"`
}

type DeleteBook struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
}
