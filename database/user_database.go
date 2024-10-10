package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/link1st/gowebsocket/v2/models"
)

// 模拟一个简单的消息队列
var loginQueue = make(chan *models.UserLoginInfo, 1000)

// DB 是全局的数据库连接池
var DB *sql.DB

// 批量插入用户登录信息到 MySQL
func BatchInsertUserLogins(userLogins []*models.UserLoginInfo) error {
	if len(userLogins) == 0 {
		return nil
	}

	query := "INSERT INTO world_user (user_id, user_nickname, login_time) VALUES "
	values := []interface{}{}
	placeholders := []string{}

	for _, login := range userLogins {
		placeholders = append(placeholders, "(?, ?, ?)")
		values = append(values, login.UserId, login.NickName, login.LoginTime)
	}

	query += strings.Join(placeholders, ",")

	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	fmt.Println("成功插入用户登录信息")
	return nil
}

// ProcessLoginQueue 处理消息队列并批量写入 MySQL
func ProcessLoginQueue() {
	for {
		select {
		case userLoginInfo := <-loginQueue:
			err := BatchInsertUserLogins([]*models.UserLoginInfo{userLoginInfo})
			if err != nil {
				fmt.Println("批量插入用户登录信息失败:", err)
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

// AddUserToQueue 将用户信息添加到登录队列
func AddUserToQueue(userLoginInfo *models.UserLoginInfo) {
	select {
	case loginQueue <- userLoginInfo:
		fmt.Println("用户登录信息已加入队列:", userLoginInfo.UserId)
	default:
		fmt.Println("登录队列已满，无法处理请求:", userLoginInfo.UserId)
	}
}
