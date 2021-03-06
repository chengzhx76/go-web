package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-web/mode"
	"sync"
)

//https://blog.csdn.net/wangshubo1989/article/details/75257614
//http://go-database-sql.org/retrieving.html

var dataSource *sql.DB
var once sync.Once

func init() {
	var err error
	once.Do(func() {
		dataSource, err = sql.Open("mysql", "pusher:pusher#cheng@tcp(139.196.35.134:3306)/pusher?charset=utf8")
	})
	if err != nil {
		fmt.Println("get db error")
	}
	//defer dataSource.Close()
}

func SaveUser(uid string, sendKey string, openid string) {
	stmt, _ := dataSource.Prepare("INSERT INTO user(uid, sendKey, openid) VALUES (?,?,?)")
	res, _ := stmt.Exec(uid, sendKey, openid)
	fmt.Println(res.LastInsertId())
}

func LoadByOpenid(openid string) (*mode.User, error) {
	var user mode.User
	err := dataSource.QueryRow("SELECT uid, sendKey FROM user WHERE=?", openid).Scan(&user.Uid, &user.SendKey)
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func LoadBySendKey(sendKey string) (*mode.User, error) {
	var user mode.User
	err := dataSource.QueryRow("SELECT uid, openid FROM user WHERE=?", sendKey).Scan(&user.Uid, &user.Openid)
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func main() {
	//db, err := sql.Open("mysql", "root:wangshubo@/test?charset=utf8")
	//checkErr(err)
	// insert
	stmt, err := dataSource.Prepare("INSERT user_info SET id=?,name=?")
	checkErr(err)

	res, err := stmt.Exec(1, "wangshubo")
	checkErr(err)

	// update
	stmt, err = dataSource.Prepare("update user_info set name=? where id=?")
	checkErr(err)

	res, err = stmt.Exec("wangshubo_update", 1)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// query
	rows, err := dataSource.Query("SELECT * FROM user_info")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string

		err = rows.Scan(&uid, &username)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
	}

	// delete
	stmt, err = dataSource.Prepare("delete from user_info where id=?")
	checkErr(err)

	res, err = stmt.Exec(1)
	checkErr(err)

	// query
	rows, err = dataSource.Query("SELECT * FROM user_info")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string

		err = rows.Scan(&uid, &username)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
	}

	dataSource.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
