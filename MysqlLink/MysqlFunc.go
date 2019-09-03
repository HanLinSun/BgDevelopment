package MysqlLink
//购买消费同步机制
import (
	"database/sql"
	"fmt"
	"log"

)
type User struct{
	UserID int `json:"user_ID" form:"user_ID"`
	UserRepo int `json:user_Repo form:user_Repo`
}
type Object struct{
	Cost int//首字母大写以允许外部调用
}

func CheckErr(err error)error{
	if err!=nil{
		fmt.Println("failure:",err.Error())
}
return err
}

var SqlDB *sql.DB//变量大写以方便其他地方调用(全局变量)

func InitSql(){
	var err error
	SqlDB,err:=sql.Open("mysql","peo:peo123@tcp(124.156.206.94)/UserCount?charset=utf8")
	if err!=nil{
		log.Fatal(err.Error())
	}
	err=SqlDB.Ping()
	if err!=nil{
		log.Fatal(err.Error())
	}
}
func (user *User)FetchUserInfo(db *sql.DB) (retn User,err error){
	InitSql()
	rows,err:=db.Query("select * from UserInfo where UserID=?",user.UserID)
	CheckErr(err)
	for rows.Next(){
		err:=rows.Scan(&retn.UserID,&retn.UserRepo)
		CheckErr(err)
	}
	defer rows.Close()
	return
}
func (user *User)FetchAlluser(db *sql.DB)(users []User,err error){
	users = make([]User,0)
	rows,err:=db.Query("select UserID,Repo from UserInfo")
	CheckErr(err)
	for rows.Next(){
		var user User
		rows.Scan(&user.UserID,&user.UserRepo)
		users = append(users,user)
	}
	defer rows.Close()
	return
}
func ShowMysqlData(db *sql.DB){
	InitSql()
	rows,err:=db.Query("select * from UserInfo")
	CheckErr(err)
	for rows.Next(){
		var id int
		var repo int
		err:=rows.Scan(&id,&repo)
		CheckErr(err)
		fmt.Println(id,repo)
	}

}
