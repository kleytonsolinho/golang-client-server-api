package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kleytonsolinho/golang-client-server-api/shared"
)

func client() {
	req, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var data shared.CurrencyPairs
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao fazer parse da resposta: %v\n", err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao criar arquivo: %v\n", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: {%s}", data.Usdbrl.Bid))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao escrever no arquivo: %v\n", err)
	}

	fmt.Println("Arquivo criado com sucesso!")
	fmt.Println("Response:" + string(res))
}
