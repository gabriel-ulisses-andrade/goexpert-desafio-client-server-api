package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

var timeoutConsultaExterna = time.Millisecond * 2000
var timeoutEscritaDB = time.Millisecond * 10

func raisePanic(err error, throwPanic bool) {
	if err != nil {
		log.Println("ERROR - ", time.DateTime, " | ", err)
		if throwPanic {
			panic(err)
		}
	}
}

func InternalServerError(err error, message string, w http.ResponseWriter) {
	if err != nil {
		log.Println("ERROR - ", time.DateTime, " | ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(message))
	}
}

func main() {
	http.HandleFunc("/cotacao", ConsultaCotacaoUSD)
	http.ListenAndServe(":8080", nil)
}

func ConsultaCotacaoUSD(w http.ResponseWriter, r *http.Request) {
	client := http.Client{Timeout: timeoutConsultaExterna}

	request, err := http.NewRequest("GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	InternalServerError(err, "Failed to create request", w)

	response, err := client.Do(request)
	InternalServerError(err, "Failed to perform request", w)

	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	InternalServerError(err, "Failed to read response body", w)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
