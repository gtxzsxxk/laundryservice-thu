package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.Delims("+-[]", "+-[]")
	r.LoadHTMLFiles("./index.html")
	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.StaticFile("/favicon512.png", "./favicon512.png")
	r.StaticFile("/manifest.json", "./manifest.json")
	r.StaticFile("/sw.js", "./sw.js")
	r.GET("/", index)
	r.GET("/dormitories", get_dorms)
	r.GET("/dormitories/:id", get_dorms_devices)
	r.Run("localhost:10080") // 监听并在 0.0.0.0:8080 上启动服务
}
