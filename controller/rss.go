package controller

import (
	"fmt"
	"gin_blog/models"
	"gin_blog/utils"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

func RssGet(c *gin.Context) {
	now := utils.GetCurrentTime()
	domain := models.ConfigConent.CommonConfig.Domain
	feed := &feeds.Feed{
		Title: "ginblog",
		Link: &feeds.Link{Href:domain},
		Description:"golang,python",
		Author: &feeds.Author{Name:"icode", Email:"hjzhaofan@gmail.com"},
		Created: now,
	}

	feed.Items = make([]*feeds.Item, 0)
	posts, err := models.ListPublishedPost("",0,0)
	if err != nil {
		logs.Error("list published post error:%v",err)
		return
	}
	for _, post := range posts {
		item := &feeds.Item{
			Id: fmt.Sprintf("%s/post/%d", domain, post.ID),
			Title: post.Title,
			Link: &feeds.Link{Href:fmt.Sprintf("%s/post/%d", domain, post.ID)},
			Description:string(post.Excerpt()),
			Created:now,
		}
		feed.Items = append(feed.Items, item)
	}
	rss, err := feed.ToRss()
	if err != nil {
		logs.Error("feed to rss error:%v",err)
		return
	}
	_, err = c.Writer.WriteString(rss)
	if err != nil {
		logs.Error("c writer write string error:%v",err)
	}
}
