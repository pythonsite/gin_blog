package controller

import (
	"gin_blog/models"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CommentPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
		post *models.Post
	)
	defer writeJson(c, res)
	s := sessions.Default(c)
	sessionUserID := s.Get(SESSION_KEY)
	userId, _ := sessionUserID.(uint)
	postId := c.PostForm("postId")
	content := c.PostForm("content")
	if len(content) == 0 {
		res["message"] = "content cannot be empty."
		return
	}
	post, err = models.GetPostById(postId)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	pid, err := strconv.ParseUint(postId, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	comment := &models.Comment{
		PostID:uint(pid),
		Content: content,
		UserID:userId,
	}
	err = comment.Insert()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	logs.Info("%s-%s-%s",post.ID, post.Title, content)
	res["succeed"] = true
}

func CommentRead(c *gin.Context) {
	var (
		id string
		_id uint64
		err error
		res = gin.H{}
	)
	defer writeJson(c, res)
	id = c.Param("id")
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	comment := new(models.Comment)
	comment.ID = uint(_id)
	err = comment.Update()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func CommentReadAll(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJson(c, res)
	err = models.SetAllCommentRead()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
