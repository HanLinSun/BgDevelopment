package main

import (
	"CountSystem/MysqlLink"
	"github.com/garyburd/redigo/redis"
	"time"
	"CountSystem/RedisLink"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
)
var pool *redis.Pool
//数据结构体声明在mysqlLink中
func InitRedis(){//redis连接池
	pool= &redis.Pool{
		MaxIdle:3, //最大空闲连接数
		MaxActive:3, //最大激活连接数
		IdleTimeout:240*time.Second,
		Dial: func() (redis.Conn, error) {
			c,err:=redis.Dial("tcp","124.156.206.94:6379")
			if err!=nil{
				return nil,err
			}
			return c,err
		},
		TestOnBorrow:func(c redis.Conn,t time.Time)error{
			if time.Since(t)<time.Minute{
				return nil
			}
			_,err:=c.Do("PING")
			return err
		},
	}
}
func main(){
	User01:=new(MysqlLink.User)
	User02:=new(MysqlLink.User)
	Obj01:=MysqlLink.Object{1000}
	InitRedis()
	conn:=pool.Get()
    db,err:=sql.Open("mysql","peo:peo123@tcp(124.156.206.94:3306)/UserCount?charset=utf8")
    User01=MysqlLink.FetchUserInfo(db,1)
	User02=MysqlLink.FetchUserInfo(db,2)
    MysqlLink.CheckErr(err)
    MysqlLink.ShowMysqlData(db)
	RedisLink.InsertIntoRedis(conn,User01)
	RedisLink.Charge(db,conn,2000,User01.UserID)
	RedisLink.Charge(db,conn,1900,User02.UserID)
    RedisLink.Purchase(conn,db,Obj01,User01.UserID)
    RedisLink.Purchase(conn,db,Obj01,User02.UserID)
}