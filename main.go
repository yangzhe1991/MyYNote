package main

import (
	"MyYNote/controllers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Router("/note/", &controllers.MainController{})
	beego.Router("/note/callback", &controllers.CallbackController{})
	beego.Router("/note/json/*", &controllers.JsonController{})
	beego.Router("/note/latex/*", &controllers.LatexController{})
	beego.Router("/note/latex/", &controllers.LatexController{})
	beego.Run()
}
