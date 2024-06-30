package main

import (
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

	go func() {
		var address BrazilAPIAddress
		req, err := http.Get(brasilApiUrl)

		if err != nil {
			panic(err)
		}
		defer req.Body.Close()

		err = json.NewDecoder(req.Body).Decode(&address)
		if err != nil {
			panic(err)
		}
		ch1 <- address
	}()

	go func() {
		var address ViaCepAddress
		req, err := http.Get(viaCepUrl)

		if err != nil {
			panic(err)
		}
		defer req.Body.Close()

		err = json.NewDecoder(req.Body).Decode(&address)
		if err != nil {
			panic(err)
		}
		ch2 <- address
	}()

	select {
	case brAPImsg := <-ch1:
		println("Brazil API:")
		fmt.Println(brAPImsg)

	case viaCepMsg := <-ch2:
		println("Via Cep:")
		fmt.Println(viaCepMsg)

	case <-time.After(5 * time.Second):
		println("Timeout")
	}

}
