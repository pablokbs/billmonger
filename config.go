package main

import (
	"io/ioutil"
	"regexp"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v2"
)

type BusinessDetails struct {
	Name      string `yaml:"name"`
	Person    string `yaml:"person"`
	Address   string `yaml:"address"`
	ImageFile string `yaml:"image_file"`
}

type BillDetails struct {
	Department   string `yaml:"department"`
	Currency     string `yaml:"currency"`
	PaymentTerms string `yaml:"payment_terms"`
	DueDate      string `yaml:"due_date"`
}

type BillToDetails struct {
	Email        string
	Name         string
	Street       string
	CityStateZip string `yaml:"city_state_zip"`
	Country      string
}

type BillableItem struct {
	Quantity    float64
	Description string
	UnitPrice   float64 `yaml:"unit_price"`
	Currency    string
}

func (b *BillableItem) Total() float64 {
	return b.UnitPrice * b.Quantity
}

func (b *BillableItem) Strings() []string {
	return []string{
		strconv.FormatFloat(b.Quantity, 'f', 2, 64),
		b.Description,
		b.Currency + " " + niceFloatStr(b.UnitPrice),
		b.Currency + " " + niceFloatStr(b.Total()),
	}
}

type BankDetails struct {
	TransferType string `yaml:"transfer_type"`
	Name         string
	Address      string
	AccountType  string `yaml:"account_type"`
	IBAN         string
	SortCode     string `yaml:"sort_code"`
}

func (b *BankDetails) Strings() []string {
	return []string{
		b.TransferType, b.Name, b.Address, b.AccountType, b.IBAN, b.SortCode,
	}
}

type BillingConfig struct {
	Business  *BusinessDetails `yaml:"business"`
	Bill      *BillDetails     `yaml:"bill"`
	BillTo    *BillToDetails   `yaml:"bill_to"`
	Billables []BillableItem   `yaml:"billables"`
	Bank      *BankDetails     `yaml:"bank"`
}

// ParseConfig parses the YAML config file which contains the
// settings for the bill we're going to process.
func ParseConfig(filename string) (*BillingConfig, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config BillingConfig
	err = yaml.Unmarshal(raw, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// niceFloatStr takes a float and gives back a monetary, human-formatted
// value.
func niceFloatStr(f float64) string {
	r := regexp.MustCompile("[0-9,]+.[0-9]{2}")
	p := message.NewPrinter(language.English)
	results := r.FindAllString(p.Sprintf("%f", f), 1)

	if len(results) < 1 {
		panic("got some ridiculous number that has no decimals")
	}

	return results[0]
}
