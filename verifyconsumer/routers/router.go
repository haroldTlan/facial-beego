// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"verifyconsumer/controllers"
)

func init() {
	ns := beego.NewNamespace("/zbx",
		//can delete or TODO
		beego.NSNamespace("/compare",
			beego.NSInclude(
				&controllers.ComparedController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
