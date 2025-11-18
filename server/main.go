package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
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
		RegistrarLog(err.Error())
		if throwPanic {
			panic(err)
		}
	}
}

func RegistrarLog(message string) {
	log.Println("LOG - ", message, " | ", time.DateTime)
}

func InternalServerError(err error, message string, w http.ResponseWriter) {
	if err != nil {
		RegistrarLog(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(message))
	}
}

func main() {
	handler := CotacaoHandler{}.Init("cotacao.db")

	http.HandleFunc("/cotacao", handler.ConsultaCotacaoUSD)
	http.ListenAndServe(":8080", nil)
}

type CotacaoHandler struct {
	DB *gorm.DB
}

func (c CotacaoHandler) Init(dsn string) CotacaoHandler {
	db, err := gorm.Open(sqlite.Open(dsn) /*, &gorm.Config{DefaultContextTimeout: time.Millisecond * 10}*/)
	raisePanic(err, true)
	db.AutoMigrate(&Cotacao{})
	c.DB = db
	return c
}

func (c CotacaoHandler) ConsultaCotacaoUSD(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	InternalServerError(err, "Ocorreu um erro ao iniciar a requisição", w)

	response, err := client.Do(request)
	InternalServerError(err, "Ocorreu um erro ao realizar a requisição", w)

	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	InternalServerError(err, "Ocorreu um erro ao ler o conteúdo da resposta", w)

	var cotacao CotacaoResponse
	err = json.Unmarshal(content, &cotacao)
	InternalServerError(err, "Ocorreu um erro ao desserializar o conteúdo da resposta", w)

	message := c.SalvarDadosCotacao(cotacao.USDBRL)
	RegistrarLog(message)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao.USDBRL)
}

func (c CotacaoHandler) SalvarDadosCotacao(cotacao Cotacao) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()

	c.DB.WithContext(ctx).Create(&cotacao)

	<-ctx.Done()
	if cotacao.ID != 0 {
		return "Cotação " + string(cotacao.ID) + " registrada com sucesso!"
	}
	return "Não foi possível registrar a cotação. O tempo limite expirou para esta operação"
}
