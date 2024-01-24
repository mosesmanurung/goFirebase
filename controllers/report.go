// controllers/report.go
package controllers

import (
	"EasyGo/models"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/client/httplib"
	"strings"
)

type ReportController struct {
	web.Controller
}

func (c *ReportController) Mileage() {
	startTime := c.GetString("start_time")
	stopTime := c.GetString("stop_time")
	carPlate := c.GetString("nopol")

	if carPlate == "" {
		c.Data["json"] = map[string]interface{}{"status": 400, "message": "Missing lstNoPOL parameter"}
		c.ServeJSON()
		return
	}

	url := "https://vtsapi.easygo-gps.co.id/api/report/total_km"
	request := httplib.Post(url)
	request.Header("token", "54DBF55357244438B965C9C3F6ECC325")
	request.JSONBody(map[string]interface{}{
		"start_time": startTime,
		"stop_time":  stopTime,
		"lstNoPOL":   []string{carPlate},
	})

	var response models.TotalKmResponse
	err := request.ToJSON(&response)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"status": 500, "message": "Internal Server Error"}
		c.ServeJSON()
		return
	}

	var responseData map[string]interface{}
	if response.ResponseCode == 0 {
		responseData = map[string]interface{}{"status": 200, "message": response.ResponseMessage}
	} else if response.ResponseCode == 1 {
		responseData = map[string]interface{}{"status": 200, "message": response.ResponseMessage, "data": response.Data[0]}
	} else {
		responseData = map[string]interface{}{"status": 500, "message": "Unknown ResponseCode"}
	}

	c.Data["json"] = responseData
	c.ServeJSON()
}

func (c *ReportController) MultipleMileage() {
	startTime := c.GetString("start_time")
	stopTime := c.GetString("stop_time")
	carPlates := c.GetString("nopol")

	// Split the comma-separated string into an array of strings
	carPlateArray := strings.Split(carPlates, ",")

	if len(carPlateArray) == 0 {
		c.Data["json"] = map[string]interface{}{"status": 400, "message": "Missing nopol parameter"}
		c.ServeJSON()
		return
	}

	url := "https://vtsapi.easygo-gps.co.id/api/report/total_km"
	request := httplib.Post(url)
	request.Header("token", "54DBF55357244438B965C9C3F6ECC325")
	request.JSONBody(map[string]interface{}{
		"start_time": startTime,
		"stop_time":  stopTime,
		"lstNoPOL":   carPlateArray,
	})

	var response models.TotalKmResponse
	err := request.ToJSON(&response)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"status": 500, "message": "Internal Server Error"}
		c.ServeJSON()
		return
	}

	var responseData map[string]interface{}
	if response.ResponseCode == 0 {
		responseData = map[string]interface{}{"status": 200, "message": response.ResponseMessage}
	} else if response.ResponseCode == 1 {
		responseData = map[string]interface{}{"status": 200, "message": response.ResponseMessage, "data": response.Data}
	} else {
		responseData = map[string]interface{}{"status": 500, "message": "Unknown ResponseCode"}
	}

	c.Data["json"] = responseData
	c.ServeJSON()
}


// func (c *ReportController) Mileage() {
// 	startTime := c.GetString("start_time")
// 	stopTime := c.GetString("stop_time")
// 	carPlate := c.GetString("nopol")

// 	url := "https://vtsapi.easygo-gps.co.id/api/report/total_km"
// 	request := httplib.Post(url)
// 	request.Header("token", "54DBF55357244438B965C9C3F6ECC325")
// 	request.JSONBody(map[string]interface{}{
// 		"start_time": startTime,
// 		"stop_time":  stopTime,
// 		"lstNoPOL":   []string{carPlate},
// 	})

// 	var response models.TotalKmResponse
// 	err := request.ToJSON(&response)
// 	if err != nil {
// 		c.Data["json"] = map[string]interface{}{"status": 500, "message": "Internal Server Error"}
// 		c.ServeJSON()
// 		return
// 	}

// 	c.Data["json"] = map[string]interface{}{"status": 200, "message": response.ResponseMessage, "Data": response.Data[0]}
// 	c.ServeJSON()
// }


// func (c *ReportController) Mileage() {
// 	startTime := c.GetString("start_time")
// 	stopTime := c.GetString("stop_time")
// 	carPlate := c.GetString("nopol")

// 	url := "https://vtsapi.easygo-gps.co.id/api/report/total_km"
// 	request := httplib.Post(url)
// 	request.Header("token", "54DBF55357244438B965C9C3F6ECC325")
// 	request.JSONBody(map[string]interface{}{
// 		"start_time": startTime,
// 		"stop_time":  stopTime,
// 		"lstNoPOL":   []string{carPlate},
// 	})

// 	var response models.TotalKmResponse
// 	err := request.ToJSON(&response)
// 	if err != nil {
// 		c.Data["json"] = map[string]interface{}{"status": 500, "message": "Internal Server Error"}
// 		c.ServeJSON()
// 		return
// 	}

// 	var responseData map[string]interface{}
// 	if response.ResponseCode == 0 {
// 		responseData = map[string]interface{}{"status": 200, "message": response.ResponseMessage}
// 	} else if response.ResponseCode == 1 {
// 		responseData = map[string]interface{}{"status": 200, "message": response.ResponseMessage, "data": response.Data[0]}
// 	} else {
// 		responseData = map[string]interface{}{"status": 500, "message": "Unknown ResponseCode"}
// 	}

// 	c.Data["json"] = responseData
// 	c.ServeJSON()
// }