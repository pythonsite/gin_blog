package controller

import (
	"gin_blog/models"
	"gin_blog/utils"
	"github.com/astaxie/beego/logs"
	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SigninGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signin.html", nil)
}

func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signup.html", nil)
}

func SignupPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJson(c, res)
	email := c.PostForm("email")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	user := &models.User{
		Email:     email,
		Telephone: telephone,
		Password:  password,
		IsAdmin:   true,
	}
	if len(user.Email) == 0 || len(user.Password) == 0 {
		res["message"] = "email or password cannot be null"
		return
	}
	user.Password = utils.Md5(user.Email + user.Password)
	err = user.Insert()
	if err != nil {
		res["message"] = "email already exists"
		return
	}
	res["succeed"] = true
}

func SigninPost(c *gin.Context) {
	var (
		err  error
		user *models.User
	)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.HTML(http.StatusOK, "auth/signin.html", gin.H{
			"message": "username or password cannot be null",
		})
		return
	}
	user, err = models.GetUserByUsername(username)
	if err != nil || user.Password != utils.Md5(username+password) {
		c.HTML(http.StatusOK, "auth/signin.html", gin.H{
			"message": "invalid username or password",
		})
		return
	}
	if user.LockState {
		c.HTML(http.StatusOK, "auth/signin.html", gin.H{
			"message": "Your account have been locked",
		})
		return
	}
	s :=  sessions.Default(c)
	s.Clear()
	s.Set(SESSION_KEY, user.ID)
	err = s.Save()
	if err != nil {
		logs.Error("session save error:%v",err)
	}
	if user.IsAdmin {
		c.Redirect(http.StatusMovedPermanently, "/admin/index")
	} else {
		c.Redirect(http.StatusMovedPermanently, "/")
	}
}

func LogoutGet(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	_ = s.Save()
	c.Redirect(http.StatusSeeOther, "/signin")
}
