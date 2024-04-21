package tax

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	*http.Response
	err error
}

func clientRequest(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	// Create HTTP client
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &Response{err: err}
	}
	return &Response{Response: res}
}

func uri(path ...string) string {
	baseURL := os.Getenv("TEST_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	if path == nil {
		return baseURL
	}
	return baseURL + "/" + strings.Join(path, "/")
}

func (r Response) Decode(v interface{}) error {
	log.Println(r.Body)
	if r.err != nil {
		return r.err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(v)
}

// func TestPostTaxHandler(t *testing.T) {
// 	taxResult := TaxResponse{
// 		Tax: 29000.0,
// 	}
// 	// Act
// 	res := clientRequest(http.MethodPost, uri("tax/calculations"), strings.NewReader(`{
// 		"totalIncome": 500000.0,
// 		"wht": 0.0,
// 		"allowances": [
// 		  {
// 			"allowanceType": "donation",
// 			"amount": 0.0
// 		  }
// 		]
// 	}`))

// 	var result IncomeDetails
// 	err := res.Decode(&result)

// 	assert.Nil(t, err)
// 	// assert.Equal(t, taxResult, calculateTax(result.TotalIncome))
// }
