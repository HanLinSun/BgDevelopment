package main

import (
	_"github.com/go-sql-driver/mysql"
	"CountSystem/Ginpack"
	"CountSystem/RedisLink"
)

//数据结构体声明在mysqlLink中

func main(){
	RedisLink.InitRedis()
	router:=Ginpack.InitRouter()

	router.Run(":8082")
	/*
	User01:=new(MysqlLink.User)
	User02:=new(MysqlLink.User)
	Obj01:=MysqlLink.Object{1000}
	InitRedis()
	conn:=pool.Get()
    db,err:=sql.Open("mysql","peo:peo123@tcp(124.156.206.94:3306)/UserCount?charset=utf8")
    if err!=nil{
    	fmt.Println("failure:",err)
	}

    go func(){
    	time.Sleep(2*time.Minute)
    	RedisLink.RedisToSql(db,conn)
	}()

    User01=MysqlLink.FetchUserInfo(db,1)
	User02=MysqlLink.FetchUserInfo(db,2)
    MysqlLink.CheckErr(err)
    MysqlLink.ShowMysqlData(db)
	RedisLink.InsertIntoRedis(conn,User01)
	go func(){
		time.Sleep(2*time.Minute)
		err=RedisLink.RedisToSql(db,conn)
		MysqlLink.CheckErr(err)
		}()
	RedisLink.Charge(db,conn,2000,User01.UserID)
	RedisLink.Charge(db,conn,1900,User02.UserID)
    RedisLink.Purchase(conn,db,Obj01,User01.UserID)
    RedisLink.Purchase(conn,db,Obj01,User02.UserID)
	*/

}