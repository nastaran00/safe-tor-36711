package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RenderHome(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Go Gin Boiler Plate",
	})
}

func Welcome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Server started successfully at" + time.Now().String(),
	})
}
func Signuplogin(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"loginSignup.html", //"new.html", //
		gin.H{
			"title": "login Page",
		})
	//c.JSON(200, gin.H{
	//	"message": "success",
	//})
}
func Homelogged(c *gin.Context) {
	//if log in:
	c.HTML(
		http.StatusOK,
		"index_logged.html",
		gin.H{
			"title": "Home",
		})
	//c.JSON(200, gin.H{
	//	"message": "success",
	//})
}
func Recommendersystem(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"recommenderPage.html",
		gin.H{
			"title": "rec Page",
		})
}
func Signuppage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"signup.html",
		gin.H{
			"title": "login Page",
		})
}
func Profile(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"profile.html",
		gin.H{
			"title": "login Page",
		})
}
func Addpeople(c *gin.Context) {
	someUser := &userAccount{}
	err := json.NewDecoder(c.Request.Body).Decode(someUser) //decode the request body into struct and failed if any error occur
	if err != nil {
		respond(c.Writer, message(false, "Invalid request"))
		return
	}

	resp := someUser.Create() //Create account
	//respond(c.Writer, resp)
	c.JSON(200, gin.H{
		"message": "Server started successfully at" + time.Now().String(),
		"res":     resp,
	})
}

func Authenticate(c *gin.Context) {
	someUser := &userAccount{}
	err := json.NewDecoder(c.Request.Body).Decode(someUser)
	if err != nil {
		respond(c.Writer, message(false, "Invalid request"))
		return
	}

	resp := login(someUser.Email, someUser.Password)
	//respond(w, resp)
	//login(someUser.Email, someUser.Password)
	c.JSON(
		http.StatusOK,
		gin.H{
			"resp":  resp,
			"title": "signed Page",
		})
}
func Editing(c *gin.Context) {
	someUser := &userAccount{}
	err := json.NewDecoder(c.Request.Body).Decode(someUser)
	if err != nil {
		respond(c.Writer, message(false, "Invalid request"))
		return
	}

	resp := profile_edit(someUser.Email, someUser.Password)
	//respond(w, resp)
	//login(someUser.Email, someUser.Password)
	c.JSON(
		http.StatusOK,
		gin.H{
			"resp":  resp,
			"title": "signed Page",
		})
}

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	someUser := &userAccount{}
	err := json.NewDecoder(r.Body).Decode(someUser) //decode the request body into struct and failed if any error occur
	if err != nil {
		respond(w, message(false, "Invalid request"))
		return
	}

	resp := someUser.Create() //Create account
	respond(w, resp)
}

func Logout(c *gin.Context) {
	resp := message(true, "Success")
	respond(c.Writer, resp)
}

func QuoteResponse(c *gin.Context) {
	resp := message(true, "Success")
	resp["data"] = "new new new new enw new newn ewjdkfh skhfd"
	respond(c.Writer, resp)
}
