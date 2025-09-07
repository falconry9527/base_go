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

var StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"address1","type":"address"},{"indexed":true,"internalType":"int256","name":"key1","type":"int256"},{"indexed":true,"internalType":"string","name":"keyStr","type":"string"},{"indexed":false,"internalType":"int256","name":"key","type":"int256"},{"indexed":false,"internalType":"int256","name":"value","type":"int256"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"int256","name":"","type":"int256"}],"name":"items","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"int256","name":"key","type":"int256"},{"internalType":"int256","name":"value","type":"int256"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/bdb2ede84fe04e41a6fc9b2c9506d8c7")
	if err != nil {
		log.Fatal(err)
	}
	contractAddress := common.HexToAddress("0xE74f46C9D1E06f764f7e2057CA9AB11ad4768981")
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
	eventSignature := "ItemSet(address,int256,string,int256,int256)"
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

func handleItemStrSetEvent(contractAbi abi.ABI, vLog types.Log) {
	fmt.Printf("发现 ItemSet 事件 (TxHash: %s, Block: %d)\n", vLog.TxHash.Hex(), vLog.BlockNumber)
	// 从topic中获取数据
	// 	side := uint8(new(big.Int).SetBytes(log.Topics[1].Bytes()).Uint64())
	// 输出解析到的事件参数

	topicKeyHash := vLog.Topics[1]   // 第一个 indexed 参数 (key) 的哈希
	topicValueHash := vLog.Topics[2] // 第二个 indexed 参数 (value) 的哈希

	fmt.Printf("Key 的哈希 (Topic[1]): %s\n", topicKeyHash.Hex()) //0x2b3d5ecfad02814a9cc25dbb78fd85dc61a80178d1d2800b4f8443e41dc24f0f
	fmt.Printf("Value 的哈希 (Topic[2]): %s\n", topicValueHash.Hex())

	fmt.Println("------------------------------------")
}

func handleItemSetEvent(contractAbi abi.ABI, vLog types.Log) {
	fmt.Printf("发现 ItemSet 事件 (TxHash: %s, Block: %d)\n", vLog.TxHash.Hex(), vLog.BlockNumber)
	//  解析日志数据
	// 1. 解析topic数据

	// 地址类型数据解析
	fromAddress := common.BytesToAddress(vLog.Topics[1].Bytes()) // 直接转换
	valueBytes := vLog.Topics[2].Bytes()
	// 数字类型解析
	value := new(big.Int).SetBytes(valueBytes) // 先将字节转换为 big.Int
	// string 类型解析： string 类型不能解析成原始数据，只能解析成对应的hash值
	topic3Hash := vLog.Topics[3] // 第二个 indexed 参数 (value) 的哈希
	fmt.Println("--topic--")

	fmt.Printf("address1 (Topic[1]): %s\n", fromAddress.Hex()) //0x2b3d5ecfad02814a9cc25dbb78fd85dc61a80178d1d2800b4f8443e41dc24f0f
	fmt.Printf("key1 (Topic[2]): %s\n", value)
	fmt.Printf("keyStr (Topic[3]): %s\n", topic3Hash.Hex())

	// 2. 解析data数据
	fmt.Println("--data--")
	var event ItemSetEvent
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
