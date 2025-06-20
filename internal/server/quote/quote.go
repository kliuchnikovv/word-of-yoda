package quote

import (
	_ "embed"
	"math/rand"

	"fmt"

	"encoding/json"

	"github.com/kliuchnikovv/word-of-yoda/domain"
)

//go:embed data/yoda_quotes.json
var quotesData []byte
var quotes struct {
	Quotes []domain.Quote `json:"quotes"`
}

func init() {
	if err := json.Unmarshal(quotesData, &quotes); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal embedded quotes: %v", err))
	}
}

func GetRandomQuote() *domain.Quote {
	return &quotes.Quotes[rand.Intn(len(quotes.Quotes))]
}
