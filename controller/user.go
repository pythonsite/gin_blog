package controller

import (
	"gin_blog/models"
	"gin_blog/utils"
	"github.com/astaxie/beego/logs"
	"strconv"

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

func UserIndex(c *gin.Context) {
	users, _ := models.ListUsers()
	logs.Info(users)
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/user.html", gin.H{
		"users": users,
		"user": user,
		"comments": models.MustListUnreadComment(),
	})
}

func UserLock(c *gin.Context) {
	var (
		err error
		_id uint64
		res = gin.H{}
		user *models.User
	)
	defer writeJson(c, res)
	id := c.Param("id")
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"]  = err.Error()
		return
	}
	user, err = models.GetUser(uint(_id))
	if err != nil {
		res["message"] = err.Error()
		return
	}
	user.LockState = !user.LockState
	err = user.Lock()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func UserDelete(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJson(c, res)
	id := c.Param("id")
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	user := &models.User{}
	user.ID = uint(uid)
	err = user.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}