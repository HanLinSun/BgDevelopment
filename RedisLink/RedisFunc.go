package RedisLink

import (
	"github.com/garyburd/redigo/redis"
	"CountSystem/MysqlLink"
	"fmt"
	"database/sql"
	"strconv"

	"time"
)
var Pool *redis.Pool
func InitRedis(){//redis连接池
	Pool= &redis.Pool{
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
func InsertIntoRedis(c redis.Conn,UserInfo *MysqlLink.User){
	fmt.Println("数据写入redis")
	_,err:=c.Do("MSET","UserID",UserInfo.UserID,"UserRepo",UserInfo.UserRepo)
	if err!=nil{//使用set会相对提高性能
		panic(err)
	}
	ShowRedisData(c)
}

func Purchase(c redis.Conn,db *sql.DB,obj MysqlLink.Object,userid int) bool{//把error写到返回值里(返回一个error类型的变量),方便debug
	idData,err:=redis.Int(c.Do("GET","UserID"))
	if idData==userid { //传入的购买者ID和Redis内ID一致，就证明该用户信息存储在Redis中
	    fmt.Println("在Redis中进行操作")
		RepoString, err := redis.String(c.Do("GET", "UserRepo"))
		if err!=nil{
			fmt.Println("failure:",err)
			return false
		}
		RepoData,err:=strconv.Atoi(RepoString)//获取该用户的存款数量(会避免一些问题)
		if obj.Cost > RepoData{
			fmt.Println("购买失败，余额不足!目前余额为:",RepoData)
			return false
		} else {
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
		err:=rows.Scan(&temp.UserID,&temp.UserRepo)
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

func ShowRedisData(c redis.Conn) error{//显示redis内的元素
   fmt.Println("redis内的元素:")
   r2,err:=redis.Strings(c.Do("MGET","UserID","UserRepo"))
   MysqlLink.CheckErr(err)
   fmt.Println(r2)
   return err
}

func RedisToSql(db *sql.DB,c redis.Conn) {//redis到sql
    temp:=new(MysqlLink.User)
    IdData,err:=redis.Int(c.Do("GET","UserID"))
    temp.UserID=IdData
	if err!=nil{
		fmt.Println("failure:",err)
	}
	RepoData,err:=redis.Int(c.Do("GET","UserRepo"))
	if err!=nil{
		fmt.Println("failure:",err)
	}
	temp.UserRepo=RepoData
	//这么写冗余代码太多,太麻烦,还容易出错。但我还没搞清Mget返回的到底是个啥东西，能不能一次直接存进结构体，所以先用笨方法
	//上述代码后面会想办法简化,这么干太麻烦
	stmt,err:=db.Prepare("update UserInfo set Repo=? where UserID=?")
	if err!=nil{
		fmt.Println("failure:",err)
	}
	result,err:=stmt.Exec(temp.UserRepo,temp.UserID)
	MysqlLink.CheckErr(err)
	result.LastInsertId()

}

func SqlToRedis(db *sql.DB,c redis.Conn,userid int) error{
	temp:=new(MysqlLink.User)
	rows,err:=db.Query("select * from UserInfo where UserID=?",userid)
	MysqlLink.CheckErr(err)
	for rows.Next(){
		err:=rows.Scan(&temp.UserID,&temp.UserRepo)
		MysqlLink.CheckErr(err)
	}//根据ID查找到的mysql对应用户数据存入临时temp
	_, err = c.Do("MSET", "UserID", temp.UserID, "UserRepo", temp.UserRepo)
	if err!=nil{
		fmt.Println("failure:",err)
	}
	ShowRedisData(c)
	return err
}

func Charge(db *sql.DB,c redis.Conn,money int,userid int)error{//充值直接对数据库操作
	stmt,err:=db.Prepare("update UserInfo set Repo=? where UserID=? ")
	MysqlLink.CheckErr(err)
	result,err:=stmt.Exec(money,userid)
	result.LastInsertId()
	SqlToRedis(db,c,userid)
	return err
}
//作业：计费系统 把它做成server端
//用户控制和消费 用client调用
//插入额度，额度值调整 ，退款
//RPC(没法负载均衡) GRPC(google RPC)
//Restful API = http服务
//paas 语音介绍:paas(平台,不需要专门找客户端,底层微服务) saas is
//把如何唤醒做成微服务 PB协议
//if 数据包满足条件 唤醒， 数据包内包含账号登录信息
//功能单元做成微服务，创建角色(把具有统一功能单元抽象出来)，把它做成GRPC微服务(server端)
//GRPC协议 pb Server收到请求，进行处理
//把计费系统写成http服务
//框架:Ginpack 先不要用gin
//先根据网上的教程熟悉gin
//注意传输协议 基于gin做个demo
//最后再把计费系统改成http服务端
//做项目 把通用功能模块抽象出来做成微服务 把它们组装成一个游戏
//功能模块思想