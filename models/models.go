package models

import (
	"errors"
	"fmt"
	_ "github.com/bylevel/pq"
	"github.com/insionng/veryhour/utils"
	"github.com/lunny/xorm"
	"os"
	//_ "github.com/mattn/go-sqlite3"
	"time"
)

var (
	Engine *xorm.Engine
)

type User struct {
	Id            int64
	Email         string
	Password      string
	Nickname      string
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
	Created       time.Time
	Hotness       float64
	Hotup         int64
	Hotdown       int64
	Hotscore      int64
	Views         int64
	LastLoginTime time.Time
	LastLoginIp   string
	LoginCount    int64
}

//category,Pid:root
type Category struct {
	Id             int64
	Pid            int64
	Uid            int64
	Ctype          int64
	Title          string
	Content        string
	Attachment     string
	Created        time.Time
	Hotness        float64
	Hotup          int64
	Hotdown        int64
	Hotscore       int64
	Views          int64
	Author         string
	NodeTime       time.Time
	NodeCount      int64
	NodeLastUserId int64
}

//node,Pid:category
type Node struct {
	Id              int64
	Pid             int64
	Uid             int64
	Ctype           int64
	Title           string
	Content         string
	Attachment      string
	Created         time.Time
	Updated         time.Time
	Hotness         float64
	Hotup           int64
	Hotdown         int64
	Hotscore        int64
	Views           int64
	Author          string
	TopicTime       time.Time
	TopicCount      int64
	TopicLastUserId int64
}

//topic,Pid:node
type Topic struct {
	Id              int64
	Cid             int64
	Nid             int64
	Uid             int64
	Ctype           int64
	Title           string
	Content         string
	Attachment      string
	Created         time.Time
	Updated         time.Time
	Hotness         float64
	Hotup           int64
	Hotdown         int64
	Hotscore        int64
	Views           int64
	Author          string
	ReplyTime       time.Time
	ReplyCount      int64
	ReplyLastUserId int64
}

//reply,Pid:topic
type Reply struct {
	Id         int64
	Uid        int64
	Pid        int64 //Topic id
	Ctype      int64
	Content    string
	Attachment string
	Created    time.Time
	Hotness    float64
	Hotup      int64
	Hotdown    int64
	Hotscore   int64
	Views      int64
	Author     string
	Email      string
	Website    string
}

type File struct {
	Id              int64
	Cid             int64
	Nid             int64
	Uid             int64
	Pid             int64
	Ctype           int64
	Filename        string
	Content         string
	Hash            string
	Location        string
	Url             string
	Size            int64
	Created         time.Time
	Updated         time.Time
	Hotness         float64
	Hotup           int64
	Hotdown         int64
	Hotscore        int64
	Views           int64
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

type Timeline struct {
	Id      int64
	Ctype   int64
	Content string
	Created time.Time
}

func init() {
	SetEngine()

	Engine.CreateTables(&User{})
	Engine.CreateTables(&Category{})
	Engine.CreateTables(&Node{})
	Engine.CreateTables(&Topic{})
	Engine.CreateTables(&Reply{})
	Engine.CreateTables(&Kvs{})
	Engine.CreateTables(&File{})
	Engine.CreateTables(&Timeline{})

	defer Engine.Close()
}

func SetEngine() {
	Engine, _ = xorm.NewEngine("postgres", "user=postgres password=jn!@#$%^&* dbname=pgsql sslmode=disable")
}

func GetEngine() *xorm.Engine {
	return Engine
}

func ConDb() *xorm.Engine {

	/*
		Engine, _ = xorm.NewEngine("sqlite3", "./test.db")
	*/
	Engine, _ := xorm.NewEngine("postgres", "user=postgres password=jn!@#$%^&* dbname=pgsql sslmode=disable")
	return Engine
}

func GetTopic(id int64) (tp Topic) {

	Engine := GetEngine()
	tp.Id = id

	Engine.Get(&tp)

	return tp
}

func GetAllTopic(offset int, limit int, path string) (allt []Topic) {

	Engine := GetEngine()

	Engine.Limit(limit, offset).OrderBy(path + " desc").Find(&allt)

	//q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("created").FindAll(&allt)
	return allt
}

func PostTopic(tp Topic) (int64, error) {

	Engine := GetEngine()
	id, err := Engine.Insert(&tp)

	return id, err
}

func PutTopic(tid int64, tp Topic) error {
	/*
				user := User{Name:"xlw"}
		rows, err := Engine.Update(&user, &User{Id:1})
		// rows, err := Engine.Where("id = ?", 1).Update(&user)
		// or rows, err := Engine.Id(1).Update(&user)
	*/
	Engine := GetEngine()
	_, err := Engine.Update(&tp, &Topic{Id: tid})
	return err

}

func DelTopic(id int64) error {

	Engine := GetEngine()

	topic := new(Topic)
	Engine.Id(id).Get(topic)

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

		if _, err := Engine.Delete(topic); err != nil {

			fmt.Println(err)
		} else {
			return err
		}

	}
	return errors.New("无法删除不存在的TOPIC ID:" + string(id))
}

/*
func main() {

	//Engine.ShowSQL = true

	var tp = Topic{Title: " haha!"}
	PostTopic(tp)
	for i := 0; i < 3; i++ {
		fmt.Println(GetTopic(int64(i)))
	}
}
*/
