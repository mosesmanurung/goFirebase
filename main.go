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

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		// AllowAllOrigins: true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "x-xsrf-token", "AxiosHeaders", "X-Requested-With", "X-CSRF-Token", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "Authorization", "Set-Cookie", "Cookie"},
		AllowCredentials: true,
	}))
	web.Run()
}
