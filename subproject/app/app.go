package app

import (
	"go-koa/contexts"
	"go-koa/middlewares"
	"net/http"
)

// App 表示框架的核心结构，负责存储中间件并调度
type App struct {
	middlewares []middlewares.Middleware // 中间件链（按注册顺序执行）
}

// 实现 http.Handler 接口
// 每次有请求进来，都会调用 ServeHTTP
func (a *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 为每个请求创建一个新的 Context，上下文中包含请求和响应对象
	context := contexts.NewContext(writer, request)

	// 调用 compose 执行中间件链
	_ = a.compose(context)
}

// NewApp 创建一个新的 App 实例
func NewApp() *App {
	return &App{}
}

// Use 注册一个中间件，加入中间件链
func (a *App) Use(mw middlewares.Middleware) {
	a.middlewares = append(a.middlewares, mw)
}

// compose 负责执行中间件链（实现“洋葱模型”调度逻辑）
func (a *App) compose(ctx *contexts.Context) error {
	var dispatch func(int) error

	// dispatch 递归调度中间件
	dispatch = func(index int) error {
		// 如果所有中间件都执行完，还没有设置状态码，则返回 404
		if index >= len(a.middlewares) {
			if ctx.GetStatusCode() == 0 {
				ctx.Error("404 Not Found", http.StatusNotFound)
			}
			return nil
		}

		// 执行当前中间件，并传入一个 next() 函数，执行下一个中间件
		return a.middlewares[index](ctx, func() error {
			return dispatch(index + 1)
		})
	}

	// 从第 0 个中间件开始执行
	return dispatch(0)
}

// Listen 启动 HTTP 服务，监听指定地址
func (a *App) Listen(addr string) error {
	return http.ListenAndServe(addr, a) // App 实现了 http.Handler
}
