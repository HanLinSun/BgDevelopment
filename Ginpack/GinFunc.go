package Ginpack

import "C"
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"CountSystem/MysqlLink"
	"strconv"
	"fmt"

	"CountSystem/RedisLink"
)

func InitRouter() *gin.Engine {
	router:=gin.Default()
	router.GET("/",IndexApi)
	router.GET("/mysql",GetUsersInDBApi)
	router.GET("/mysql/:id",GetUserInDBApi)
	//router.GET("/redis",PurchaseApi)
	return router
}
func IndexApi(c *gin.Context){
	c.String(http.StatusOK,"Begin to work")
}
func GetUsersInDBApi(c *gin.Context){
	var u MysqlLink.User
	users,err:=u.FetchAlluser(MysqlLink.SqlDB)
	MysqlLink.CheckErr(err)
	c.JSON(http.StatusOK,gin.H{
		"Users":users,
	})
}
func GetUserInDBApi(c *gin.Context){
	cID:=c.Param("id")//字符串ID
	RealID,err:=strconv.Atoi(cID)//int ID
	if err!=nil{
		fmt.Println("failure:",err.Error())
	}
	var user MysqlLink.User
	user.UserID=RealID
	temp,err:=user.FetchUserInfo(MysqlLink.SqlDB)
	if err!=nil{
		fmt.Println("failure:",err.Error())
	}
	c.JSON(http.StatusOK,gin.H{
		"user_ID":temp.UserID,
		"user_Repo":temp.UserRepo,
	})
}
func RedisInsertApi(c *gin.Context){
	cID:=c.Param("id")//
	RealID,err:=strconv.Atoi(cID)
	if err!=nil{
		fmt.Println("failure:",err.Error())
	}
	var user MysqlLink.User
	user.UserID=RealID
	temp,err:=user.FetchUserInfo(MysqlLink.SqlDB)
	RedisLink.InsertIntoRedis(RedisLink.Pool.Get(),&temp)
}
func PurchaseApi(c *gin.Context){
	//var obj =MysqlLink.Object{1100}

}
func ShowMysqlDataApi(c *gin.Context){

}
func ShowRedisDataApi(c *gin.Context){

}