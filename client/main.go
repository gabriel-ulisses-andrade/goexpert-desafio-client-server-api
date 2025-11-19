package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	ID         int32  `gorm:"primaryKey;autoIncrement"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func raisePanic(err error, throwPanic bool) {
	if err != nil {
		log.Println("LOG - ", err, " | ", time.DateTime)
		if throwPanic {
			panic(err)
		}
	}
}

func main() {
	cotacao, err := ConsultarCotacao()
	raisePanic(err, true)
	CriarArquivo(cotacao)
}

func CriarArquivo(cotacao Cotacao) {
	f, err := os.Create("cotacao.txt")
	raisePanic(err, true)
	_, err = f.WriteString("Dólar: " + cotacao.Bid)
	raisePanic(err, true)
	defer f.Close()
}

func ConsultarCotacao() (Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return Cotacao{}, err
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return Cotacao{}, err
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return Cotacao{}, err
	}

	var cotacao Cotacao
	err = json.Unmarshal(content, &cotacao)
	if err != nil {
		return Cotacao{}, err
	}

	return cotacao, nil
}
