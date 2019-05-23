package controller

import (
	"gin_blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PageIndex(c *gin.Context) {
	pages, _ := models.ListAllPage()
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/page.html", gin.H{
		"pages": pages,
		"user": user,
		"comments": models.MustListUnreadComment(),
	})
}

func PageNew(c *gin.Context) {
	c.HTML(http.StatusOK, "page/new.html", nil)
}

func PageCreate(c *gin.Context) {
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished
	page := &models.Page{
		Title: title,
		Body: body,
		IsPublished:published,
	}
	err := page.Insert()
	if err != nil {
		c.HTML(http.StatusOK, "page/new.html", gin.H{
			"message": err.Error(),
			"page": page,
		})
	}
}
