package routers

import (
	"EasyGo/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.MainController{})
	web.Router("/report/mileage", &controllers.ReportController{}, "get:Mileage")
	web.Router("/report/mileage/multiple", &controllers.ReportController{}, "get:MultipleMileage")
	web.Router("/firebasefile", &controllers.FirebaseFileController{}, "post:Post;get:GetAll")
	web.Router("/firebasefile/multiple", &controllers.FirebaseFileController{}, "post:PostMulti;get:GetAll")
	web.Router("/firebasefile/multiple/:id", &controllers.FirebaseFileController{}, "put:PutMulti;get:GetAll")
	web.Router("/firebasefile/:id", &controllers.FirebaseFileController{}, "get:GetOne;put:Put")
	web.Router("/firebaselink", &controllers.FirebaseFileController{}, "post:PostLink")
}
