package routers

import (
	"EasyGo/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.MainController{})
	web.Router("/report/mileage", &controllers.ReportController{}, "get:Mileage")
	web.Router("/report/mileage/multiple", &controllers.ReportController{}, "get:MultipleMileage")
	web.Router("/firebasefile", &controllers.FirebaseFileController{}, "post:Post")
	// web.Router("/firebasefile/multiple", &controllers.FirebaseFileController{}, "post:PostMulti")
}
