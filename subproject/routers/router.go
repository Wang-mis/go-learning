package routers

import (
	"go-koa/contexts"
	"go-koa/middlewares"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

func NewCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}
}

// Router 路由器结构体，用于管理路由规则
// prefix: 路由前缀
// routes: 路由表，按 HTTP 方法分类存储
type Router struct {
	prefix string
	trie   *RouterTrieNode
}

// NewRouter 创建一个新的路由器实例
// 参数 prefix 表示路由前缀
func NewRouter(prefix string) *Router {
	root, _ := newRouterTrie(prefix)

	return &Router{
		prefix: prefix,
		trie:   root,
	}
}

// AddRoute 添加一条路由规则
// 参数 method 表示 HTTP 请求方法
// 参数 path 表示路由路径
// 参数 handler 表示对应的处理函数
func (router *Router) AddRoute(
	method string,
	path string,
	handler middlewares.Middleware,
	cors *CORSConfig,
) error {
	path = router.prefix + path

	leaf, err := router.trie.addPath(path)
	if err != nil {
		return err
	}

	err = leaf.addMethod(method, handler, cors)
	if err != nil {
		return err
	}

	return nil
}

// Get 注册 GET 请求的路由规则
// 参数 path 表示路由路径
// 参数 handler 表示处理函数
func (router *Router) Get(
	path string,
	handler middlewares.Middleware,
	cors *CORSConfig,
) error {
	return router.AddRoute(http.MethodGet, path, handler, cors)
}

// Post 注册 POST 请求的路由规则
// 参数 path 表示路由路径
// 参数 handler 表示处理函数
func (router *Router) Post(
	path string,
	handler middlewares.Middleware,
	cors *CORSConfig,
) error {
	return router.AddRoute(http.MethodPost, path, handler, cors)
}

func setCORS(config *CORSConfig, context *contexts.Context, origin string, method string) {
	if slices.Contains(config.AllowedOrigins, "*") {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	} else if slices.Contains(config.AllowedHeaders, origin) {
		context.Writer.Header().Set("Access-Control-Allow-Headers", origin)
	}

	context.Writer.Header().Set("Access-Control-Allow-Methods", method)
	context.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
	context.Writer.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(config.AllowCredentials))
}

// Routes 返回一个中间件，用于匹配并执行对应的路由处理函数
func (router *Router) Routes() middlewares.Middleware {
	return func(context *contexts.Context, next middlewares.Next) error {
		method := context.Request.Method
		path := context.Request.URL.Path
		// 匹配路由
		leaf, params := router.trie.findPath(path)
		// 未匹配到路由，继续执行下一个中间件
		if leaf == nil {
			return next()
		}
		// 匹配到路由
		if method == http.MethodOptions {
			// 处理预检请求
			requestMethod := context.Request.Header.Get("Access-Control-Request-Method")
			origin := context.Request.Header.Get("Origin")
			_, exist := leaf.handlers[requestMethod]
			if exist {
				config := leaf.corsConfigs[method]
				setCORS(config, context, origin, requestMethod)
				context.SetStatusCode(http.StatusNoContent)
			}

			return next()
		}
		// 处理真正的请求
		handler, exist := leaf.handlers[method]
		if !exist {
			return next()
		}
		// 设置cors配置
		config := leaf.corsConfigs[method]
		origin := context.Request.Header.Get("Origin")
		setCORS(config, context, origin, method)

		context.Params = params
		err := handler(context, next)
		if err != nil {
			return err
		}

		return next()
	}
}
