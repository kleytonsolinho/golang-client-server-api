package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/kleytonsolinho/golang-client-server-api/shared"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(shared.CurrencyPriceRoute, getCurrencyPriceHandler)
	http.ListenAndServe(":8080", mux)
}

func getCurrencyPriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != shared.CurrencyPriceRoute {
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

func getCurrency() (*shared.CurrencyPairs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", shared.AwesomeAPI, nil)
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

	var result shared.CurrencyPairs
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func currencyPriceRepository(currency *shared.CurrencyPairs) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/currency-db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	currencyPrice := NewCurrencyPrice(currency.Usdbrl.Bid)
	err = insertCurrencyPrice(ctx, db, currencyPrice)
	if err != nil {
		panic(err)
	}
}

func insertCurrencyPrice(ctx context.Context, db *sql.DB, currencyPrice *shared.CurrencyPrice) error {
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

func NewCurrencyPrice(price string) *shared.CurrencyPrice {
	return &shared.CurrencyPrice{
		ID:    uuid.New().String(),
		Code:  "USDBRL",
		Price: price,
	}
}
