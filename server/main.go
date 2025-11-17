package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CotacaoResponse struct {
	USDBRL Cotacao `json:"USDBRL"`
}

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
		RegisterLog(err.Error())
		if throwPanic {
			panic(err)
		}
	}
}

func RegisterLog(message string) {
	log.Println("LOG - ", message, " | ", time.DateTime)
}

func InternalServerError(err error, message string, w http.ResponseWriter) {
	if err != nil {
		RegisterLog(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(message))
	}
}

func main() {
	handler := CotacaoHandler{}.Init("root:root@tcp(localhost:3306)/cotacao?charset=utf8mb4&parseTime=True&loc=Local")

	http.HandleFunc("/cotacao", handler.ConsultaCotacaoUSD)
	http.ListenAndServe(":8080", nil)
}

type CotacaoHandler struct {
	DB *gorm.DB
}

func (c CotacaoHandler) Init(dsn string) CotacaoHandler {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DefaultContextTimeout: time.Millisecond * 10})
	raisePanic(err, true)
	db.AutoMigrate(&Cotacao{})
	c.DB = db
	return c
}

func (c CotacaoHandler) ConsultaCotacaoUSD(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2000)
	defer cancel()

	http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	request, err := http.NewRequest("GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	InternalServerError(err, "Failed to create request", w)

	response, err := client.Do(request)
	InternalServerError(err, "Failed to perform request", w)

	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	InternalServerError(err, "Failed to read response body", w)

	var cotacao CotacaoResponse
	err = json.Unmarshal(content, &cotacao)
	InternalServerError(err, "Failed to create request", w)

	c.DB.Create(&cotacao.USDBRL)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao.USDBRL)
}
