package main

import (
	"database/sql"
	//	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/siddontang/go/log"
	"io/ioutil"
	"net/http"
	//	"regexp"
	//	"time"
)

var db *sql.DB
var err error
var headers = map[string]string{
	"User-Agent": "MQQBrowser/26 Mozilla/5.0 (Linux; U; Android 2.3.7; zh-cn; MB200 Build/GRJ22; CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
	"Referer":    "https://yun.baidu.com/share/home?uk=325913312#category/type=0"}

var uinfoId int64 = 0

//Mysql初始化
func init() {
	db, err = sql.Open("mysql", "root@(127.0.0.1:3306)/baidu")
	if err != nil {
		log.Error("数据库连接出错")
	}
	db.SetMaxOpenConns(50)
}

func main() {

	getFollows(1644403944, 0)
}

func getFollows(uk int64, start int) {
	ifExist := ifKeyExist(uk)
	if !ifExist {
		setUk(uk)
		res, _ := HttpGet("http://so.com", headers)
		fmt.Println(res)

	} else {

	}
}

func recureFollow() {

}

func HttpGet(url string, headers map[string]string) (result string, err error) {

	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("数据读取异常")
		return "", err
	}
	defer resp.Body.Close()
	return string(body), nil
}

func setUk(uk int64) {
	res, err := db.Exec("INSERT INTO avaiuk(uk) VALUES(?)", uk)
	checkErr(err)
	id, err := res.LastInsertId()

	uinfoId = id
	checkErr(err)
	log.Info("insert avaiuk，uk:", uk, ",Id:", uinfoId)
}

func ifKeyExist(uk int64) bool {
	var id int
	err = db.QueryRow("SELECT id FROM avaiuk where uk=?", uk).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Fatal("read data error")
		return false
	default:
		return true
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
