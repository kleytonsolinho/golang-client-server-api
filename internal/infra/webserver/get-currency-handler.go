package webserver

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kleytonsolinho/golang-client-server-api/internal/dto"
	"github.com/kleytonsolinho/golang-client-server-api/internal/entity"
	"github.com/kleytonsolinho/golang-client-server-api/internal/infra/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetCurrencyPriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := r.Context()
	log.Println("Request iniciada")
	defer log.Println("Request finalizada")

	select {
	case <-time.After(210 * time.Millisecond):
		log.Println("Request processada com sucesso")
	case <-ctx.Done():
		log.Println("Request cancelada pelo cliente")
	}

	currency, err := getCurrency()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	currencyPriceRepository(currency)
	log.Println("Dados salvos no banco de dados com sucesso")

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currency)
}

func getCurrency() (*dto.CurrencyPairs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error de timeout ao fazer a requisição")
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result dto.CurrencyPairs
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func currencyPriceRepository(currency *dto.CurrencyPairs) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.CurrencyPrice{})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	currencyPrice, err := entity.NewCurrencyPrice(currency.Usdbrl.Bid)
	if err != nil {
		panic(err)
	}

	currencyDB := database.NewCurrencyDb(db)

	err = currencyDB.Create(ctx, currencyPrice)
	if err != nil {
		panic(err)
	}
}
