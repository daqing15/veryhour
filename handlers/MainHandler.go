package handlers

import (
	"github.com/insionng/veryhour/libs"
	//"../models"
	"github.com/insionng/veryhour/utils"
)

type MainHandler struct {
	libs.BaseHandler
}

func (self *MainHandler) Get() {

	self.TplNames = "default.html"
	content, _ := self.BaseHandler.RenderString()
	utils.WriteFile("./doc/", "default.html", content)
	self.Redirect("/doc/default.html", 302)

}
