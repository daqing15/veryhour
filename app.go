package main

import (
	"github.com/insionng/veryhour/core"
	"github.com/insionng/veryhour/handlers"
	//"./handlers/root"
	"github.com/astaxie/beego"
	//"runtime"
)

func main() {

	beego.SetStaticPath("/doc", "./doc")
	beego.Router("/", &handlers.MainHandler{})
	/*
		beego.Router("/topic/:tid([0-9]+)", &handlers.TopicHandler{})

		beego.Router("/node/:nid([0-9]+)", &handlers.NodeHandler{})

		beego.Router("/topic/delete/:tid([0-9]+)", &handlers.TopicDeleteHandler{})
		beego.Router("/topic/edit/:tid([0-9]+)", &handlers.TopicEditHandler{})

		beego.Router("/category/:cid([0-9]+)", &handlers.MainHandler{})
		beego.Router("/search", &handlers.SearchHandler{})

		beego.Router("/register", &handlers.RegHandler{})
		beego.Router("/login", &handlers.LoginHandler{})
		beego.Router("/logout", &handlers.LogoutHandler{})

		beego.Router("/like/topic/:tid([0-9]+)", &handlers.LikeTopicHandler{})
		beego.Router("/hate/topic/:tid([0-9]+)", &handlers.HateTopicHandler{})

		beego.Router("/like/node/:nid([0-9]+)", &handlers.LikeNodeHandler{})
		beego.Router("/hate/node/:nid([0-9]+)", &handlers.HateNodeHandler{})

		beego.Router("/new/category", &handlers.NewCategoryHandler{})
		beego.Router("/new/node", &handlers.NewNodeHandler{})
		beego.Router("/new/topic", &handlers.NewTopicHandler{})
		beego.Router("/new/reply/:tid([0-9]+)", &handlers.NewReplyHandler{})

		beego.Router("/modify/category", &handlers.ModifyCategoryHandler{})
		beego.Router("/modify/node", &handlers.ModifyNodeHandler{})

		beego.Router("/node/delete/:nid([0-9]+)", &handlers.NodeDeleteHandler{})
		beego.Router("/node/edit/:nid([0-9]+)", &handlers.NodeEditHandler{})

		beego.Router("/delete/reply/:rid([0-9]+)", &handlers.DeleteReplyHandler{})
	*/
	beego.Router("/core/topic", &core.TopicHandler{})
	//beego.Router("/core/node", &core.NodeHandler{})
	/*
		beego.Router("/root", &root.RMainHandler{})
		beego.Router("/root-login", &root.RLoginHandler{})
		beego.Router("/root/account", &root.RAccountHandler{})
	*/
	beego.SessionOn = true
	//runtime.GOMAXPROCS(2)
	beego.Run()
}
