package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["verifyconsumer/controllers:ComparedController"] = append(beego.GlobalControllerRouter["verifyconsumer/controllers:ComparedController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["verifyconsumer/controllers:ComparedController"] = append(beego.GlobalControllerRouter["verifyconsumer/controllers:ComparedController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
