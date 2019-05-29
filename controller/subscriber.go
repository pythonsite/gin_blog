package controller

import (
	"fmt"
	"gin_blog/models"
	"gin_blog/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
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

func SubuscribeGet(c *gin.Context) {
	count, _ := models.CountSubscriber()
	c.HTML(http.StatusOK, "other/subscribe.html",gin.H{
		"total": count,
	})
}

func Subscribe(c *gin.Context) {
	mail := c.PostForm("mail")
	var err error
	if len(mail) > 0 {
		var subscriber *models.Subscriber
		subscriber, err = models.GetSubscriberByEmail(mail)
		if err == nil {
			if !subscriber.VerifyState && utils.GetCurrentTime().After(subscriber.OutTime) {
				err = sendActiveEmail(subscriber)
				if err == nil {
				//	TODO 订阅功能
				}
			}
		}
	}
}

func sendActiveEmail(subscriber *models.Subscriber) (err error) {
	uuid := utils.UUID()
	duration, _ := time.ParseDuration("30m")
	subscriber.OutTime = utils.GetCurrentTime().Add(duration)
	subscriber.SecretKey = uuid
	signature := utils.Md5(subscriber.Email + uuid + subscriber.OutTime.Format("20060102150405"))
	subscriber.Signature = signature
	err = sendMail(subscriber.Email, "[icode blog]邮箱验证",fmt.Sprintf("%s/active?sid=%s", models.ConfigConent.CommonConfig.Domain, signature))
	if err != nil {
		return
	}
	err = subscriber.Update()
	return
}
