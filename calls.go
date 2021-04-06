package main

import (
	"fmt"
	"log"
	"os"

	"github.com/VictorAvelar/mollie-api-go/v2/mollie"
	"github.com/gin-gonic/gin"
)

var (
	config *mollie.Config
	client *mollie.Client
)

func init() {
	_ = os.Setenv(mollie.APITokenEnv, "TOKEN")

	config := mollie.NewConfig(true, mollie.APITokenEnv)
	client, _ = mollie.NewClient(nil, config)
}

func main() {
	r := gin.Default()
	r.GET("/payment", call)

	r.Run()
}

func call(c *gin.Context) {
	pp := workflow()

	c.JSON(200, gin.H{
		"customerId": pp.ID,
		"Links":      pp.Links.Checkout.Href,
	})
}

func workflow() *mollie.Payment {

	m := CreateMollie()

	cs := CreateCustomer(m, "barbare", "code@yt.barb")

	cs = GetCustomer(m, cs.ID)

	var p = mollie.Payment{
		Amount: &mollie.Amount{
			Currency: "EUR",
			Value:    "1.00",
		},
		Description: "Mon paiement",
		RedirectURL: "http://localhost:3000",
		WebhookURL:  "https://webshop.example.org/payments/webhook/",
	}

	pp, err := m.Customers.CreatePayment(cs.ID, p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New payment ID %s, currency : %s, "+
		"value : %s , checkout : %s",
		pp.ID, pp.Amount.Currency, pp.Amount.Value, pp.Links.Checkout.Href)

	return pp
}

func CreateMollie() *mollie.Client {
	config := mollie.NewConfig(false, mollie.APITokenEnv)
	m, err := mollie.NewClient(nil, config)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func CreateCustomer(m *mollie.Client, name string, email string) *mollie.Customer {
	var c = mollie.Customer{
		Name:  name,
		Email: email,
	}
	cs, err := m.Customers.Create(c)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID %s\n", cs.ID)
	return cs
}

func GetCustomer(m *mollie.Client, customerId string) *mollie.Customer {

	cs, err := m.Customers.Get(customerId)
	if err != nil {
		log.Fatal(err)
	}
	return cs
}
