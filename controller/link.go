package controller

import (
	"gin_blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func LinkIndex(c *gin.Context) {
	links, _ := models.ListLinks()
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/link.html",gin.H{
		"links": links,
		"user":user,
		"comments": models.MustListUnreadComment(),
	})
}

func LinkCreate(c *gin.Context) {
	var (
		err error
		res = gin.H{}
		_sort int64
	)
	defer writeJson(c, res)
	name := c.PostForm("name")
	url := c.PostForm("url")
	sort := c.PostForm("sort")
	if len(name) == 0 || len(url) == 0 {
		res["message"] = "error parameter"
		return
	}
	_sort, err = strconv.ParseInt(sort, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	link := &models.Link{
		Name:name,
		Url:url,
		Sort:int(_sort),
	}
	err = link.Insert()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func LinkUpdate(c *gin.Context) {
	var (
		_id uint64
		_sort int64
		err error
		res = gin.H{}
	)
	defer writeJson(c, res)
	id := c.Param("id")
	name := c.PostForm("name")
	url := c.PostForm("url")
	sort := c.PostForm("sort")
	if len(id) == 0 || len(name) == 0 || len(url) == 0 {
		res["message"] = "error parameter"
		return
	}
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	_sort, err = strconv.ParseInt(sort,10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	link := models.Link{
		Name:name,
		Url:url,
		Sort:int(_sort),
	}
	link.ID = uint(_id)
	err = link.Update()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func LinkDelete(c *gin.Context) {
	var (
		err error
		_id uint64
		res = gin.H{}
	)
	defer writeJson(c, res)
	id := c.Param("id")
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	link := new(models.Link)
	link.ID = uint(_id)
	err = link.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}