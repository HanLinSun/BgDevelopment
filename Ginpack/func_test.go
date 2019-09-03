package Ginpack
//http服务器
import (
	"testing"
	"github.com/gin-gonic/gin"
	"net/http"

)
type Login struct{
	Name string `json:"name"`
	Password string `json:"Password"`
}
func GetInfos(c *gin.Context){
	name:=c.Query("name")
	lastname:=c.DefaultQuery("lastname","默认值")
	c.String(http.StatusOK,"Hello %s %s",name,lastname)
}
func PostInfos(c *gin.Context){
	name:=c.PostForm("name")
	lastname:=c.DefaultPostForm("lastname","默认值")
	c.String(http.StatusOK,"Hello %s %s",name,lastname)
}
func LoginForm(c *gin.Context){
	form:=Login{}
	if c.Bind(&form)==nil{
		if form.Name=="root" && form.Password=="root"{
			c.JSON(200,gin.H{"status":"Successful"})
		} else {
			c.JSON(203,gin.H{"status":"账号或者密码错误"})
		}
	}

}
//GIN 处理http

func Test01(t *testing.T) {
	r := gin.Default()
	r.GET("/ping", GetInfos)
	r.POST("/ping",PostInfos)
	r.POST("/LoginForm",LoginForm)
	r.Run(":8082")
	r.GET("/user/:name/:password",ginHandler)
	r.POST("",printer)
}
func printer(c *gin.Context){
	c.String(http.StatusOK,"Hello %s %s","a","b")
}
func ginHandler(c *gin.Context){
	name:=c.Param("name")
	pwd:=c.Param("password")
	c.String(http.StatusOK,"Hello %s %s",name,pwd)
}
// 