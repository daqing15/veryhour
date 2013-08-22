package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

type User struct {
	Id      int `orm:"auto"` // 设置为auto主键
	Name    string
	Profile *Profile `orm:"rel(one)"` // OneToOne relation
}

type Profile struct {
	Id   int `orm:"auto"`
	Age  int16
	User *User `orm:"reverse(one)"` // 设置反向关系(可选)
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(User), new(Profile))

	//orm.RegisterDriver("postgres", orm.DR_Postgres)
	orm.RegisterDataBase("default", "postgres", "user=postgres password=jn!@#$%^&* dbname=pgsql sslmode=disable", 30)

}

func main() {
	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库

	profile := &Profile{}

	profile.Age = 30

	user := &User{}
	user.Profile = profile
	user.Name = "slene"

	fmt.Println(o.Insert(profile))
	fmt.Println(o.Insert(user))
}
