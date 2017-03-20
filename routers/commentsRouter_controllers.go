package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["cloud-platform-ua/controllers:UserController"] = append(beego.GlobalControllerRouter["cloud-platform-ua/controllers:UserController"],
		beego.ControllerComments{
			Method: "Register",
			Router: `/register`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["cloud-platform-ua/controllers:UserController"] = append(beego.GlobalControllerRouter["cloud-platform-ua/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["cloud-platform-ua/controllers:UserController"] = append(beego.GlobalControllerRouter["cloud-platform-ua/controllers:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}
