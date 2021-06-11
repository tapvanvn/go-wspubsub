package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapvanvn/gorouter/v2"
)

//Data type response
type Data struct {
	Success   bool        `json:"success"`
	ErrorCode int         `json:"error_code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"response,omitempty"`
}

//NotFound response not found
func NotFound(context *gorouter.RouteContext) {

	context.Handled = true
	context.W.WriteHeader(http.StatusNotFound)
	fmt.Println(context.R.URL.Path, 404)
}

//ServerError response internal server error
func ServerError(context *gorouter.RouteContext, errorCode int, message string, data interface{}) {

	context.W.WriteHeader(http.StatusInternalServerError)

	responseData := &Data{Success: false,
		ErrorCode: errorCode,
		Message:   message,
		Data:      data}

	if data, err := json.Marshal(responseData); err == nil {

		context.W.Write(data)
	}
	context.Handled = true

	fmt.Println(context.R.URL.Path, 500)
}

//InvalidParam response invalid param
func InvalidParam(context *gorouter.RouteContext, errorCode int, message string, data interface{}) {

	context.W.WriteHeader(http.StatusBadRequest)

	responseData := &Data{Success: false,
		ErrorCode: errorCode,
		Message:   message,
		Data:      data}

	if data, err := json.Marshal(responseData); err == nil {
		fmt.Println(string(data))
		context.W.Write(data)
	}
	context.Handled = true
	fmt.Println("message", responseData.Message)
	fmt.Println(context.R.URL.Path, 400)
}

//BadRequest response for bad request
func BadRequest(context *gorouter.RouteContext, errorCode int, message string, data interface{}) {

	context.W.WriteHeader(http.StatusBadRequest)

	responseData := &Data{Success: false,
		ErrorCode: errorCode,
		Message:   message,
		Data:      data}

	if data, err := json.Marshal(responseData); err == nil {

		fmt.Println(string(data))
		context.W.Write(data)
	}
	context.Handled = true
	fmt.Println("message", responseData.Message)
	fmt.Println(context.R.URL.Path, 400)
}

//Conflict response a conflict
func Conflict(context *gorouter.RouteContext, errorCode int, message string, data interface{}) {

	context.W.WriteHeader(http.StatusConflict)

	responseData := &Data{Success: false,
		ErrorCode: errorCode,
		Message:   message,
		Data:      data}

	if data, err := json.Marshal(responseData); err == nil {

		context.W.Write(data)
	}
	context.Handled = true

	fmt.Println(context.Path, 409)
}

//Unauthorize response an unauthorize response
func Unauthorize(context *gorouter.RouteContext, errorCode int, message string, data interface{}) {

	context.W.WriteHeader(http.StatusUnauthorized)

	responseData := &Data{Success: false,
		ErrorCode: errorCode,
		Message:   message,
		Data:      data}

	if data, err := json.Marshal(responseData); err == nil {

		context.W.Write(data)
	}
	context.Handled = true

	fmt.Println(context.R.URL.Path, 401)
}

//Success response success
func Success(context *gorouter.RouteContext, data interface{}) {

	responseData := &Data{Success: true,
		Data: data}

	if data, err := json.Marshal(responseData); err == nil {

		context.W.Write(data)
	}
	context.Handled = true

	fmt.Println(context.R.URL.Path, 200)

}

//SendFile send file to client
func SendFile(context *gorouter.RouteContext, fileName string, content *[]byte) {

	context.W.Header().Set("Content-Type", "application/octet-stream")
	context.W.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	context.W.Header().Set("Content-Transfer-Encoding", "binary")
	context.W.Header().Set("Expires", "0")
	context.W.Write(*content)
	context.Handled = true

}
