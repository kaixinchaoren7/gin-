package main

import (
	"net/http" // 引入 net/http 标准库，用于 HTTP 状态码常量

	"github.com/gin-gonic/gin"
)

// 首先在终端执行：go get -u github.com/gin-gonic/gin Go Modules 会自动处理依赖下载
func main() {
	// ---------------1. 环境准备与第一个gin应用
	// 创建一个默认的路由引擎
	// Default() 包含了 Logger 和 Recovery 中间件，方便调试和异常恢复
	r := gin.Default()

	// 2. 定义一个 GET 路由 和对应的处理函数
	// 当访问 /ping 时，会执行后面的匿名函数
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong!",
			"status":  "ok",
		})
	})
	// 4. 启动 HTTP 服务，默认监听在 0.0.0.0:8080
	// 也可以指定端口，例如 r.Run(":9090")
	err := r.Run()
	if err != nil {
		panic("Failed to start Gin server: " + err.Error())
	}

	// -------------------------------------- 二：（核心功能）路由
	r.GET("/a", getting)
	r.POST("/b", posting)
	// r.PUT("/c",putting)
	// r.DELETE("/d",deleting)
	// r.PATCH("/e",patching)
	// r.OPTIONS("/f",optionsing)

	// （1）路由参数 (单个/多个)
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": "user id is " + id,
			"status":  "ok",
		})
	})
	// （2）查询参数 例如 /search?query=gin&page=1
	r.GET("/search", func(c *gin.Context) {
		// 使用 c.Query() 获取查询参数，如果不存在，返回空字符串
		query := c.Query("query")
		// 使用 c.DefaultQuery() 获取查询参数，如果不存在，返回指定的默认值
		page := c.DefaultQuery("page", "1")
		// 也可以用 c.GetQuery()，它返回 (value, ok)
		limit, ok := c.GetQuery("limit")
		if !ok {
			limit = "10" // 设置默认值
		}
		c.JSON(http.StatusOK, gin.H{
			"query": query,
			"page":  page,
			"limit": limit,
		})
	})
	// （3）路由分组
	// 创建一个 /api/v1 的路由组
	v1 := r.Group("/api/v1")
	{ // 可以用大括号增加可读性
		v1.GET("/users", func(c *gin.Context) { /* 获取用户列表 */
			c.JSON(http.StatusOK, gin.H{"users": []string{"user1", "user2"}})
		})
		v1.POST("/users", func(c *gin.Context) { /* 创建用户 */
			c.JSON(http.StatusCreated, gin.H{"message": "User created"})
		})
	}
	// 另一个分组
	v2 := r.Group("/api/v2")
	{
		v2.GET("/products", func(c *gin.Context) { /* ... */ })
	}

	//-----------------------------------------------

	//------------------------------------------(三) 核心功能：请求处理
	// （1）获取表单数据
	r.POST("/login", func(c *gin.Context) {
		userName := c.PostForm("username")
		// password := c.DefaultPostForm("password", "123456")  为nil时给默认值
		c.JSON(http.StatusOK, gin.H{
			"status":   "logged in",
			"username": userName,
		})
	})

	// （2）获取 JSON 数据
	r.POST("/")
	type User struct {
		Username string `json:"username" binding:"required"` // 必填字段
		Password string `json:"password" binding:"required"`
		Age      int    `json:"age" binding:"gte=0,lte=130"` // 年龄大于等于0，小于等于130
	}
	r.POST("/register", func(c *gin.Context) {
		var user User
		// ShouldBindJSON 会尝试将请求体中的 JSON 绑定到 user 结构体
		// 如果 JSON 格式错误或不满足 binding 校验规则，会返回 error
		if err := c.ShouldBindJSON(&user); err != nil {
			// 返回 400 Bad Request 错误，并附带错误信息
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return // 终止处理
		}
		// 绑定成功，可以使用 user 对象了
		c.JSON(http.StatusOK, gin.H{
			"message":  "Registration successful",
			"username": user.Username,
			"age":      user.Age,
		})
	})

	//------------------------------------------------

	//------------------------------------------ 四：（核心功能）响应处理
	//（1）返回 JSON 数据
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": "user id is " + id,
			"status":  "ok",
		})
	})
	// （2）返回 String
	r.GET("hello", func(c *gin.Context) {
		name := c.DefaultQuery("name", "Guest")
		c.String(http.StatusOK, "Hello, %s!", name)
	})
	// (3) 返回 HTML 模板
	//  在 r := gin.Default() 之后
	// 加载 templates 目录下所有 .html 文件
	// r.LoadHTMLGlob("templates/*")
	//  或者指定具体文件
	// r.LoadHTMLFiles("templates/index.html", "templates/login.html")
	//渲染模板
	//处理函数
	r.GET("/index", func(c *gin.Context) {
		// 传递给模板的数据
		data := gin.H{
			"title":   "我的主页",
			"message": "欢迎来到 Gin 的世界!",
		}
		c.HTML(http.StatusOK, "index.html", data)
	})
	//（4）重定向
	r.GET("/redirect", func(c *gin.Context) {
		// 301 永久重定向 或 302 临时重定向
		c.Redirect(http.StatusMovedPermanently, "https://xxx.com/")
	})
	//-----------------------------------------------------------

}

func getting(c *gin.Context) { /* ... */ }
func posting(c *gin.Context) { /* ... */ }
