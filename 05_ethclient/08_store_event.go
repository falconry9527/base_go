package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"int256","name":"key","type":"int256"},{"indexed":false,"internalType":"int256","name":"value","type":"int256"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"int256","name":"","type":"int256"}],"name":"items","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"int256","name":"key","type":"int256"},{"internalType":"int256","name":"value","type":"int256"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/bdb2ede84fe04e41a6fc9b2c9506d8c7")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x897D159F4b7AF148D3931C465dba822CB8DADc96")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6920583),
		// ToBlock:   big.NewInt(2394201),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	// 获取查询合约的方法
	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatal(err)
	}

	// 获取不同合约的签名hash
	eventSignature := "ItemSet(int256,int256)"
	hashItemSet := crypto.Keccak256Hash([]byte(eventSignature))
	fmt.Printf("事件 ItemSet 的签名哈希: %s\n", hashItemSet.Hex())

	for _, vLog := range logs {
		switch vLog.Topics[0].String() {
		case hashItemSet.Hex():
			handleItemSetEvent(contractAbi, vLog)
			continue
		default:
		}
	}
}

func handleItemSetEvent(contractAbi abi.ABI, vLog types.Log) {
	fmt.Printf("发现 ItemSet 事件 (TxHash: %s, Block: %d)\n", vLog.TxHash.Hex(), vLog.BlockNumber)
	//  解析日志数据
	var event ItemSetEvent
	// 使用 ABI 解包日志数据到 event 结构体
	// 因为 key 和 value 在事件定义中都没有 indexed，所以它们都在 Data 字段中
	err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
	if err != nil {
		log.Printf("解析日志数据失败: %v", err)
	}
	// 输出解析到的事件参数
	fmt.Printf("Key: %s\n", event.Key.String())
	fmt.Printf("Value: %s\n", event.Value.String())
	fmt.Println("------------------------------------")
}

type ItemSetEvent struct {
	Key   *big.Int
	Value *big.Int
}
