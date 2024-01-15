package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/kleytonsolinho/golang-client-server-api/internal/dto"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	currency, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data dto.CurrencyPairs
	err = json.Unmarshal(currency, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao fazer parse da resposta: %v\n", err)
	}

	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println("Erro ao carregar o fuso horário:", err)
		return
	}
	currentTime := time.Now().In(location)
	formattedTime := currentTime.Format("2006-01-02-15:04:05")

	file1, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao criar arquivo: %v\n", err)
	}
	defer file1.Close()

	file2, err := os.Create("tmp/cotacao-" + formattedTime + ".txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao criar arquivo: %v\n", err)
	}
	defer file2.Close()

	_, err = file1.WriteString(fmt.Sprintf("Dólar: {%s}", data.Usdbrl.Bid))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao escrever no arquivo: %v\n", err)
	}

	_, err = file2.WriteString(fmt.Sprintf("Dólar: {%s}", data.Usdbrl.Bid))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ao escrever no arquivo: %v\n", err)
	}

	fmt.Println("Arquivo criado com sucesso!")
	fmt.Println("Response:" + string(currency))
}
