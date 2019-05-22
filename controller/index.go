package controller

import (
	"gin_blog/models"
	"github.com/gin-gonic/gin"
	"gin_blog/system"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
	"strconv"
)

func IndexGet(c *gin.Context) {
	var (
		pageIndex int
		pageSize = system.ConfigConent.Page.PageSize
		total int
		page string
		err error
		posts []*models.Post
		policy *bluemonday.Policy
	)
	page = c.Query("page")
	pageIndex, _ = strconv.Atoi(page)
	if pageIndex <=0 {
		pageIndex = 1
	}
	posts, err = models.ListPublishedPost("", pageIndex, pageSize)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	total, err = models.CountPostByTag("")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	policy = bluemonday.StrictPolicy()
	for _, post := range posts {
		post.Tags, _ = models.ListTagByPostId(strconv.FormatUint(uint64(post.ID), 10))
		post.Body = policy.Sanitize(string(blackfriday.Run([]byte(post.Body))))
	}
	//user, _  = c.Get(CONTEXT_USER_KEY)
	//c.HTML(http.StatusOK, "index/index.html", gin.H{
	//
	//})

}
