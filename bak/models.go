package models

import (
	"../utils"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/coocood/qbs"
	//_ "github.com/lib/pq" //当某时间字段表现为0001-01-01 07:36:42+07:36:42形式的时候 会读不出数据
	_ "github.com/bylevel/pq"
	//_ "github.com/mattn/go-sqlite3"
	"errors"
	"os"
	"time"
)

const (
	DbName         = "./data/sqlite.db"
	DbUser         = "root"
	mysqlDriver    = "mymysql"
	mysqlDrvformat = "%v/%v/"
	pgDriver       = "postgres"
	pgDrvFormat    = "user=%v dbname=%v sslmode=disable"
	sqlite3Driver  = "sqlite3"
	dbtypeset      = "pgsql"
)

type User struct {
	Id            int64
	Email         string `qbs:"index"`
	Password      string
	Nickname      string `qbs:"index"`
	Realname      string
	Avatar        string
	Avatar_min    string
	Avatar_max    string
	Birth         time.Time
	Province      string
	City          string
	Company       string
	Address       string
	Postcode      string
	Mobile        string
	Website       string
	Sex           int64
	Qq            string
	Msn           string
	Weibo         string
	Ctype         int64
	Role          int64
	Created       time.Time `qbs:"index"`
	Hotness       float64   `qbs:"index"`
	Hotup         int64     `qbs:"index"`
	Hotdown       int64     `qbs:"index"`
	Hotscore      int64     `qbs:"index"`
	Views         int64     `qbs:"index"`
	LastLoginTime time.Time
	LastLoginIp   string
	LoginCount    int64
}

//category,Pid:root
type Category struct {
	Id             int64
	Pid            int64 `qbs:"index"`
	Uid            int64 `qbs:"index"`
	Ctype          int64
	Title          string
	Content        string
	Attachment     string
	Created        time.Time `qbs:"index"`
	Hotness        float64   `qbs:"index"`
	Hotup          int64     `qbs:"index"`
	Hotdown        int64     `qbs:"index"`
	Hotscore       int64     `qbs:"index"`
	Views          int64     `qbs:"index"`
	Author         string
	NodeTime       time.Time
	NodeCount      int64
	NodeLastUserId int64
}

//node,Pid:category
type Node struct {
	Id              int64
	Pid             int64 `qbs:"index"`
	Uid             int64 `qbs:"index"`
	Ctype           int64
	Title           string
	Content         string
	Attachment      string
	Created         time.Time `qbs:"index"`
	Updated         time.Time `qbs:"index"`
	Hotness         float64   `qbs:"index"`
	Hotup           int64     `qbs:"index"`
	Hotdown         int64     `qbs:"index"`
	Hotscore        int64     `qbs:"index"`
	Views           int64     `qbs:"index"`
	Author          string
	TopicTime       time.Time
	TopicCount      int64
	TopicLastUserId int64
}

//topic,Pid:node
type Topic struct {
	Id              int64
	Cid             int64 `qbs:"index"`
	Nid             int64 `qbs:"index"`
	Uid             int64 `qbs:"index"`
	Ctype           int64
	Title           string
	Content         string
	Attachment      string
	Created         time.Time `qbs:"index"`
	Updated         time.Time `qbs:"index"`
	Hotness         float64   `qbs:"index"`
	Hotup           int64     `qbs:"index"`
	Hotdown         int64     `qbs:"index"`
	Hotscore        int64     `qbs:"index"`
	Views           int64     `qbs:"index"`
	Author          string
	ReplyTime       time.Time
	ReplyCount      int64
	ReplyLastUserId int64
}

//reply,Pid:topic
type Reply struct {
	Id         int64
	Uid        int64 `qbs:"index"`
	Pid        int64 `qbs:"index"` //Topic id
	Ctype      int64
	Content    string
	Attachment string
	Created    time.Time `qbs:"index"`
	Hotness    float64   `qbs:"index"`
	Hotup      int64     `qbs:"index"`
	Hotdown    int64     `qbs:"index"`
	Hotscore   int64     `qbs:"index"`
	Views      int64     `qbs:"index"`
	Author     string
	Email      string
	Website    string
}

type File struct {
	Id              int64
	Cid             int64 `qbs:"index"`
	Nid             int64 `qbs:"index"`
	Uid             int64 `qbs:"index"`
	Pid             int64 `qbs:"index"`
	Ctype           int64
	Filename        string
	Content         string
	Hash            string
	Location        string
	Url             string
	Size            int64
	Created         time.Time `qbs:"index"`
	Updated         time.Time `qbs:"index"`
	Hotness         float64   `qbs:"index"`
	Hotup           int64     `qbs:"index"`
	Hotdown         int64     `qbs:"index"`
	Hotscore        int64     `qbs:"index"`
	Views           int64     `qbs:"index"`
	ReplyTime       time.Time
	ReplyCount      int64
	ReplyLastUserId int64
}

// k/v infomation
type Kvs struct {
	Id int64
	/*
		Cid int64
		Nid int64
		Tid int64
		Rid int64
	*/
	K string
	V string
}

//Subscribe 订阅之邮件列表
type Subscribe struct {
	Id      int64
	Email   string
	Ctype   int64 //0为接受订阅，1为用户取消订阅。
	Created time.Time
}

func RegisterDb() {

	switch {
	case dbtypeset == "sqlite":
		qbs.Register("sqlite3", "./data/sqlite.db", "", qbs.NewSqlite3())

	case dbtypeset == "mysql":
		qbs.Register("mysql", "qbs_test@/qbs_test?charset=utf8&parseTime=true&loc=Local", "dbname", qbs.NewMysql())

	case dbtypeset == "pgsql":
		// "postgres://postgres:LhS88root@110.76.39.205:5432/postgres"
		//qbs.Register("postgres", "user=postgres password=LhS88root dbname=test sslmode=disable", "test", qbs.NewPostgres())
		qbs.Register("postgres", "user=postgres password=jn!@#$%^&* dbname=pgsql sslmode=disable", "pgsql", qbs.NewPostgres())

	}
	//设置最小连接池大小
	qbs.ChangePoolSize(2)
	//设置最大的数据库连接
	qbs.SetConnectionLimit(3, false)
}

func init() {
	RegisterDb()
	var install string
	//如果不存在配置文件，则跳过安装过程
	if conf, err := goconfig.LoadConfigFile("./conf/config.conf"); err != nil {
		//为了保障数据安全，只有存在配置文件，并且配置文件设置允许全新安装，这样才会触发安装过程，不然一概忽略！

		//install = "true"则开始全新安装，install = "false"则跳过安装
		install = "false"
	} else {
		//如果存在配置文件，但 是否全新安装的状态设置 读取错误则设置为允许安装状态
		if install, err = conf.GetValue("install", "allow"); err != nil {
			install = "true"
		}

		if install == "true" {

			q, err := qbs.GetQbs()
			defer q.Close()
			if err != nil {
				fmt.Println("Veryhour system connect to database error:", err)

			} else {
				mg, _ := qbs.GetMigration()
				defer mg.Close()

				mg.CreateTableIfNotExists(new(User))
				mg.CreateTableIfNotExists(new(Category))
				mg.CreateTableIfNotExists(new(Node))
				mg.CreateTableIfNotExists(new(Topic))
				mg.CreateTableIfNotExists(new(Reply))
				mg.CreateTableIfNotExists(new(Kvs))
				mg.CreateTableIfNotExists(new(File))
				mg.CreateTableIfNotExists(new(Subscribe))

				//用户等级划分：正数是普通用户，负数是管理员各种等级划分，为0则尚未注册
				if GetUserByRole(-1000).Role != -1000 {
					AddUser("root@localhost", "root", "系统默认管理员", utils.Encrypt_hash("rootpass", nil), -1000)
					fmt.Println("Default User:root,Password:rootpass")

					if GetAllTopic(0, 0, "id") == nil {
						//分類默認數據
						AddCategory("Category！", "This is Category！")

						AddNode("Node！", "This is Node!", 1, 1)
						SetTopic(0, 1, 1, 1, 0, "Topic Title", `<p>This is Topic!</p>`, "root", "")

					}
				}

				if GetKV("author") != "Insion" {
					SetKV("author", "Insion")
					SetKV("title", "Veryhour")
					SetKV("title_en", "Veryhour")
					SetKV("keywords", "Veryhour,")
					SetKV("description", "Veryhour,")

					SetKV("company", "Veryhour")
					SetKV("copyright", "2013 Copyright Veryhour .All Right Reserved")
					SetKV("site_email", "info@verywave.com")
				}
				//成功安装数据后把安装状态改为禁止安装
				conf.SetValue("install", "allow", "false")
				goconfig.SaveConfigFile(conf, "./conf/config.conf")

				fmt.Println("Veryhour system is successfully installed!")
			}

		} else {
			fmt.Println("Veryhour system starts running..")
		}
	}

}

func Counts() (categorys int, nodes int, topics int, menbers int) {
	q, _ := qbs.GetQbs()
	defer q.Close()

	var categoryz []*Category
	if e := q.FindAll(&categoryz); e != nil {
		categorys = 0
		fmt.Println(e)
	} else {
		categorys = len(categoryz)
	}

	var nodez []*Node
	if e := q.FindAll(&nodez); e != nil {
		nodes = 0
		fmt.Println(e)
	} else {
		nodes = len(nodez)
	}

	var topicz []*Topic
	if e := q.FindAll(&topicz); e != nil {
		topics = 0
		fmt.Println(e)
	} else {
		topics = len(topicz)
	}

	var menberz []*User
	if e := q.FindAll(&menberz); e != nil {
		menbers = 0
		fmt.Println(e)
	} else {
		menbers = len(menberz)
	}

	return categorys, nodes, topics, menbers
}

func TopicCount() (today int, this_week int, this_month int) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	var topict, topicw, topicm []*Topic
	k := time.Now()

	//一天之前
	d, _ := time.ParseDuration("-24h")
	t := k.Add(d)
	e := q.Where("created>?", t).FindAll(&topict)
	if e != nil {
		today = 0
		fmt.Println(e)
	} else {
		today = len(topict)
	}

	//一周之前
	w := k.Add(d * 7)
	e = q.Where("created>?", w).FindAll(&topicw)
	if e != nil {
		this_week = 0
		fmt.Println(e)
	} else {
		this_week = len(topicw)
	}

	//一月之前
	m := k.Add(d * 30)
	e = q.Where("created>?", m).FindAll(&topicm)
	if e != nil {
		this_month = 0
		fmt.Println(e)
	} else {
		this_month = len(topicm)
	}

	return today, this_week, this_month
}

func SetTopic(id int64, cid int64, nid int64, uid int64, ctype int64, title string, content string, author string, attachment string) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	var tp Topic
	if q.WhereEqual("id", id).Find(&tp); tp.Id == 0 {
		_, err := q.Save(&Topic{Id: id, Cid: cid, Nid: nid, Uid: uid, Ctype: ctype, Title: title, Content: content, Author: author, Attachment: attachment})
		return err
	} else {
		type Topic struct {
			Cid        int64
			Nid        int64
			Uid        int64
			Ctype      int64
			Title      string
			Content    string
			Author     string
			Attachment string
		}

		_, err := q.WhereEqual("id", id).Update(&Topic{Cid: cid, Nid: nid, Uid: uid, Ctype: ctype, Title: title, Content: content, Author: author, Attachment: attachment})
		return err
	}
	return nil
}

func AddFile(ctype int64, location string, url string) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.Save(&File{Ctype: ctype, Location: location, Url: url})
	return err
}

func DelFile(id int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	f := GetFile(id)

	if utils.Exist("." + f.Location) {
		if err := os.Remove("." + f.Location); err != nil {
			return err
			fmt.Println(err)
		}
	}

	//不管实际路径中是否存在文件均删除该数据库记录，以免数据库记录陷入死循环无法删掉
	_, err := q.Delete(&f)
	fmt.Println(err)
	return err
}

func GetFile(id int64) (f File) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.Where("id=?", id).Find(&f)
	return f
}

func GetAllFile() (f []*File) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.OrderByDesc("id").FindAll(&f)
	return f
}

func GetAllFileByCtype(ctype int64) (f []*File) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.WhereEqual("ctype", ctype).OrderByDesc("id").FindAll(&f)
	return f
}

func SaveFile(f File) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, e := q.Save(&f)
	return e
}

func SetFile(id int64, pid int64, ctype int64, filename string, content string, hash string, location string, url string, size int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	var f File
	if q.WhereEqual("id", id).Find(&f); f.Id == 0 {
		_, err := q.Save(&File{Id: id, Pid: pid, Ctype: ctype, Filename: filename, Content: content, Hash: hash, Location: location, Url: url, Size: size})
		return err
	} else {
		type File struct {
			Pid      int64
			Ctype    int64
			Filename string
			Content  string
			Hash     string
			Location string
			Url      string
			Size     int64
		}
		_, err := q.WhereEqual("id", id).Update(&File{Pid: pid, Ctype: ctype, Filename: filename, Content: content, Hash: hash, Location: location, Url: url, Size: size})

		return err
	}
	return nil
}

func AddKV(k string, v string) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.Save(&Kvs{K: k, V: v})
	return err
}

func SetKV(k string, v string) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	var kvs Kvs
	if q.Where("k=?", k).Find(&kvs); kvs.Id == 0 {
		_, err := q.Save(&Kvs{K: k, V: v})
		return err
	} else {
		type Kvs struct {
			K string
			V string
		}

		_, err := q.WhereEqual("k", k).Update(&Kvs{K: k, V: v})

		return err
	}
	return nil
}

func GetKV(k string) (v string) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	var kvs Kvs
	q.Where("k=?", k).Find(&kvs)
	return kvs.V
}

func AddUser(email string, nickname string, realname string, password string, role int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.Save(&User{Email: email, Nickname: nickname, Realname: realname, Password: password, Role: role, Created: time.Now()})

	return err
}

func SaveUser(usr User) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, e := q.Save(&usr)
	return e
}

func UpdateUser(uid int, ur User) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.WhereEqual("id", int64(uid)).Update(&ur)
	return err
}

func GetUser(id int64) (user User) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.Where("id=?", id).Find(&user)
	return user
}

func DelUser(uid int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	usr := GetUser(uid)
	_, err := q.Delete(&usr)

	return err
}

func GetUserByRole(role int) (user User) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.Where("role=?", int64(role)).Find(&user)
	return user
}

func GetAllUserByRole(role int) (user []*User) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.Where("role=?", int64(role)).OrderByDesc("id").FindAll(&user)
	return user
}

func GetUserByNickname(nickname string) (user User) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.Where("nickname=?", nickname).Find(&user)
	return user
}

func AddCategory(title string, content string) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.Save(&Category{Title: title, Content: content, Created: time.Now()})

	return err
}

func SaveCategory(cat Category) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.Save(&cat)
	return err
}

func AddNode(title string, content string, cid int64, uid int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	if _, err := q.Save(&Node{Pid: cid, Uid: uid, Title: title, Content: content, Created: time.Now()}); err != nil {
		return err
	}

	type Category struct {
		NodeTime       time.Time
		NodeCount      int64
		NodeLastUserId int64
	}

	if _, err := q.WhereEqual("id", cid).Update(&Category{NodeTime: time.Now(), NodeCount: int64(len(GetAllNodeByCid(cid, 0, 0, 0, "id"))), NodeLastUserId: uid}); err != nil {
		return err
	}
	/*
		ctr := GetCategory(cid)
		ctr.NodeTime = time.Now()
		ctr.NodeCount = int64(len(GetAllNodeByCid(cid, 0, 0, "id")))
		ctr.NodeLastUserId = int64(uid)
		if _, err := q.Save(&ctr); err != nil {
			return err
		}
	*/
	return nil
}

func SetNode(id int64, title string, content string, cid int64, uid int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	var nd Node
	if q.WhereEqual("id", id).Find(&nd); nd.Id == 0 {
		_, err := q.Save(&Node{Id: id, Pid: cid, Uid: uid, Title: title, Content: content})
		return err
	} else {
		type Node struct {
			Pid     int64
			Uid     int64
			Title   string
			Content string
		}

		_, err := q.WhereEqual("id", id).Update(&Node{Pid: cid, Uid: uid, Title: title, Content: content})
		return err
	}
	return nil
}

func AddTopic(title string, content string, cid int64, nid int64, uid int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	if _, err := q.Save(&Topic{Cid: cid, Nid: nid, Title: title, Content: content, Created: time.Now()}); err != nil {
		return err
	}

	type Node struct {
		TopicTime       time.Time
		TopicCount      int64
		TopicLastUserId int64
	}

	if _, err := q.WhereEqual("id", nid).Update(&Node{TopicTime: time.Now(), TopicCount: int64(len(GetAllTopicByNid(nid, 0, 0, 0, "id"))), TopicLastUserId: uid}); err != nil {
		return err
	}
	/*
		nd := GetNode(nid)
		nd.TopicTime = time.Now()
		nd.TopicCount = int64(len(GetAllTopicByNid(nid, 0, 0, "id")))
		nd.TopicLastUserId = int64(uid)
		if _, err := q.Save(&nd); err != nil {
			return err
		}
	*/
	return nil
}

func AddReply(tid int64, uid int64, content string, author string, email string, website string) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	if _, err := q.Save(&Reply{Pid: tid, Uid: uid, Content: content, Created: time.Now(), Author: author, Email: email, Website: website}); err != nil {
		return err
	}

	type Topic struct {
		ReplyTime       time.Time
		ReplyCount      int64
		ReplyLastUserId int64
	}

	if _, err := q.WhereEqual("id", tid).Update(&Topic{ReplyTime: time.Now(), ReplyCount: int64(len(GetReplyByPid(tid, 0, 0, "id"))), ReplyLastUserId: uid}); err != nil {
		return err
	}
	/*
		tp := GetTopic(tid)
		tp.ReplyCount = int64(len(GetReplyByPid(tid, 0, 0, "id")))
		tp.ReplyTime = time.Now()
		tp.ReplyLastUserId = int64(uid)
		if _, err := q.Save(&tp); err != nil {
			return err
		}
	*/
	return nil
}

func SaveNode(nd *Node) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.Save(&nd)
	return err
}

func DelNodePlus(nid int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	node := GetNode(nid)
	_, err := q.Delete(&node)

	for i, v := range GetAllTopicByNid(nid, 0, 0, 0, "id") {
		if i > 0 {
			DelTopic(v.Id)
			for ii, vv := range GetReplyByPid(v.Id, 0, 0, "id") {
				if ii > 0 {
					DelReply(vv.Id)
				}
			}
		}
	}

	return err
}

func DelCategory(id int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	category := GetCategory(id)
	_, err := q.Delete(&category)

	return err
}

func DelTopic(id int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	topic := new(Topic)

	q.Where("id=?", id).Find(topic)

	if topic.Attachment != "" {

		if utils.Exist("." + topic.Attachment) {
			if err := os.Remove("." + topic.Attachment); err != nil {
				//return err
				//可以输出错误，但不要反回错误，以免陷入死循环无法删掉
				fmt.Println("DEL TOPIC", id, err)
			}
		}
	}

	//不管实际路径中是否存在文件均删除该数据库记录，以免数据库记录陷入死循环无法删掉
	if topic.Id == id {

		if _, err := q.Delete(topic); err != nil {

			fmt.Println(err)
		} else {
			return err
		}

	}
	return errors.New("无法删除不存在的TOPIC ID:" + string(id))
}

func DelNode(nid int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	node := GetNode(nid)
	_, err := q.Delete(&node)

	return err
}

func DelReply(tid int64) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	reply := GetReply(tid)
	_, err := q.Delete(&reply)

	return err
}

func GetAllCategory() (allc []*Category) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.FindAll(&allc)
	return allc
}

func GetAllNode() (alln []*Node) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	//q.OrderByDesc("id").FindAll(&alln)
	q.OrderByDesc("created").FindAll(&alln)
	return alln
}

func GetAllTopic(offset int, limit int, path string) (allt []*Topic) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("created").FindAll(&allt)
	return allt
}

func GetAllNodeByCid(cid int64, offset int, limit int, ctype int64, path string) (alln []*Node) {
	//排序首先是热值优先，然后是时间优先。
	q, _ := qbs.GetQbs()
	defer q.Close()
	switch {
	case path == "asc":
		if ctype != 0 {
			condition := qbs.NewCondition("pid=?", cid).And("ctype=?", ctype)
			q.Condition(condition).Offset(offset).Limit(limit).FindAll(&alln)
		} else {
			if cid == 0 {
				q.Offset(offset).Limit(limit).FindAll(&alln)
			} else {
				q.WhereEqual("pid", cid).Offset(offset).Limit(limit).FindAll(&alln)
			}

		}
	case path == "views" || path == "topic_count":
		if ctype != 0 {
			condition := qbs.NewCondition("pid=?", cid).And("ctype=?", ctype)
			q.Condition(condition).OrderByDesc(path).Offset(offset).Limit(limit).FindAll(&alln)

		} else {
			if cid == 0 {
				q.OrderByDesc(path).Offset(offset).Limit(limit).FindAll(&alln)
			} else {
				q.WhereEqual("pid", cid).OrderByDesc(path).Offset(offset).Limit(limit).FindAll(&alln)
			}

		}
	default:
		if ctype != 0 {

			condition := qbs.NewCondition("pid=?", cid).And("ctype=?", ctype)
			q.Condition(condition).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("topic_count").OrderByDesc("created").FindAll(&alln)

		} else {
			if cid == 0 {
				q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("topic_count").OrderByDesc("created").FindAll(&alln)
			} else {
				q.WhereEqual("pid", cid).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("topic_count").OrderByDesc("created").FindAll(&alln)
			}
		}

	}
	return alln
}

func GetAllTopicByCid(cid int64, offset int, limit int, ctype int64, path string) (allt []*Topic) {
	//排序首先是热值优先，然后是时间优先。
	q, _ := qbs.GetQbs()
	defer q.Close()

	switch {
	case path == "asc":
		if ctype != 0 {
			condition := qbs.NewCondition("cid=?", cid).And("ctype=?", ctype)
			q.Condition(condition).Offset(offset).Limit(limit).FindAll(&allt)

		} else {
			q.Where("cid=?", cid).Offset(offset).Limit(limit).FindAll(&allt)

		}
	case path == "views" || path == "reply_count":
		if ctype != 0 {
			condition := qbs.NewCondition("cid=?", cid).And("ctype=?", ctype)
			q.Condition(condition).OrderByDesc(path).Offset(offset).Limit(limit).FindAll(&allt)

		} else {
			if cid == 0 {
				q.OrderByDesc(path).Offset(offset).Limit(limit).FindAll(&allt)
			} else {
				q.WhereEqual("cid", cid).OrderByDesc(path).Offset(offset).Limit(limit).FindAll(&allt)
			}

		}
	default:
		if ctype != 0 {

			condition := qbs.NewCondition("cid=?", cid).And("ctype=?", ctype)
			q.Condition(condition).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)

		} else {
			if cid == 0 {
				q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)

			} else {
				q.WhereEqual("cid", cid).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			}
		}

	}
	return allt
}

func GetAllTopicByCidNid(cid int64, nid int64, offset int, limit int, ctype int64, path string) (allt []*Topic) {

	q, _ := qbs.GetQbs()
	defer q.Close()

	switch {
	case path == "asc":
		if ctype != 0 {
			condition := qbs.NewCondition("cid=?", cid).And("nid=?", nid).And("ctype=?", ctype)
			q.Condition(condition).Offset(offset).Limit(limit).FindAll(&allt)

		} else {

			condition := qbs.NewCondition("cid=?", cid).And("nid=?", nid)
			q.Condition(condition).Offset(offset).Limit(limit).FindAll(&allt)

		}
	default:
		if ctype != 0 {
			condition := qbs.NewCondition("cid=?", cid).And("nid=?", nid).And("ctype=?", ctype)
			q.Condition(condition).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)

		} else {

			condition := qbs.NewCondition("cid=?", cid).And("nid=?", nid)
			q.Condition(condition).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)

		}

	}
	return allt
}

func GetAllTopicByNid(nodeid int64, offset int, limit int, ctype int64, path string) (allt []*Topic) {
	//排序首先是热值优先，然后是时间优先。
	q, _ := qbs.GetQbs()
	defer q.Close()

	switch {
	case path == "asc":
		if nodeid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				condition := qbs.NewCondition("nid=?", nodeid).And("ctype=?", ctype)
				q.Condition(condition).Offset(offset).Limit(limit).FindAll(&allt)

			} else {
				q.Where("nid=?", nodeid).Offset(offset).Limit(limit).FindAll(&allt)

			}
		}
	default:
		if nodeid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				condition := qbs.NewCondition("nid=?", nodeid).And("ctype=?", ctype)
				q.Condition(condition).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)

			} else {
				q.Where("nid=?", nodeid).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)

			}
		}
	}
	return allt
}

func SearchTopic(content string, offset int, limit int, path string) (allt []*Topic) {
	//排序首先是热值优先，然后是时间优先。
	if content != "" {
		q, _ := qbs.GetQbs()
		defer q.Close()
		keyword := "%" + content + "%"
		condition := qbs.NewCondition("title like ?", keyword).Or("content like ?", keyword)
		q.Condition(condition).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
		//q.Where("title like ?", keyword).Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("created").FindAll(&allt)
		return allt
	}
	return nil
}

func GetCategory(id int64) (category Category) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.Where("id=?", id).Find(&category)
	return category
}

func GetNode(id int64) (node *Node) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.WhereEqual("id", id).Find(&node)
	return node
}

func GetTopic(id int64) (tp *Topic) {
	q, _ := qbs.GetQbs()
	defer q.Close()

	q.Where("id=?", id).Find(&tp)
	//fmt.Println(*topic) //输出内容
	//fmt.Println(&topic) //输出内存地址
	return tp
}

func PostTopic(tp Topic) (int64, error) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.Save(&tp)

	return tp.Id, err
}

func PutTopic(tid int64, tp Topic) error {
	/*
		由于无法预知前端提交过来的数据哪些字段是更新过的哪些是没更新的，导致无法预知需要修改什么字段。
		以及qbs更新的限制，这里采用完整数据覆盖式更新，前端先读出完整旧数据，然后修改对应的字段后，完整的提交到这里覆盖式保存。
	*/
	q, _ := qbs.GetQbs()
	defer q.Close()

	_, err := q.WhereEqual("id", tid).Update(&tp) //不建临时结构，直接覆盖更新！
	return err
}

func SaveTopic(tp *Topic) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.Save(tp)

	return err
}

func UpdateCategory(cid int64, cg Category) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.WhereEqual("id", int64(cid)).Update(&cg)
	return err
}

func UpdateNode(nid int64, nd *Node) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.WhereEqual("id", int64(nid)).Update(&nd)

	return err
}

func UpdateTopic(tid int64, tp *Topic) error {
	q, _ := qbs.GetQbs()
	defer q.Close()
	_, err := q.WhereEqual("id", int64(tid)).Update(tp)
	return err
}

func EditNode(nid int64, cid int64, uid int64, title string, content string) error {
	nd := GetNode(nid)
	nd.Pid = cid
	nd.Title = title
	nd.Content = content
	nd.Updated = time.Now()
	if err := UpdateNode(nid, nd); err != nil {
		return err
	}

	q, _ := qbs.GetQbs()
	defer q.Close()

	type Category struct {
		NodeTime       time.Time
		NodeCount      int64
		NodeLastUserId int64
	}

	if _, err := q.WhereEqual("id", cid).Update(&Category{NodeTime: time.Now(), NodeCount: int64(len(GetAllNodeByCid(cid, 0, 0, 0, "id"))), NodeLastUserId: int64(uid)}); err != nil {
		return err
	}

	return nil
}

func EditTopic(tid int64, nid int64, cid int64, uid int64, title string, content string) error {
	tpc := GetTopic(tid)
	tpc.Cid = int64(cid)
	tpc.Nid = int64(nid)
	tpc.Title = title
	tpc.Content = content
	tpc.Updated = time.Now()

	if err := UpdateTopic(tid, tpc); err != nil {
		return err
	}

	q, _ := qbs.GetQbs()
	defer q.Close()

	type Node struct {
		TopicTime       time.Time
		TopicCount      int64
		TopicLastUserId int64
	}

	if _, err := q.WhereEqual("id", nid).Update(&Node{TopicTime: tpc.Created, TopicCount: int64(len(GetAllTopicByNid(nid, 0, 0, 0, "id"))), TopicLastUserId: int64(uid)}); err != nil {
		return err
	}

	return nil
}

func GetAllReply() (allr []*Reply) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.OrderByDesc("id").FindAll(&allr)
	return allr
}

func GetReply(id int64) (reply Reply) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	q.Where("id=?", id).Find(&reply)
	return reply
}

func GetReplyByPid(tid int64, offset int, limit int, path string) (allr []*Reply) {
	q, _ := qbs.GetQbs()
	defer q.Close()
	if tid == 0 {
		q.Offset(offset).Limit(limit).OrderByDesc(path).FindAll(&allr)
	} else {
		//最热回复
		//q.Where("pid=?", tid).Offset(offset).Limit(limit).OrderByDesc("hotness").FindAll(&allr)
		q.WhereEqual("pid", tid).Offset(offset).Limit(limit).OrderByDesc(path).FindAll(&allr)
	}
	return allr
}

/*
func main() {

	ct()

	for i := 0; i < 100; i++ {
		AddCategory("我系标题", "我系内容啊~")
	}
	for i := 0; i < 100; i++ {
		AddUser("insion@lihuaer.com", "insion", "huhjj897857hggfgjhghsjg")
	}
	for i := 0; i < 100; i++ {
		AddNode("node title", "node content")
	}
	for i := 0; i < 100; i++ {
		AddTopic("topic title", "topic content")
	}
	for i := 0; i < 100; i++ {
		AddReply(int64(i), "a reply's content")
	}
	cc := GetAllCategory()
	for _, info := range cc {
		fmt.Println(info.Content)
	}

	c := GetCategory(1)
	fmt.Println(c.Content)

	n := GetNode(1)
	fmt.Println(n.Title)

	t := GetTopic(1)
	fmt.Println(t.Content)

	r := GetReply(1)
	fmt.Println(r.Content)

	for _, info := range GetAllCategory() {
		fmt.Println(info.Title)
	}
	for _, info := range GetAllNode() {
		fmt.Println(info.Content)
	}
	for _, info := range GetAllTopic() {
		fmt.Println(info.Title)
	}
	for _, info := range GetAllReply() {
		fmt.Println(info.Content)
	}
}
*/
