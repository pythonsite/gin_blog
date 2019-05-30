package main

import (
	"gin_blog/controller"
	"gin_blog/models"
	"gin_blog/utils"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	setTemplate(router)
	setSessions(router)
	router.Use(SharedData())

	router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))
	router.GET("/", controller.IndexGet)
	router.GET("/index", controller.IndexGet)
	router.GET("/rss", controller.RssGet)

	router.GET("/signin", controller.SigninGet)
	router.POST("/signin", controller.SigninPost)

	router.GET("/logout", controller.LogoutGet)

	router.GET("/signup", controller.SignupGet)
	router.POST("/signup", controller.SignupPost)

	router.GET("/post/:id", controller.PostGet)
	router.GET("/link/:id", controller.LinkGet)
	router.GET("/page/:id",controller.PageGet)
	router.GET("/tag/:tag", controller.TagGet)
	router.GET("/archives/:year/:month", controller.ArchiveGet)

	router.GET("/subscribe", controller.SubuscribeGet)
	router.POST("/subscribe", controller.Subscribe)

	visitor := router.Group("/visitor")
	visitor.Use(AuthRequired())
	{
		visitor.POST("/new_comment", controller.CommentPost)
	}


	authorized := router.Group("/admin")
	authorized.Use(AdminScopeRequired())
	{
		authorized.GET("/index", controller.AdminIndex)
		// page
		authorized.GET("/page", controller.PageIndex)
		authorized.GET("/new_page", controller.PageNew)
		authorized.POST("/new_page", controller.PageCreate)
		authorized.GET("/page/:id/edit", controller.PageEdit)
		authorized.POST("/page/:id/edit", controller.PageUpdate)
		authorized.POST("/page/:id/publish", controller.PagePublish)
		authorized.POST("/page/:id/delete", controller.PageDelete)

		// post
		authorized.GET("/post", controller.PostIndex)
		authorized.GET("/new_post", controller.PostNew)
		authorized.POST("/new_post",controller.PostCreate)
		authorized.GET("/post/:id/edit", controller.PostEdit)
		authorized.POST("/post/:id/edit", controller.PostUpdate)
		authorized.POST("/post/:id/delete", controller.PostDelete)

		// tag
		authorized.POST("/new_tag", controller.TagCreate)

		authorized.GET("/user", controller.UserIndex)
		authorized.POST("/user/:id/lock", controller.UserLock)
		authorized.POST("/user/:id/delete", controller.UserDelete)

		// subscriber
		authorized.GET("/subscriber", controller.SubscriberIndex)
		authorized.POST("/subscriber", controller.SubscriberPost)

		// link
		authorized.GET("/link", controller.LinkIndex)
		authorized.POST("/new_link", controller.LinkCreate)
		authorized.POST("/link/:id/edit", controller.LinkUpdate)
		authorized.POST("/link/:id/delete", controller.LinkDelete)

	}
	err := router.Run(":8090")
	if err != nil {
		logs.Error("router run errror:%s",err)
	}
}

func setTemplate(engine *gin.Engine) {
	funcMap := template.FuncMap{
		"dateFormat": utils.DateFormat,
		"substring": utils.Substring,
		"isOdd":      utils.IsOdd,
		"isEven":     utils.IsEven,
		"truncate":   utils.Truncate,
		"add":        utils.Add,
		"minus":      utils.Minus,
		"listtag":    models.GetTagsStr,
	}
	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob(filepath.Join(getCurrentDirectory(),"./views/**/*"))
}

func getCurrentDirectory() string{
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.Error(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func setSessions(router *gin.Engine) {
	store := sessions.NewCookieStore([]byte("ginblog"))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))
}

func AdminScopeRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get(controller.CONTEXT_USER_KEY);user != nil {
			if u, ok := user.(*models.User);ok && u.IsAdmin {
				c.Next()
				return
			}
		}
		logs.Error("user not authorized to visit %s", c.Request.RequestURI)
		c.HTML(http.StatusForbidden, "errors/error.html", gin.H{
			"message": "Forbiden!",
		})
		c.Abort()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get(controller.CONTEXT_USER_KEY);user != nil {
			if u, ok := user.(*models.User); ok && u.IsAdmin {
				c.Next()
				return
			}
		}
		logs.Error("user not authorized to visit %s", c.Request.RequestURI)
		c.HTML(http.StatusForbidden, "errors/error.html", gin.H{
			"message": "Forbidden!",
		})
		c.Abort()
	}
}

func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if uID := session.Get(controller.SESSION_KEY); uID != nil {
			user, err := models.GetUser(uID)
			if err == nil {
				c.Set(controller.CONTEXT_USER_KEY, user)
			}
		}
		c.Next()
	}
}