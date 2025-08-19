package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 通过 Etherscan API 获取合约 ABI
func getContractABI(contractAddress string) (string, error) {
	etherscanAPI := "https://api-sepolia.etherscan.io/api"
	apiKey := "2BDYF5KGJGH5VAIIPXVFB7K6KWDVQ36M1G"

	url := fmt.Sprintf("%s?module=contract&action=getabi&address=%s&apikey=%s",
		etherscanAPI, contractAddress, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result["status"].(string) == "1" {
		return result["result"].(string), nil
	}

	return "", fmt.Errorf("failed to get ABI: %s", result["message"])
}

func main() {
	// USDT 合约地址
	contractAddress := "0x5560e1c2E0260c2274e400d80C30CDC4B92dC8ac"

	abi, err := getContractABI(contractAddress)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Contract ABI: %s\n", abi)
}
