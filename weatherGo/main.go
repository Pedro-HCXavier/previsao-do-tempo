package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"log"
	"net/url"
)

type WeatherResponse struct {
	Weather []struct {
		Descricao string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temperatura float64 `json:"temp"`
		Sensacao float64 `json:"feels_like"`
		Maxima float64 `json:"temp_max"`
		Minima float64 `json:"temp_min"` 
	} `json:"main"`
}

type GeoResponse struct {
	Nome string `json:"name"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// retorna respectivamente, cidade, lat e lon
func geolocalizacao(cidade string) (string, float64, float64) {
	cidadeCodificada := url.QueryEscape(cidade)
	apiKey := "ce4b16a8bc0e9a5d23c1cb7934123213"
	apiUrl := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s", cidadeCodificada, apiKey)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatal("status da requisição:", res.StatusCode)
	}


	var localizacao []GeoResponse
	if err := json.NewDecoder(res.Body).Decode(&localizacao) ; err != nil {
		log.Fatal("algo deu errado durante a requisição da geolocalização", err)
	}
	if len(localizacao) > 0 {
		return localizacao[0].Nome, localizacao[0].Lat, localizacao[0].Lon
	}

	return "",0,0
}

func descreverClima(weatherData WeatherResponse) string {	
	if len(weatherData.Weather) > 0 {
		return weatherData.Weather[0].Descricao
	}
	return "Sem descrição"
}

func main() {
	cidade, lat, lon := geolocalizacao("São Paulo")
	


	apiKey := "8a39857feed93b77ffda50940099f2c6"
	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric&lang=pt", lat, lon, apiKey)

	res, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("Algo deu errado durante a requisição, statusCode:%d\n", res.StatusCode)
		return
	}

	var weather WeatherResponse
	if err := json.NewDecoder(res.Body).Decode(&weather) ; err != nil {
		fmt.Println("erro ao decodificar o body", err)
	}

	descricao := descreverClima(weather)
	fmt.Printf("Clima em %s:\n%s\n%v Cº\nSensação de: %v Cº\nMáxima: %v Cº\nMínima: %v Cº\n", cidade, descricao, weather.Main.Temperatura, weather.Main.Sensacao, weather.Main.Maxima, weather.Main.Minima)
}