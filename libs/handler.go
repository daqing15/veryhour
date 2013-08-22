package libs

import (
	//"../models"
	"github.com/insionng/veryhour/utils"
	"github.com/astaxie/beego"
	"runtime"
//	"time"
)

var (
	sess_username string
	sess_uid      int64
	sess_role     int64
	sess_email    string
	sess_core     bool
	bc            *beego.BeeCache
)

type BaseHandler struct {
	beego.Controller
}

type AuthHandler struct {
	BaseHandler
}

type RootAuthHandler struct {
	BaseHandler
}

type CoreAuthHandler struct {
	BaseHandler
}

type RootHandler struct {
	BaseHandler
}

func init() {
	bc = beego.NewBeeCache()
	bc.Every = 259200 //該單位為秒，0為不過期，259200 三天,604800 即一個星期清空一次緩存
	bc.Start()
}

//用户等级划分：正数是普通用户，负数是管理员各种等级划分，为0则尚未注册
func (self *BaseHandler) Prepare() {
	sess_username, _ = self.GetSession("username").(string)
	sess_uid, _ = self.GetSession("userid").(int64)
	sess_role, _ = self.GetSession("userrole").(int64)
	sess_email, _ = self.GetSession("useremail").(string)

	if sess_role == 0 {
		self.Data["Userid"] = 0
		self.Data["Username"] = ""
		self.Data["Userrole"] = 0
		self.Data["Useremail"] = ""
	} else {
		self.Data["Userid"] = sess_uid
		self.Data["Username"] = sess_username
		self.Data["Userrole"] = sess_role
		self.Data["Useremail"] = sess_email
	}
	/*
		self.Data["categorys"] = models.GetAllCategory()
		self.Data["nodes"] = models.GetAllNode()
		self.Data["topics_5s"] = models.GetAllTopic(0, 5, "id")
		self.Data["topics_10s"] = models.GetAllTopic(0, 10, "id")
		self.Data["nodes_10s"] = models.GetAllNodeByCid(0, 0, 10, 0, "id")
		self.Data["replys_5s"] = models.GetReplyByPid(0, 0, 5, "id")
		self.Data["replys_10s"] = models.GetReplyByPid(0, 0, 10, "id")

		self.Data["author"] = models.GetKV("author")
		self.Data["title"] = models.GetKV("title")
		self.Data["title_en"] = models.GetKV("title_en")
		self.Data["keywords"] = models.GetKV("keywords")
		self.Data["description"] = models.GetKV("description")

		self.Data["company"] = models.GetKV("company")
		self.Data["copyright"] = models.GetKV("copyright")
		self.Data["site_email"] = models.GetKV("site_email")

		self.Data["tweibo"] = models.GetKV("tweibo")
		self.Data["sweibo"] = models.GetKV("sweibo")
		self.Data["timenow"] = time.Now()
		self.Data["statistics"] = models.GetKV("statistics")
	*/

}

//会员或管理员前台权限认证
func (self *AuthHandler) Prepare() {
	self.BaseHandler.Prepare()

	if sess_role == 0 {
		self.Ctx.Redirect(302, "/login")
	}
}

//管理员前台权限认证
func (self *RootAuthHandler) Prepare() {
	self.BaseHandler.Prepare()
	if sess_role != -1000 {
		self.Ctx.Redirect(302, "/login")
	}
}

//CORE API认证
func (self *CoreAuthHandler) Prepare() {

	self.BaseHandler.Prepare()
	sess_core, _ = self.GetSession("core").(bool)
	if sess_role != -1000 {
		self.Ctx.Redirect(302, "/login")
	}
}

//管理员后台后台认证
func (self *RootHandler) Prepare() {
	self.BaseHandler.Prepare()

	if !utils.IsSpider(self.Ctx.Request.UserAgent()) {
		if sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			self.Data["remoteproto"] = self.Ctx.Request.Proto
			self.Data["remotehost"] = self.Ctx.Request.Host
			self.Data["remoteos"] = runtime.GOOS
			self.Data["remotearch"] = runtime.GOARCH
			self.Data["remotecpus"] = runtime.NumCPU()
			self.Data["golangver"] = runtime.Version()
		}
	} else {
		self.Ctx.Redirect(302, "/")
	}
}
