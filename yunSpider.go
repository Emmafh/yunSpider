package main

import (
	"database/sql"
	"encoding/json"
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

type follow struct {
	//Request_id int64
	Total_count int
	Follow_list []follow_list
	Errno       int
}
type follow_list struct {
	Follow_uk      int64
	Follow_uname   string
	Avatar_url     string
	Intro          string
	Pubshare_count int
	Follow_count   int
	Fans_count     int
	User_type      int
	Is_vip         int
}

//Mysql初始化
func init() {
	db, err = sql.Open("mysql", "root@(127.0.0.1:3306)/baidu")
	if err != nil {
		log.Error("数据库连接出错")
	}
	db.SetMaxOpenConns(50)
}

func main() {
	//getFollowList(271528372, 0)
	//	fmt.Println(int64(224 / 24))
	getTotalFollow(271528372)

}

func getTotalFollow(uk int64) (err error) {
	totalNum, err := getFollowList(uk, 0)
	if err == nil {
		pageNum := int(totalNum/24 + 1)
		for i := 1; i < pageNum; i++ {
			start := i * 24
			getFollowList(uk, start)
		}

	} else {
		log.Error(err.Error())
	}
	return err
}

func getFollowList(uk int64, start int) (followNum int, err error) {
	url := "http://yun.baidu.com/pcloud/friend/getfollowlist?query_uk=%d&limit=24&start=%d&bdstoken=e6f1efec456b92778e70c55ba5d81c3d&channel=chunlei&clienttype=0&web=1&logid=MTQ3NDA3NDg5NzU4NDAuMzQxNDQyMDY2MjA5NDA4NjU="
	real_url := fmt.Sprintf(url, uk, start)
	res, err := HttpGet(real_url, headers)
	if err != nil {
		return 0, err
	}

	var f follow
	err = json.Unmarshal([]byte(res), &f)
	if err != nil {
		return 0, err
	}

	for _, v := range f.Follow_list {
		ifExist := ifKeyExist(v.Follow_uk)
		if !ifExist {
			setUinfo(v)
		}
	}
	return f.Total_count, nil
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

func setUinfo(i follow_list) {
	res, err := db.Exec("INSERT INTO uinfo(uk,uname,avatar_url,intro,pubshare_count,follow_count,fans_count,user_type,is_vip) VALUES(?,?,?,?,?,?,?,?,?)", i.Follow_uk, i.Follow_uname, i.Avatar_url, i.Intro, i.Pubshare_count, i.Follow_count, i.Fans_count, i.User_type, i.Is_vip)
	checkErr(err)
	id, err := res.LastInsertId()

	uinfoId = id
	checkErr(err)
	log.Info("insert uinfo，uk:", i.Follow_uk, ",Id:", uinfoId)
}

func ifKeyExist(uk int64) bool {
	var id int
	err = db.QueryRow("SELECT id FROM uinfo where uk=?", uk).Scan(&id)
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
