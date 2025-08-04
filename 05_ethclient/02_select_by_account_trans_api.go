package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	etherscanAPI = "https://api.etherscan.io/api"
	apiKey       = "2BDYF5KGJGH5VAIIPXVFB7K6KWDVQ36M1G"
	account      = "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
)

type EtherscanResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

type Transaction struct {
	BlockNumber string `json:"blockNumber"`
	TimeStamp   string `json:"timeStamp"`
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
}

func main() {
	//1.获取所有交易
	url := fmt.Sprintf("%s?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=desc&apikey=%s",
		etherscanAPI, account, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result EtherscanResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}

	if result.Status != "1" {
		log.Fatalf("Etherscan API 错误: %s", result.Message)
	}

	fmt.Printf("找到 %d 笔交易\n", len(result.Result))
	for _, tx := range result.Result {
		fmt.Printf("区块: %s, 时间: %s, 哈希: %s\n", tx.BlockNumber, tx.TimeStamp, tx.Hash)
		fmt.Printf("从: %s 到: %s 金额: %s wei\n", tx.From, tx.To, tx.Value)
		fmt.Println("----------------------------------")
	}
}
