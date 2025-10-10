package contexts

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	// 由中间件添加的键值对
	Keys map[string]any
	// 动态路由参数
	Params map[string]string
	// 状态码
	statusCode int
}

func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Writer:  writer,
		Request: request,
		Keys:    make(map[string]any),
		Params:  make(map[string]string),
	}
}

func (context *Context) SetStatusCode(code int) {
	context.statusCode = code
	context.Writer.WriteHeader(code)
}

func (context *Context) NotFound() {
	context.statusCode = http.StatusNotFound
	http.NotFound(context.Writer, context.Request)
}

func (context *Context) GetStatusCode() int {
	return context.statusCode
}

func (context *Context) Error(s string, serverError int) {
	context.statusCode = serverError
	http.Error(context.Writer, s, serverError)
}

func (context *Context) ServeFile(filePath string) {
	context.statusCode = http.StatusOK
	http.ServeFile(context.Writer, context.Request, filePath)
}

func (context *Context) WriteString(s string) error {
	return context.Write([]byte(s))
}

func (context *Context) Write(bytes []byte) error {
	if context.statusCode == 0 {
		context.SetStatusCode(http.StatusOK)
	}

	_, err := context.Writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (context *Context) JSON(data any) error {
	context.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	if context.statusCode == 0 {
		context.SetStatusCode(http.StatusOK)
	}

	return json.NewEncoder(context.Writer).Encode(data)
}

func (context *Context) HTML(data string) error {
	context.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	return context.WriteString(data)
}

func (context *Context) Redirect(location string) {
	if context.statusCode == 0 {
		context.SetStatusCode(http.StatusFound)
	}

	http.Redirect(context.Writer, context.Request, location, context.statusCode)
}

func (context *Context) Query(key string) string {
	return context.Request.URL.Query().Get(key)
}

func (context *Context) FormValue(key string) (string, error) {
	err := context.Request.ParseForm()
	if err != nil {
		return "", err
	}

	return context.Request.PostForm.Get(key), nil
}

func (context *Context) JSONBody(receiver *any) error {
	if context.Request.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("Content-Type不是application/json类型")
	}

	return json.NewDecoder(context.Request.Body).Decode(receiver)
}
