package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Estado     string `json:"estado"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go <CEP> <WEATHER_API_KEY>")
		return
	}

	cep := os.Args[1]
	apiKey := os.Args[2]

	// Consulta ViaCEP
	viaCEPurl := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(viaCEPurl)
	if err != nil {
		fmt.Println("Erro ao consultar ViaCEP:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler o corpo da resposta: %v\n", err)
		return
	}

	var cepData ViaCEPResponse
	if err := json.Unmarshal(bodyBytes, &cepData); err != nil {
		fmt.Println("Erro ao decodificar JSON do ViaCEP:", err)
		return
	}

	fmt.Printf("Localidade: %s, %s\n", cepData.Localidade, cepData.Estado)

	// Consulta WeatherAPI
	weatherAPIurl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, cepData.Localidade)
	resp, err = http.Get(weatherAPIurl)
	if err != nil {
		fmt.Println("Erro ao consultar WeatherAPI:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler o corpo da resposta: %v\n", err)
		return
	}

	var weatherData WeatherAPIResponse
	if err := json.Unmarshal(bodyBytes, &weatherData); err != nil {
		fmt.Println("Erro ao decodificar JSON do WeatherAPI:", err)
		return
	}

	tempC := weatherData.Current.TempC
	tempF := tempC*1.8 + 32 // Conversão de Celsius para Fahrenheit
	tempK := tempC + 273    // Conversão de Celsius para Kelvin

	fmt.Printf("Temperatura atual em %s: %.1f°C\n", cepData.Localidade, tempC)
	fmt.Printf("Temperatura atual em Fahrenheit: %.1f°F\n", tempF)
	fmt.Printf("Temperatura atual em Kelvin: %.1fK\n", tempK)
}
