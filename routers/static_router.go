package routers

import (
	"go-koa/contexts"
	"go-koa/middlewares"
	"os"
	"path"
)

func route(filePath string) middlewares.Middleware {
	return func(ctx *contexts.Context, next middlewares.Next) error {
		ctx.ServeFile(filePath)
		return nil
	}
}

func Static(prefix string, staticDir string, cors *CORSConfig) (*Router, error) {
	router := NewRouter(prefix)

	var helper func(string, string) error
	helper = func(prefix string, dir string) error {
		// 读取目录
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		// 为每一个文件添加路由
		for _, dirEntry := range dirEntries {
			filePath := path.Join(dir, dirEntry.Name())
			// 递归托管目录
			if dirEntry.IsDir() {
				err := helper(prefix+dirEntry.Name()+"/", filePath)
				if err != nil {
					return err
				}
			}
			// 托管文件
			err := router.Get(prefix+dirEntry.Name(), route(filePath), cors)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := helper("/", staticDir)
	if err != nil {
		return nil, err
	}

	return router, nil
}
