package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

type Message struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  int    `json:"status"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TraceId string      `json:"traceId,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

var errorMap map[int]Message
var successMessage Message

func LoadMessages(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open message file: %v", err)
	}
	defer file.Close()

	var data struct {
		Errors  []Message `json:"errors"`
		Success Message   `json:"success"`
	}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return fmt.Errorf("could not decode message file: %v", err)
	}

	errorMap = make(map[int]Message)
	for _, msg := range data.Errors {
		errorMap[msg.Status] = msg
	}

	successMessage = data.Success
	return nil
}

func HandleError(c *gin.Context, status int, traceId string, data interface{}, errs interface{}) {
	if msg, ok := errorMap[status]; ok {
		response := ErrorResponse{
			Success: false,
			Code:    msg.Code,
			Message: msg.Message,
			TraceId: traceId,
			Data:    data,
			Errors:  errs,
		}
		c.JSON(msg.Status, response)
	} else {
		c.JSON(500, gin.H{"success": false, "message": "Internal Server Error"})
	}
}

func HandleSuccess(c *gin.Context, data interface{}, count int) {
	c.JSON(successMessage.Status, gin.H{
		"success": true,
		"errors":  nil,
		"code":    successMessage.Code,
		"message": successMessage.Message,
		"data":    data,
		"count":   count,
	})
}
