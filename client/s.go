package main

import (
	"../utils"
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

/*
package main

import "time"

func work(ch chan int) {  //工作方法
    //做事
    <- ch  //完事后在 ch里面抽一个数。
}

func main() {
    ch := make(chan int, 100)  //限制100个goroutine
    for i := 0; ; i++ {
        ch <- 1              //如果当有已经有100个goroutine在工作，那么此事就会堵塞直到有goroutine工作完
        go work(ch)
    }
}
*/

func work(ch chan int, i int) { //工作方法
	//做事
	fmt.Println("#", i)
	TPost()
	TPut(i)
	TDel(i)
	<-ch //完事后在 ch里面抽一个数。
}

func main() {
	ch := make(chan int, 100) //限制100个goroutine

	for i := 0; ; i++ {
		ch <- 1 //如果当有已经有100个goroutine在工作，那么此事就会堵塞直到有goroutine工作完
		go work(ch, i)
	}
}

func TPost() {

	content := []byte(`{"Cid":1,"Nid":5,"Uid":10,"Ctype":8,"Title":"Title防止擅改！itleTitle","Content":"IFuckin加密后gYou!!!!!"}`)
	actionurl := "http://localhost/core/topic"
	if sco, err := utils.SendingPackets(true, "POST", actionurl, content, publicKey); err != nil || sco.StatusCode != 200 {
		fmt.Println("Post发送失败！")
		fmt.Println(sco)
	} else {
		fmt.Println("Post发送成功！")
		fmt.Println(sco)
	}
}

func TPut(i int) {

	content := []byte(`{"Cid":2,"Nid":4,"Uid":8,"Ctype":88,"Title":"PutTitle防止擅改PutPut","Content":"PutPut加密后FuckingYou!!!!!"}`)
	actionurl := "http://localhost/core/topic/" + strconv.Itoa(i)

	if sco, err := utils.SendingPackets(true, "PUT", actionurl, content, publicKey); err != nil || sco.StatusCode != 200 {
		fmt.Println("Put发送失败！")
		fmt.Println(sco)
	} else {
		fmt.Println("Put发送成功！")
		fmt.Println(sco)
	}
}

func TDel(i int) {

	content := []byte(`{"Id":` + strconv.Itoa(i) + `}`)
	actionurl := "http://localhost/core/topic/" + strconv.Itoa(i)
	fmt.Println(string(content))
	fmt.Println(actionurl)
	if sco, err := utils.SendingPackets(true, "DELETE", actionurl, content, publicKey); err != nil || sco.StatusCode != 200 {
		fmt.Println("Del发送失败！")
		fmt.Println(sco)
	} else {
		//	Delete success!

		fmt.Println("Del发送成功！")
		fmt.Println(sco.Body)

	}
}
