package core

import (
	"github.com/insionng/veryhour/libs"
	"github.com/insionng/veryhour/models"
	"github.com/insionng/veryhour/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

var (
	publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZsfv1qscqYdy4vY+P4e3cAtmv
ppXQcRvrF1cB4drkv0haU24Y7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0Dgacd
wYWd/7PeCELyEipZJL07Vro7Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NL
AUeJ6PeW+DAkmJWF6QIDAQAB
-----END PUBLIC KEY-----
`)

	privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`)
)

type TopicHandler struct {
	libs.BaseHandler
}

//创建帖子
func (self *TopicHandler) Post() {

	if hash := self.GetString("hash"); hash != "" {

		if rsa_decrypt_content, err := utils.ReceivingPackets(true, hash, "POST", self.Ctx.RequestBody, publicKey, privateKey); err == nil {
			var tp models.Topic
			json.Unmarshal(rsa_decrypt_content, &tp)
			if tid, err := models.PostTopic(tp); err != nil {

				self.Data["json"] = "Post failed!"
			} else {
				self.Data["json"] = `{"TopicId:"` + strconv.Itoa(int(tid)) + `}`
			}

			self.ServeJson()
		} else {

			fmt.Println("401 Unauthorized!")
			self.Abort("401")
		}
	} else {

		fmt.Println("401 Unauthorized!")
		self.Abort("401")
	}
}

//获取帖子
func (self *TopicHandler) Get() {
	tid, _ := self.GetInt(":objectId") //beego api模式下，提交的参数名总是唤作objectId

	if tid > 0 {
		tp := models.GetTopic(tid)

		self.Data["json"] = tp

	} else {
		tps := models.GetAllTopic(0, 0, "id")
		self.Data["json"] = tps
	}
	self.ServeJson()
}

//更新帖子
func (self *TopicHandler) Put() {

	if hash := self.GetString("hash"); hash != "" {

		if rsa_decrypt_content, err := utils.ReceivingPackets(true, hash, "PUT", self.Ctx.RequestBody, publicKey, privateKey); err == nil {

			tid, _ := self.GetInt(":objectId")
			var tp models.Topic
			json.Unmarshal(rsa_decrypt_content, &tp)

			if err := models.PutTopic(tid, tp); err != nil {
				self.Data["json"] = "Update failed!"
			} else {
				self.Data["json"] = "Update success!"
			}
			self.ServeJson()
		} else {

			fmt.Println("401 Unauthorized!")
			self.Abort("401")
		}
	} else {

		fmt.Println("401 Unauthorized!")
		self.Abort("401")
	}
}

//删除帖子
func (self *TopicHandler) Delete() {
	if hash := self.GetString("hash"); hash != "" {

		if rsa_decrypt_content, err := utils.ReceivingPackets(true, hash, "DELETE", self.Ctx.RequestBody, publicKey, privateKey); err == nil {

			tid, _ := self.GetInt(":objectId")
			var tp *models.Topic
			json.Unmarshal(rsa_decrypt_content, &tp)
			if tid == tp.Id && tid > 0 {
				if e := models.DelTopic(tid); e != nil {
					self.Data["json"] = "Delete failed!"

				} else {

					self.Data["json"] = "Delete success!"
				}
				//self.ServeJson()
				self.Ctx.WriteString(self.Data["json"].(string))
			} else {

				fmt.Println("401 Unauthorized!")
				self.Abort("401")
			}

		} else {

			fmt.Println("401 Unauthorized!")
			self.Abort("401")
		}
	} else {

		fmt.Println("401 Unauthorized!")
		self.Abort("401")
	}
}
