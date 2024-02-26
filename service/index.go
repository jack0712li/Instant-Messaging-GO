package service


import "github.com/gin-gonic/gin"


// GetIndex
// @Tags index
// @Success 200 {string} Hello World
// @Router /index [get]
func GetIndex(c *gin.Context){

	c.JSON(200, gin.H{
		"message": "Hello World",
	})

}