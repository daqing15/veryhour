package core

import (
	"../libs"
	"../models"
	"encoding/json"
	"strconv"
)

type NodeHandler struct {
	libs.BaseHandler
}

//创建
func (self *NodeHandler) Post() {
	var tp models.Topic
	/*
		var b []byte
		b = []byte(`{"Cid":2,"Nid":4,"Uid":6,"Ctype":8,"Title":"TitleTitleTitleTitle","Content":"iamContentContentContentContent!!!!!!"}`)
		json.Unmarshal(b, &tp)
	*/
	//self.Ctx.RequestBody的内容必须是json格式
	json.Unmarshal(self.Ctx.RequestBody, &tp)
	if tid, err := models.PostTopic(tp); err != nil {

		self.Data["json"] = "Post failed!"
	} else {
		self.Data["json"] = `{"TopicId:"` + strconv.Itoa(int(tid)) + `}`
	}

	self.ServeJson()
}

//获取
func (self *NodeHandler) Get() {
	nid, _ := self.GetInt(":objectId") //beego api模式下，提交的参数名总是唤作objectId

	if nid > 0 {
		nd := models.GetNode(nid)

		self.Data["json"] = nd

	} else {
		nds := models.GetAllNode()
		self.Data["json"] = nds
	}
	self.ServeJson()
}

//更新
func (self *NodeHandler) Put() {
	tid, _ := self.GetInt(":objectId")
	var tp models.Topic
	json.Unmarshal(self.Ctx.RequestBody, &tp)

	if err := models.PutTopic(tid, tp); err != nil {
		self.Data["json"] = "Update failed!"
	} else {
		self.Data["json"] = "Update success!"
	}
	self.ServeJson()
}

//删除
func (self *NodeHandler) Delete() {

	tid, _ := self.GetInt(":objectId")

	if e := models.DelTopic(tid); e != nil {
		self.Data["json"] = "Delete failed!"
	} else {

		self.Data["json"] = "Delete success!"
	}
	self.ServeJson()
}
