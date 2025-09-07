package main

import (
	store "base_go/05_ethclient/callFun"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 1. 连接到以太坊节点（不需要私钥！）
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/bdb2ede84fe04e41a6fc9b2c9506d8c7")
	if err != nil {
		log.Fatal("连接节点失败:", err)
	}
	defer client.Close()

	// 3. 创建合约实例
	instance, err := store.NewStore(common.HexToAddress("0x897D159F4b7AF148D3931C465dba822CB8DADc96"), client)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 定义要查询的键（int256类型）
	keysToQuery := []*big.Int{
		big.NewInt(1), // 查询 key = 1
		big.NewInt(2), // 查询 key = 42
		big.NewInt(3), // 查询 key = 100
	}

	// 5. 批量查询映射中的值
	fmt.Println("开始查询 mapping(int256 => int256) public items:")
	fmt.Println("==============================================")

	for _, key := range keysToQuery {
		value, err := instance.Items(nil, key)
		if err != nil {
			log.Printf("查询键 %s 失败: %v", key.String(), err)
			continue
		}
		fmt.Printf("items[%s] = %s\n", key.String(), value.String())
	}

}
