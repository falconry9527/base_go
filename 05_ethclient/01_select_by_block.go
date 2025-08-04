package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	// 测试链
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/bdb2ede84fe04e41a6fc9b2c9506d8c7")
	if err != nil {
		log.Fatal(err)
	}
	// 查询最新的区块号
	header1, err1 := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err1)
	}
	fmt.Println("Latest block number:", header1.Number.String())

	// 获取指定区块头的信息
	blockNumber := big.NewInt(5671744)
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("-----header -------------")
	fmt.Println(header.Number.Uint64())     // 5671744
	fmt.Println(header.Time)                // 1527211625
	fmt.Println(header.Difficulty.Uint64()) // 3217000136609065
	fmt.Println(header.Hash().Hex())        // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9

	// 获取指定区块的信息
	// block, err := client.BlockByNumber(context.Background(), blockNumber)
	block, err := client.BlockByHash(context.Background(), header.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("-----block-------------")
	fmt.Println(block.Number().Uint64())     // 5671744
	fmt.Println(block.Time())                // 1712798400
	fmt.Println(block.Difficulty().Uint64()) // 0
	fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	fmt.Println(len(block.Transactions()))   // 70

	// 遍历每一笔交易： 交易树
	//for idx, tx := range block.Transactions() {
	idx := 143
	tx := block.Transactions()[143]

	fmt.Println("\n-----Transactions -------------")
	fmt.Printf("=== 交易 %d/%d ===", idx+1, len(block.Transactions()))
	fmt.Println("交易哈希:", tx.Hash().Hex())
	// 交易数据
	fmt.Println("区块ID:", block.Number())
	fmt.Println("交易哈希:", tx.Hash().Hex())
	fmt.Println("发送方:", getSender(tx))
	fmt.Println("接收方:", getReceiver(tx))
	fmt.Println("金额:", weiToEther(tx.Value()), "ETH")
	fmt.Println("Gas Limit:", tx.Gas())
	fmt.Println("Gas Price:", weiToGwei(tx.GasPrice()), "Gwei")
	fmt.Println("Nonce:", tx.Nonce())

	// 收据数据
	fmt.Println("=== 收据数据 ===")
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		fmt.Println("警告: 无法获取交易收据 -", err)
		// return
	}
	printTransactionReceiptDetails(receipt, block.Number())
	// 获取余额
	fmt.Println("=== 余额信息 ===")
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(getSender(tx)), block.Number())
	if err != nil {
		fmt.Println("警告: 获取余额 -", err)
		// return
	}
	fmt.Println("发送方:", getSender(tx), "余额", balance)

	// 遍历交易收据： 收据树
	//receiptSlice, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithHash(block.Hash(), false))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, receipt := range receiptSlice {
	//	printTransactionReceiptDetails(receipt, block.Number())
	//}

}

func printTransactionReceiptDetails(receipt *types.Receipt, blockNumber *big.Int) {
	// 收据信息
	fmt.Println("区块ID", receipt.BlockNumber.Uint64())
	fmt.Println("交易Hash:", receipt.TxHash.Hex())
	fmt.Println("状态:", receiptStatus(receipt.Status))
	fmt.Println("实际Gas消耗:", receipt.GasUsed)
	fmt.Println("区块确认数:", blockNumber.Uint64()-receipt.BlockNumber.Uint64())
	fmt.Println("合约地址:", receipt.ContractAddress.Hex())
	fmt.Println("日志数量:", len(receipt.Logs))
}

// 辅助函数
func getSender(tx *types.Transaction) string {
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return "未知"
	}
	return sender.Hex()
}

func getReceiver(tx *types.Transaction) string {
	if tx.To() == nil {
		return "合约创建"
	}
	return tx.To().Hex()
}

func receiptStatus(status uint64) string {
	if status == 1 {
		return "成功"
	}
	return "失败"
}

func weiToEther(wei *big.Int) string {
	f := new(big.Float).SetInt(wei)
	f = new(big.Float).Quo(f, big.NewFloat(1e18))
	return f.String()
}

func weiToGwei(wei *big.Int) string {
	f := new(big.Float).SetInt(wei)
	f = new(big.Float).Quo(f, big.NewFloat(1e9))
	return f.String()
}
