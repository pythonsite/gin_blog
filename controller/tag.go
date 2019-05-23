package controller

import (
	"gin_blog/models"
	"github.com/gin-gonic/gin"
)

func TagCreate(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJson(c, res)
	name := c.PostForm("value")
	tag:= &models.Tag{Name:name}
	err = tag.Insert()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
	res["data"] = tag
}