package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		TestConfirmTransaction()
		wg.Done()
	}()
	go func() {
		TestConfirmTransaction()
		wg.Done()
	}()
	go func() {
		TestConfirmTransaction()
		wg.Done()
	}()
	go func() {
		TestConfirmTransaction()
		wg.Done()
	}()
	wg.Wait()
}
func TestConfirmTransaction() {

	url := "http://localhost:8080/api/v1/wallet/confirm"
	method := "POST"

	payload := strings.NewReader(`{
  "transaction_id": 30
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
