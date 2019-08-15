package RedisLink

import (
	"github.com/garyburd/redigo/redis"
	"CountSystem/MysqlLink"
	"fmt"
	"database/sql"
	"strconv"
)
func InsertIntoRedis(c redis.Conn,UserInfo *MysqlLink.User){
	fmt.Println("数据写入redis")
	_,err:=c.Do("MSET","UserID",UserInfo.UserID,"UserName",UserInfo.UserName,"UserRepo",UserInfo.UserRepo)
	if err!=nil{
		panic(err)
	}
	ShowRedisData(c)

}

func Purchase(c redis.Conn,db *sql.DB,obj MysqlLink.Object,userid int) bool{
	idData,err:=redis.Int(c.Do("GET","UserID"))//这
	if idData==userid { //传入的购买者ID和Redis内ID一致，就证明该用户信息存储在Redis中
	    fmt.Println("在Redis中进行操作")
		RepoString, err := redis.String(c.Do("GET", "UserRepo"))
		if err!=nil{
			fmt.Println("failure:",err)
			return false
		}
		RepoData,err:=strconv.Atoi(RepoString)//获取该用户的存款数量
		if obj.Cost > RepoData{
			fmt.Println("购买失败，余额不足!目前余额为:",RepoData)
			return false
		} else{
			RepoData=RepoData-obj.Cost
			_,err:=c.Do("SET","UserRepo",RepoData)
			MysqlLink.CheckErr(err)
			fmt.Println("购买成功，购买后的余额为:",RepoData)
		}
		RedisToSql(db,c)//数据同步至数据库

	} else {//若传入数据ID和Redis内不匹配，则访问数据库
	fmt.Println("该用户数据不在redis中,访问数据库")
	temp:=new(MysqlLink.User)
	rows,err:=db.Query("Select * from UserInfo where UserID=?",userid)
	MysqlLink.CheckErr(err)
	for rows.Next(){
		err:=rows.Scan(&temp.UserID,&temp.UserName,&temp.UserRepo)
		MysqlLink.CheckErr(err)
	}//数据库内和ID匹配的用户数据保存至临时temp变量
	if obj.Cost>temp.UserRepo{//商品价格大于用户余额
		fmt.Println("购买失败，余额不足.目前余额为:",temp.UserRepo)
		return false
	} else{

		temp.UserRepo=temp.UserRepo-obj.Cost
		fmt.Println("购买成功，余额为:",temp.UserRepo)
		stmt,err:=db.Prepare("update UserInfo set Repo=? where UserID=?")
		MysqlLink.CheckErr(err)
		result,err:=stmt.Exec(temp.UserRepo,temp.UserID)
		MysqlLink.CheckErr(err)
		result.LastInsertId()
	}
	SqlToRedis(db,c,userid)
	}
	MysqlLink.CheckErr(err)
   return true
}
func ShowRedisData(c redis.Conn){//显示redis内的元素
   fmt.Println("redis内的元素:")
   r2,err:=redis.Strings(c.Do("MGET","UserID","UserName","UserRepo"))
   MysqlLink.CheckErr(err)
   fmt.Println(r2)

}
func RedisToSql(db *sql.DB,c redis.Conn){//redis到sql
    temp:=new(MysqlLink.User)
    IdData,err:=redis.Int(c.Do("GET","UserID"))
    temp.UserID=IdData
	NameData,err:=redis.String(c.Do("GET","UserName"))
	MysqlLink.CheckErr(err)
	temp.UserName=NameData
	MysqlLink.CheckErr(err)
	RepoData,err:=redis.Int(c.Do("GET","UserRepo"))
	MysqlLink.CheckErr(err)
	temp.UserRepo=RepoData
	//这么写冗余代码太多,太麻烦,还容易出错。但我还没搞清Mget返回的到底是个啥东西，能不能一次直接存进结构体，所以先用笨方法
	//上述代码后面会想办法简化,这么干太麻烦(也许可以用Map试试)
	stmt,err:=db.Prepare("update UserInfo set Repo=? where UserID=?")
	MysqlLink.CheckErr(err)
	result,err:=stmt.Exec(temp.UserRepo,temp.UserID)
	MysqlLink.CheckErr(err)
	result.LastInsertId()
}
func SqlToRedis(db *sql.DB,c redis.Conn,userid int){
	temp:=new(MysqlLink.User)
	rows,err:=db.Query("select * from UserInfo where UserID=?",userid)
	MysqlLink.CheckErr(err)
	for rows.Next(){
		err:=rows.Scan(&temp.UserID,&temp.UserName,&temp.UserRepo)
		MysqlLink.CheckErr(err)
	}//根据ID查找到的mysql对应用户数据存入临时temp
	_, err = c.Do("MSET", "UserID", temp.UserID, "UserName", temp.UserName, "UserRepo", temp.UserRepo)
	if err!=nil{
		fmt.Println("failure:",err)
	}
	ShowRedisData(c)
}
func Charge(db *sql.DB,c redis.Conn,money int,userid int){//充值直接对数据库操作
	stmt,err:=db.Prepare("update UserInfo set Repo=? where UserID=? ")
	MysqlLink.CheckErr(err)
	result,err:=stmt.Exec(money,userid)
	result.LastInsertId()
	SqlToRedis(db,c,userid)
}