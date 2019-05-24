package controller

import (
	"gin_blog/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func SubscriberIndex(c *gin.Context) {
	subscribers, _ := models.ListSubscriber(false)
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/subscriber.html",gin.H{
		"subscribers": subscribers,
		"user": user,
		"comments": models.MustListUnreadComment(),
	})
}

func sendEmailToSubscribers(subject, body string) (err error) {
	var (
		subscribers []*models.Subscriber
		emails = make([]string,0)
	)
	subscribers,err = models.ListSubscriber(true)
	if err != nil {
		return
	}
	for _, subscriber := range subscribers {
		emails = append(emails, subscriber.Email)
	}
	if len(emails) == 0 {
		err = errors.New("no subscribers!")
		return
	}
	err = sendMail(strings.Join(emails, ";"), subject, body)
	return
}

func SubscriberPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJson(c, res)
	email := c.PostForm("mail")
	subject := c.PostForm("subject")
	body := c.PostForm("body")
	if len(email) > 0 {
		err = sendMail(email, subject, body)
	} else {
		err = sendEmailToSubscribers(subject, body)
	}
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func SbuscribeGet(c *gin.Context) {
	count, _ := models.CountSubscriber()
	c.HTML(http.StatusOK, "other/subscribe.html",gin.H{
		"total": count,
	})
}
