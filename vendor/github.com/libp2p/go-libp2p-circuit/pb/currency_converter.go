package main

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	BaseCurrency    string
	TargetCurrency  string
	Rate            float64
	LastUpdated     time.Time
}

type CurrencyConverter struct {
	rates map[string]ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]ExchangeRate),
	}
}

func (c *CurrencyConverter) AddRate(base, target string, rate float64) {
	key := base + "_" + target
	c.rates[key] = ExchangeRate{
		BaseCurrency:   base,
		TargetCurrency: target,
		Rate:           rate,
		LastUpdated:    time.Now(),
	}
}

func (c *CurrencyConverter) Convert(amount float64, base, target string) (float64, error) {
	if base == target {
		return amount, nil
	}

	key := base + "_" + target
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", base, target)
	}

	return amount * rate.Rate, nil
}

func (c *CurrencyConverter) GetSupportedPairs() []string {
	var pairs []string
	for key := range c.rates {
		pairs = append(pairs, key)
	}
	return pairs
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.92)
	converter.AddRate("EUR", "USD", 1.09)
	converter.AddRate("USD", "JPY", 148.50)
	
	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	fmt.Printf("Supported pairs: %v\n", converter.GetSupportedPairs())
}package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	cache      map[string]ExchangeRateResponse
	lastUpdate time.Time
}

func NewCurrencyConverter(apiKey string) *CurrencyConverter {
	return &CurrencyConverter{
		apiKey:  apiKey,
		baseURL: "https://api.exchangerate-api.com/v4/latest/",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: make(map[string]ExchangeRateResponse),
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	rates, err := c.getExchangeRates(from)
	if err != nil {
		return 0, err
	}

	rate, exists := rates.Rates[to]
	if !exists {
		return 0, fmt.Errorf("currency %s not supported", to)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) getExchangeRates(baseCurrency string) (*ExchangeRateResponse, error) {
	if cached, exists := c.cache[baseCurrency]; exists {
		if time.Since(c.lastUpdate) < 30*time.Minute {
			return &cached, nil
		}
	}

	url := c.baseURL + baseCurrency
	if c.apiKey != "" {
		url += "?api_key=" + c.apiKey
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rates ExchangeRateResponse
	if err := json.Unmarshal(body, &rates); err != nil {
		return nil, err
	}

	c.cache[baseCurrency] = rates
	c.lastUpdate = time.Now()

	return &rates, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies(baseCurrency string) ([]string, error) {
	rates, err := c.getExchangeRates(baseCurrency)
	if err != nil {
		return nil, err
	}

	currencies := make([]string, 0, len(rates.Rates))
	for currency := range rates.Rates {
		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func main() {
	converter := NewCurrencyConverter("")

	amount := 100.0
	from := "USD"
	to := "EUR"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)

	currencies, err := converter.GetSupportedCurrencies("USD")
	if err != nil {
		fmt.Printf("Error fetching currencies: %v\n", err)
		return
	}

	fmt.Printf("Supported currencies: %v\n", currencies[:10])
}