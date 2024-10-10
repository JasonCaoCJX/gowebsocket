// Package main 实现一个基于websocket分布式聊天(IM)系统。
package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/link1st/gowebsocket/v2/database"
	"github.com/link1st/gowebsocket/v2/lib/redislib"
	"github.com/link1st/gowebsocket/v2/routers"
	"github.com/link1st/gowebsocket/v2/servers/grpcserver"
	"github.com/link1st/gowebsocket/v2/servers/task"
	"github.com/link1st/gowebsocket/v2/servers/websocket"
)

func main() {
	initConfig()
	initFile()
	initRedis()
	initMysql()
	router := gin.Default()

	// 初始化路由
	routers.Init(router)
	routers.WebsocketInit()

	// 定时任务
	task.Init()

	// 服务注册
	task.ServerInit()
	go websocket.StartWebSocket()

	// grpc
	go grpcserver.Init()
	go open()
	httpPort := viper.GetString("app.httpPort")
	_ = http.ListenAndServe(":"+httpPort, router)

	// 确保在程序退出时关闭数据库连接
	defer database.DB.Close()
}

// 初始化日志
func initFile() {
	// 禁用控制台颜色，因为在将日志写入文件时不需要颜色
	gin.DisableConsoleColor()
	// 从配置文件中获取日志文件的路径
	logFile := viper.GetString("app.logFile")
	// 创建或打开日志文件，并将其设置为gin框架的日志输出目标
	f, _ := os.Create(logFile)
	gin.DefaultWriter = io.MultiWriter(f)
}

// 初始化项目配置
func initConfig() {
	// 设置配置文件名和路径
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".")
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 打印出配置项的值
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config redis:", viper.Get("redis"))
	fmt.Println("config redis:", viper.Get("mysql"))
}

// 初始化Redis客户端
func initRedis() {
	redislib.NewClient()
}

func initMysql() {
	// 初始化数据库连接
	host := viper.GetString("mysql.host")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	dbname := viper.GetString("mysql.database")
	mysqlConnect := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbname)

	var err error
	database.DB, err = sql.Open("mysql", mysqlConnect)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 检查数据库连接
	if err = database.DB.Ping(); err != nil {
		log.Fatalf("数据库无法连接: %v", err)
	}

	// 配置数据库连接池
	database.DB.SetMaxOpenConns(25)
	database.DB.SetMaxIdleConns(25)
	database.DB.SetConnMaxLifetime(5 * time.Minute)

	// 启动后台任务处理登录队列
	go database.ProcessLoginQueue()

	// 其他初始化代码...
	fmt.Println("服务已启动")
}

func open() {
	time.Sleep(1000 * time.Millisecond)
	httpUrl := viper.GetString("app.httpUrl")
	httpUrl = "http://" + httpUrl + "/home/index"
	fmt.Println("访问页面体验:", httpUrl)
	cmd := exec.Command("open", httpUrl)
	_, _ = cmd.Output()
}
