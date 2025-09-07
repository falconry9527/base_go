## 流程
```
abigen --abi=SimpleStorage.abi --pkg=main --type=SimpleStorage --out=simpleStorage.go
1. 生成 abi 文件： store_sol_Store.abi
solcjs --abi store.sol

2. 生成可调用内部方法的 go文件
abigen --abi=store_sol_Store.abi --pkg=store --out=store.go

3. 部署合约 remix
地址： 0x897D159F4b7AF148D3931C465dba822CB8DADc96

4. 调用合约

```