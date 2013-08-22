package handlers

import (
	"../libs"
	"../models"
)

type TopicHandler struct {
	libs.BaseHandler
}

func (self *TopicHandler) Get() {
	tid, _ := self.GetInt(":tid")
	tid_handler := models.GetTopic(tid)
	self.TplNames = "view.html"

	if tid_handler.Id > 0 {

		tid_handler.Views = tid_handler.Views + 1
		models.UpdateTopic(tid, tid_handler)

		self.Data["article"] = tid_handler
		self.Data["replys"] = models.GetReplyByPid(tid, 0, 0, "id")

		tps := models.GetAllTopicByCid(tid_handler.Cid, 0, 0, 0, "asc")

		if tps != nil && tid != 0 {

			for i, v := range tps {

				if v.Id == tid {
					prev := i - 1
					next := i + 1

					for i, v := range tps {
						if prev == i {
							self.Data["previd"] = v.Id
							self.Data["prev"] = v.Title
						}
						if next == i {
							self.Data["nextid"] = v.Id
							self.Data["next"] = v.Title
						}
					}
				}
			}
		}

	} else {
		self.Redirect("/", 302)
	}

}
