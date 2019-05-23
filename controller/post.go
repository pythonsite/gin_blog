package controller

import (
	"gin_blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostIndex(c *gin.Context) {
	posts, _ := models.ListAllPost("")
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/post.html", gin.H{
		"posts": posts,
		"Active": "posts",
		"user": user,
		"comments": models.MustListUnreadComment(),
	})

}