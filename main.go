package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Estado     string `json:"estado"`
	Erro       bool   `json:"erro"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func main() {
	http.HandleFunc("/clima", climaHandler)
	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}

func climaHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	apiKey := r.URL.Query().Get("apiKey")

	// Validação do CEP (sempre 8 dígitos numéricos)
	if !regexp.MustCompile(`^\d{8}$`).MatchString(cep) {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	// Consulta ViaCEP
	viaCEPurl := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(viaCEPurl)
	if err != nil {
		http.Error(w, `{"message": "erro ao consultar ViaCEP"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var cepData ViaCEPResponse
	if err := json.Unmarshal(body, &cepData); err != nil || cepData.Erro {
		http.Error(w, `{"message": "can not find zipcode"}`, http.StatusInternalServerError)
		return
	}

	// Consulta WeatherAPI
	weatherAPIurl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, cepData.Localidade)
	resp, err = http.Get(weatherAPIurl)
	if err != nil {
		http.Error(w, `{"message": "erro ao consultar WeatherAPI"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	var weatherData WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		http.Error(w, `{"message": "erro ao processar dados do clima"}`, http.StatusInternalServerError)
		return
	}

	tempC := weatherData.Current.TempC
	tempF := tempC*1.8 + 32 // Conversão de Celsius para Fahrenheit
	tempK := tempC + 273    // Conversão de Celsius para Kelvin

	// Resposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]float64{
		"tempC": tempC,
		"tempF": tempF,
		"tempK": tempK,
	})
}
