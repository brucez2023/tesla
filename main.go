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

	// Parse the body to JSON
	err = context.BindJSON(&jsonBody)
	if err != nil {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// Use the delimiter of : to split up the body it should have 4 parts
	parts := strings.Split(jsonBody.Data, ":")
	if len(parts) != 4 {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	
	// First check if we can cover the first part into integer
	newDevice.deviceId, err = strconv.Atoi(parts[0])
	if err != nil {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// Second part needs to be checked if we can convert into 64bit integer
	newDevice.time, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// Not sure why this is the case, but depending on terminal or postman tool to call the curl command it can ignore the single quotes
	// Because of this have to hand the case of where there is either a single quote or no single quote.
	// Third part needs to check if the string is 'Temperature'
	if !(parts[2] == "Temperature" || parts[2] == "\'Temperature\'") {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// Fourth Part needs to make sure it's a 64bit float
	newDevice.temperature, err = strconv.ParseFloat(parts[3], 64)
	if err != nil {
		errors = append(errors, jsonBody.Data)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	tempFormat(context, newDevice)

}

// tempFormat formats the json properly back
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
