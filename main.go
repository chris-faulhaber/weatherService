package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ErrResponse struct {
	Cod        string   `json:"cod"`
	Message    string   `json:"message"`
	Parameters []string `json:"parameters"`
}

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

var apiKey string
var endPoint = "https://api.openweathermap.org/data/2.5/weather"

func main() {
	apiKey = os.Getenv("WEATHER_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("Required API key not set in environment variable WEATHER_API_KEY")
	}

	http.HandleFunc("/weather", weatherHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	lat, ok := getParam(w, r, "lat")
	if !ok {
		return
	}

	long, ok := getParam(w, r, "long")
	if !ok {
		return
	}

	request := fmt.Sprintf("%s?units=imperial&lat=%s&lon=%s&APPID=%s", endPoint, lat, long, apiKey)
	log.Println(request)

	resp, err := http.Get(request)
	if err != nil {
		returnErr(w, fmt.Sprintf("Couldn't get data from Open Weather API {%s}\n", err.Error()))
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var errResponse ErrResponse
		_ = json.Unmarshal(body, &errResponse)
		returnErr(w, fmt.Sprintf("return code not 200, code {%d} - {%s}", resp.StatusCode, errResponse.Message))
		return
	}

	var weatherResponse WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		returnErr(w, fmt.Sprintf("problem unmarshaling {%s}", err.Error()))
		return
	}
	log.Print(weatherResponse)

	var tempCondition string
	if weatherResponse.Main.Temp < 33 {
		tempCondition = "cold"
	} else if weatherResponse.Main.Temp > 79 {
		tempCondition = "hot"
	} else {
		tempCondition = "moderate"
	}

	response := fmt.Sprintf("The weather condition outside is %s and the tempature is %s", weatherResponse.Weather[0].Description, tempCondition)
	_, err = w.Write([]byte(response))
	if err != nil {
		returnErr(w, fmt.Sprintf("Couldn't write response {%s}\n", err.Error()))
		return
	}

}

func returnErr(w http.ResponseWriter, msg string) {
	log.Print(msg)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(msg))
}

func getParam(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Url Param '%s' is missing", key)
		_, _ = w.Write([]byte(msg))
		return "", ok
	}
	return keys[0], ok

}
