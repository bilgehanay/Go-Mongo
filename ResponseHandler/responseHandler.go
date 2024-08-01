package ResponseHandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
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

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TraceId string      `json:"traceId,omitempty"`
	Data    interface{} `json:"data"`
	Count   int         `json:"count,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

var response map[int]Message

func LoadMessages() error {
	viper.AddConfigPath("response.json")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var data struct {
		Messages []Message `json:"handler"`
	}

	if err := viper.Unmarshal(&data); err != nil {
		return err
	}
	response = make(map[int]Message)
	for _, msg := range data.Messages {
		response[msg.Code] = msg
	}
	fmt.Println(response)
	return nil
}

func HandleError(c *gin.Context, code int, traceId string, data interface{}, errs interface{}) {
	if msg, ok := response[code]; ok {
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
	if msg, ok := response[10000]; ok {
		c.JSON(msg.Status, gin.H{
			"success": true,
			"errors":  nil,
			"code":    msg.Code,
			"message": msg.Message,
			"data":    data,
			"count":   count,
		})
	}
}

func New() *Response {
	return &Response{
		Success: false,
		Code:    0,
		Message: "",
		TraceId: "",
		Data:    nil,
		Errors:  nil,
	}
}

func (r Response) SendError(c *gin.Context, code int) {
	err := response[code]
	r.Message = err.Message
	r.Code = code
	fmt.Println(response)
	c.JSON(err.Status, r)
	return
}

func (r Response) SendSuccess(c *gin.Context) {
	r.Success = true
	r.Message = "OK"
	r.Code = 10000
	c.JSON(http.StatusOK, r)
	return
}
