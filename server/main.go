package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kleytonsolinho/golang-client-server-api/internal/infra/webserver"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", webserver.GetCurrencyPriceHandler)
	http.ListenAndServe(":8080", mux)
}
