package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BrazilAPIAddress struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type ViaCepAddress struct {
	Cep          string `json:"cep"`
	State        string `json:"uf"`
	City         string `json:"localidade"`
	Neighborhood string `json:"bairro"`
	Street       string `json:"logradouro"`
}

const (
	brasilApiUrl = "https://brasilapi.com.br/api/cep/v1/01153000"
	viaCepUrl    = "https://viacep.com.br/ws/01153000/json/"
)

func main() {

	ch1 := make(chan BrazilAPIAddress)
	ch2 := make(chan ViaCepAddress)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()

	go getAddressFromBrazilAPI(ctx, ch1)
	go getAddressFromViaCep(ctx, ch2)

	select {
	case brAPImsg := <-ch1:
		println("Brazil API:")
		fmt.Println(brAPImsg)

	case viaCepMsg := <-ch2:
		println("Via Cep:")
		fmt.Println(viaCepMsg)

	case <-ctx.Done():
		println("Timeout")
	}

}

func getAddressFromBrazilAPI(ctx context.Context, ch chan<- BrazilAPIAddress) {
	var address BrazilAPIAddress
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, brasilApiUrl, nil)

	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&address)
	if err != nil {
		panic(err)
	}
	ch <- address
}

func getAddressFromViaCep(ctx context.Context, ch chan<- ViaCepAddress) {
	var address ViaCepAddress
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, viaCepUrl, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&address)
	if err != nil {
		panic(err)
	}
	ch <- address
}
