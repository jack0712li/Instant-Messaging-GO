package service

import (
	"text/template"
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetIndex
// @Tags index
// @Success 200 {string} Hello World
// @Router /index [get]
func GetIndex(c *gin.Context){

	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")

	// c.JSON(200, gin.H{
	// 	"message": "welcome",
	// })
}


func ToRegister(c *gin.Context) {
	reg, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	reg.Execute(c.Writer, "register")
}


func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html", 
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/main.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html")
	if err != nil {
		panic(err)
	}

	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	ind.Execute(c.Writer, user)

	// c.JSON(200, gin.H{
	// 	"message": "welcome",
	// })
}