package models

type TotalKmResponse struct {
    ResponseCode    int     `json:"ResponseCode"`
    ResponseMessage string  `json:"ResponseMessage"`
    Data            []TotalKmData `json:"Data"`
}

type TotalKmData struct {
    CarPlate string  `json:"car_plate"`
    TotalKm  float64 `json:"total_km"`
}
