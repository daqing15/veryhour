package handlers

import (
	"github.com/insionng/veryhour/libs"
	//"../models"
	//"../utils"
)

type MainHandler struct {
	libs.BaseHandler
}

func (self *MainHandler) Get() {

	self.TplNames = "index.html"

}
