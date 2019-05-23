package controller

import (
	"gin_blog/models"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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

func PostNew(c *gin.Context) {
	c.HTML(http.StatusOK, "post/new.html", nil)
}

func PostCreate(c *gin.Context) {
	tags := c.PostForm("tags")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished

	post := &models.Post{
		Title: title,
		Body: body,
		IsPublished: published,
	}
	err := post.Insert()
	if err != nil {
		c.HTML(http.StatusOK, "post/new.html", gin.H{
			"post": post,
			"message": err.Error(),
		})
		return
	}
	if len(tags) > 0 {
		tagArr := strings.Split(tags, ",")
		for _, tag := range tagArr {
			tagId, err := strconv.ParseUint(tag, 10, 64)
			if err != nil {
				continue
			}
			pt := &models.PostTag{
				PostId:post.ID,
				TagId:uint(tagId),
			}
			err = pt.Insert()
			logs.Error("insert db error:%v",err)
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/post")
}

func PostEdit(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil {
		Handle404(c)
		return
	}
	post.Tags, _ = models.ListTagByPostId(id)
	c.HTML(http.StatusOK, "post/modify.html", gin.H{
		"post": post,
	})
}

func PostUpdate(c *gin.Context) {
	id := c.Param("id")
	tags := c.PostForm("tags")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished
	pid, err := strconv.ParseUint(id,10,64)
	if err != nil {
		Handle404(c)
		return
	}
	post := &models.Post{
		Title: title,
		Body: body,
		IsPublished:published,
	}
	post.ID = uint(pid)
	err = post.Update()
	if err != nil {
		c.HTML(http.StatusOK, "post/modify.html", gin.H{
			"post": post,
			"message": err.Error(),
		})
		return
	}
	//删除tag
	err = models.DeletePostTagByPostId(post.ID)
	if len(tags) > 0 {
		tagArr := strings.Split(tags, ",")
		for _, tag := range tagArr {
			tagId, err := strconv.ParseUint(tag, 10, 64)
			if err != nil {
				continue
			}
			pt := &models.PostTag{
				PostId:post.ID,
				TagId: uint(tagId),
			}
			_ = pt.Insert()
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/post")
}
