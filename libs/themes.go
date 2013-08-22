package libs

import (
	"github.com/insionng/veryhour/plugins/goconfig"
	"github.com/astaxie/beego"
)

func init() {
	themes := GetThemes()
	beego.SetStaticPath("/static", themes+"/static")
	beego.SetViewsPath(themes + "/views")
}

func GetThemes() string {
	var themes string

	if conf, err := goconfig.LoadConfigFile("./conf/config.conf"); err != nil {

		//如果不存在配置文件，则设为默认主题路径
		themes = "./themes/default"
	} else { //如果存在配置文件

		// 主题设置读取错误 即section不存在 或 字段为空 则重置主题为默认主题并保存到配置文件
		if themes, err = conf.GetValue("themes", "path"); err != nil {
			conf.SetValue("themes", "path", "./themes/default")
			goconfig.SaveConfigFile(conf, "./conf/config.conf")
			themes = "./themes/default"
		}
	}
	return themes
}
