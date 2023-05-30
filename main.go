package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)


type device struct{
	deviceId int
	time int64
	temperature float64
}

type inputBody struct{
	Data string `json:"data"`
}

var errors = []string{}

func main() {
	router := gin.Default()
	router.POST("/temp", postDevice)
	router.GET("/errors", getErrors)
	router.DELETE("/errors", clearBuffer)

	router.Run(":8080")
}

func postDevice(context *gin.Context) {

	var newDevice device
	var err error
	var jsonBody inputBody

	err = context.BindJSON(&jsonBody)
	if err != nil {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	parts := strings.Split(jsonBody.Data, ":")
	if len(parts) != 4 {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	
	newDevice.deviceId, err = strconv.Atoi(parts[0])
	if err != nil {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	newDevice.time, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if parts[2] != "Temperature" {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	newDevice.temperature, err = strconv.ParseFloat(parts[3], 64)
	if err != nil {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	tempFormat(context, newDevice)

}

func tempFormat(context *gin.Context, d device) {
	if d.temperature >= 90 {
		context.IndentedJSON(http.StatusOK,
			gin.H{"overtemp": true, "device_id": d.deviceId, "formatted_time": time.Unix(0, d.time * 1000000)})
	} else {
		context.IndentedJSON(http.StatusOK, gin.H{"overtemp": false})
	}
}

func getErrors(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, gin.H{"errors": errors})
}

func clearBuffer(context *gin.Context) {
	errors = []string{}
	context.IndentedJSON(http.StatusOK, gin.H{"msg": "errors deleted"})
}
