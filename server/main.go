package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/kleytonsolinho/golang-client-server-api/dto"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", getCurrencyPriceHandler)
	http.ListenAndServe(":8080", mux)
}

func getCurrencyPriceHandler(w http.ResponseWriter, r *http.Request) {
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
	db.AutoMigrate(&dto.CurrencyPrice{})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	currencyPrice := NewCurrencyPrice(currency.Usdbrl.Bid)
	currencyDB := NewCurrency(db)

	err = insertCurrencyPrice(ctx, currencyDB, currencyPrice)
	if err != nil {
		panic(err)
	}
}

type Currency struct {
	DB *gorm.DB
}

func NewCurrency(db *gorm.DB) *Currency {
	return &Currency{
		DB: db,
	}
}

func insertCurrencyPrice(ctx context.Context, db *gorm.DB, currencyPrice *dto.CurrencyPrice) error {
	stmt, err := db.Prepare("insert into currency_price(id, code, price) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, currencyPrice.ID, currencyPrice.Code, currencyPrice.Price)
	if err != nil {
		log.Println("Error ao salvar no banco de dados")
		return err
	}
	return nil
}

func NewCurrencyPrice(price string) *dto.CurrencyPrice {
	return &dto.CurrencyPrice{
		ID:    uuid.New().String(),
		Code:  "USDBRL",
		Price: price,
	}
}
