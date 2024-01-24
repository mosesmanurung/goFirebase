// package main

// import (
// 	_ "EasyGo/routers"
// 	"github.com/beego/beego/v2/server/web"
// )

// func main() {
// 	web.Run()
// }

package main

import (
	// "EasyGo/models"
	_ "EasyGo/routers"

	"github.com/astaxie/beego/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// // Inisialisasi database
	// err := orm.RegisterDataBase("default", "mysql", "root:Sesmo@12345@tcp(127.0.0.1:3306)/firebase?charset=utf8")
	// if err != nil {
	// 	panic(err)
	// }

	// // Load models
	// orm.RegisterModel(new(models.FirebaseFile))

	// // Auto create table
	// err = orm.RunSyncdb("default", false, true)
	// if err != nil {
	// 	panic(err)
	// }

	// web.Run()

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:Sesmo@12345@tcp(127.0.0.1:3306)/firebase?charset=utf8")

	web.Run()
}
