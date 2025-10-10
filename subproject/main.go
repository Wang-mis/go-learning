package main

import (
	"fmt"
	"go-koa/app"
	"go-koa/contexts"
	"go-koa/middlewares"
	"go-koa/routers"
	"log"
)

func hello(ctx *contexts.Context, _ middlewares.Next) error {
	user := ctx.Params["user"]
	info := ctx.Params["info"]

	return ctx.HTML(
		fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Hello</title>
		</head>
		<body>
		<h2>Hello %s! %s.</h2>
		</body>
		</html>
		`, user, info),
	)
}

func main() {
	webapp := app.NewApp()

	staticRouter, err := routers.Static("/public", "./public", routers.NewCORSConfig())
	if err != nil {
		log.Fatal(err)
	}
	webapp.Use(middlewares.LogMw)
	webapp.Use(middlewares.ErrorMw)

	router := routers.NewRouter("/data/:user/:info")

	err = router.Get("/hello", hello, routers.NewCORSConfig())
	if err != nil {
		log.Fatal(err)
	}

	webapp.Use(staticRouter.Routes())
	webapp.Use(router.Routes())

	err = webapp.Listen(":8000")
	if err != nil {
		fmt.Println("开启监听错误：", err)
		return
	}
}
