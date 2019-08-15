package MysqlLink
//购买消费同步机制
import (
	"database/sql"
	"fmt"
	)
type User struct{
	UserID int
	UserName string
	UserRepo int
}
type Object struct{
	Cost int//首字母大写以允许外部调用
}
func CheckErr(err error){
	if err!=nil{
		panic(err)
}
}
func FetchUserInfo(db *sql.DB,id int)(*User){
	temp:=new(User)
	rows,err:=db.Query("select * from UserInfo where UserID=?",id)
	CheckErr(err)
	for rows.Next(){
		err:=rows.Scan(&temp.UserID,&temp.UserName,&temp.UserRepo)
		CheckErr(err)
	}
	return temp
}
func ShowMysqlData(db *sql.DB){
	rows,err:=db.Query("select * from UserInfo")
	CheckErr(err)
	for rows.Next(){
		var id int
		var name string
		var repo int
		err:=rows.Scan(&id,&name,&repo)
		CheckErr(err)
		fmt.Println(id,name,repo)
	}

}
