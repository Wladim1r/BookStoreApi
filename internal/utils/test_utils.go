package utils

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func toFloat64(num interface{}) float64 {
	switch v := num.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	default:
		return 0
	}
}

func Success_GetUserBook(t *testing.T, expected map[string]interface{}, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	actualBook := response["data"].(map[string]interface{})

	assert.Equal(t, expected["title"], actualBook["title"],
		"Titles of books do not match")

	assert.Equal(t, expected["author"], actualBook["author"],
		"Authors of books do not match")

	assert.Equal(t, toFloat64(expected["price"]), toFloat64(actualBook["price"]),
		"Prices of books do not match")
}

func Success_GetUserBooks(t *testing.T, expected map[string]interface{}, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response) // gin.H{} -> map[string]interface{}
	assert.NoError(t, err, "Failed to parse JSON")

	expectedMeta := expected["meta"].(map[string]interface{})
	actualMeta := response["meta"].(map[string]interface{})

	assert.Equal(t, toFloat64(expectedMeta["total"]),
		actualMeta["total"], "Total count mismatch")

	assert.Equal(t, toFloat64(expectedMeta["user_id"]),
		actualMeta["user_id"], "user_id mismatch")

	expectedBooks := expected["data"].([]interface{})
	actualBooks := response["data"].([]interface{})

	for i := range expectedBooks {
		expectedBook := expectedBooks[i].(map[string]interface{})
		actualBook := actualBooks[i].(map[string]interface{})

		assert.Equal(t, expectedBook["title"], actualBook["title"],
			"Titles of books do not match in book №", i)

		assert.Equal(t, expectedBook["author"], actualBook["author"],
			"Authors of books do not match in book №", i)

		assert.Equal(t, toFloat64(expectedBook["price"]), toFloat64(actualBook["price"]),
			"Prices of books do not match in book №", i)
	}
}

func WrongID(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedResponse := "Authentication required"

	assert.Equal(t, expectedResponse, response["error"])
}

func ServiceError(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedResponse := "Database connection failed"
	assert.Equal(t, expectedResponse, response["error"])

}

func WrongParamID(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedResponse := "Invalid ID in request"
	assert.Equal(t, expectedResponse, response["error"])
}

func NotFound(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedResponse := "Book with entred ID does not exist"
	assert.Equal(t, expectedResponse, response["error"])
}

func Success_PostBook(t *testing.T, expected map[string]interface{}, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	actualBook := response["Created Book"].(map[string]interface{})

	assert.Equal(t, expected["title"], actualBook["title"],
		"Titles of books do not match")

	assert.Equal(t, expected["author"], actualBook["author"],
		"Authors of books do not match")

	assert.Equal(t, toFloat64(expected["price"]), toFloat64(actualBook["price"]),
		"Prices of books do not match")
}

func InvalidBodyRequest(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedError := "Invalid body request"
	assert.Equal(t, expectedError, response["error"])
}

func Success_UpdateBook(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedMessage := "Alterations have been done"

	assert.Equal(t, expectedMessage, response["message"])
}

func Success_DeleteBook(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedMessage := "Book was successfully deleted"

	assert.Equal(t, expectedMessage, response["message"])
}

func Success_Register(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedMessage := "User created"

	assert.Equal(t, expectedMessage, response["message"])
}

func HashPasswordError(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedMessage := "Failed to hash password"

	assert.Equal(t, expectedMessage, response["error"])
}

func Success_Login(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedPrefix := "Your token: eyJ"

	assert.True(t, strings.HasPrefix(response["message"].(string), expectedPrefix))
}

func NotRegistred(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedMessage := "You have not registered yet"

	assert.Equal(t, expectedMessage, response["error"])
}

func IncorrectPassword(t *testing.T, actual string) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(actual), &response)
	assert.NoError(t, err, "Failed to parse JSON")

	expectedMessage := "Entred incorrect password"

	assert.Equal(t, expectedMessage, response["error"])
}
