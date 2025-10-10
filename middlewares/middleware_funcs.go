package middlewares

import (
	"fmt"
	"go-koa/contexts"
	"net/http"
	"time"
)

func LogMw(context *contexts.Context, next Next) error {
	now := time.Now().Format(time.RFC3339)
	uri := context.Request.URL.RequestURI()
	method := context.Request.Method
	fmt.Printf("[%s] START %s %s\n", now, method, uri)

	defer func() {
		now := time.Now().Format(time.RFC3339)
		statusCode := context.GetStatusCode()
		statusText := http.StatusText(statusCode)
		fmt.Printf("[%s] COMPLETE %s %s %d %s\n", now, method, uri, statusCode, statusText)
	}()

	return next()
}

func ErrorMw(context *contexts.Context, next Next) error {
	err := next()
	if err != nil {
		fmt.Println("错误: ", err)
		context.Error(err.Error(), http.StatusInternalServerError)
	}

	return err
}
