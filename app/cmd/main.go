package main

import (
	"context"
	"fmt"

	api "github.com/shun-ideguchi/golang-api-playground"
)

func main() {
	client, err := api.NewClient(api.WithApiKey("API_KEY"))
	if err != nil {
		fmt.Println(err.Error())
	}
	bank, err := client.GetBank(context.Background(), "0001", &api.GetParameter{})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(bank)
	}
}
