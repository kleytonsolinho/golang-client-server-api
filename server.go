package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

	currency, err := getCurrency()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currency)
}

func getCurrency() (*shared.CurrencyPairs, error) {
	res, err := http.Get(shared.AwesomeAPI)
	if err != nil {
		return nil, err
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
