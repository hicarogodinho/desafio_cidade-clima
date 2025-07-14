package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClimaHandler(t *testing.T) {
	apiKey := "ff4e9f3ecf62466396a141841251407"

	tests := []struct {
		name           string
		cep            string
		expectedStatus int
	}{
		{
			name:           "CEP válido",
			cep:            "58700070",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "CEP inválido",
			cep:            "38750", // menos de 8 dígitos
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "CEP inexistente",
			cep:            "00000000", // CEP que não existe
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/clima?cep="+tt.cep+"&apiKey="+apiKey, nil)

			w := httptest.NewRecorder()
			climaHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("handler retornou o status %v, esperado %v", resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}
