package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	keyFile      = flag.String("keyfile", "", "-keyfile xxx.json")
	privateKey   = flag.String("privatekey", "", "-privatekey xxx")
	rpcAddr      = flag.String("rpc", "http://123.56.10.119:8545", "-rpc http://123.56.10.119:8545")
	contractAddr = flag.String("contractAddr", "0xcc33f3073f3e645f1a8ca094098cca68d8c1087c", "-contractAddr 0xcc33f3073f3e645f1a8ca094098cca68d8c1087c")
	accountAddr  = flag.String("account", "", "-account 0x63fdb173af269faf42a85a6a5964bb72830b8151")
	sendAmounts  = flag.Int64("amounts", 100, "-amounts 100")
	cmd          = flag.String("cmd", "", "-cmd tokenInfo|balanceOf|sendToken")
)

var (
	// 秘钥文件内容 json 串
	jsonKey []byte
)

func loadKeyfile(fileName string) error {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("read file failed,err:%s\n", err.Error())
		return err
	}
	jsonKey = d
	return nil
}

func usage(prog string) {
	fmt.Printf("get token information:\n")
	fmt.Printf("  %s -cmd tokenInfo -rpc http://123.56.10.119:8545 -contractAddr 0xcc33f3073f3e645f1a8ca094098cca68d8c1087c\n", prog)
	fmt.Printf("get balance of one account\n")
	fmt.Printf("  %s -cmd balanceOf -rpc http://123.56.10.119:8545 -contractAddr 0xcc33f3073f3e645f1a8ca094098cca68d8c1087c\n", prog)
	fmt.Printf("    -account 0x49e7888acb220790b363e7061a8a9b46d58bfdc8\n")
	fmt.Printf("send token to one account\n")
	fmt.Printf("  %s -cmd sendToken -rpc http://123.56.10.119:8545 -contractAddr 0xcc33f3073f3e645f1a8ca094098cca68d8c1087c\n", prog)
	fmt.Printf("    -keyfile xxx.json -privatekey xxx -account 0xxx -amounts 0x2710\n")
}

func getTokenInfo() error {
	// RPC 拨号
	dial, err := rpc.Dial(*rpcAddr)
	if err != nil {
		fmt.Printf("rpc dail faield,err:%s\n", err.Error())
		return err
	}
	defer dial.Close()

	cli := ethclient.NewClient(dial)
	defer cli.Close()
	// 操作合约对象
	tk, err := NewToken(common.HexToAddress(*contractAddr), cli)
	if err != nil {
		fmt.Printf("newToken failed,err:%s\n", err.Error())
		return err
	}
	// 获取代币名称
	name, err := tk.Name(nil)
	if err != nil {
		fmt.Printf("get token Name failed,err:%s\n", err.Error())
		return err
	}
	fmt.Printf("TokenName:%s\n", name)
	// 获取代币精确位数
	dec, err := tk.Decimals(nil)
	if err != nil {
		fmt.Printf("get token Decimals failed,err:%s\n", err.Error())
		return err
	}
	fmt.Printf("Decimals:%d\n", dec)
	// 获取代币符号
	sym, err := tk.Symbol(nil)
	if err != nil {
		fmt.Printf("get token symbol failed,err:%s\n", err.Error())
		return err
	}
	fmt.Printf("TokenSymbol:%s\n", sym)
	// 获取代币发行量
	total, err := tk.TotalSupply(nil)
	if err != nil {
		fmt.Printf("get token totalsupply failed,err:%s\n", err.Error())
		return err
	}
	fmt.Printf("TokenTotalSupply:%s\n", total.String())
	return nil
}

func getBalanceOf() error {
	// RPC 拨号
	dial, err := rpc.Dial(*rpcAddr)
	if err != nil {
		fmt.Printf("rpc dail faield,err:%s\n", err.Error())
		return err
	}
	defer dial.Close()

	cli := ethclient.NewClient(dial)
	defer cli.Close()
	// 操作合约对象
	tk, err := NewToken(common.HexToAddress(*contractAddr), cli)
	if err != nil {
		fmt.Printf("newToken failed,err:%s\n", err.Error())
		return err
	}
	// 获取代币余额
	var out = new(*big.Int)
	err = tk.TokenCaller.contract.Call(nil, &out, "balanceOf", common.HexToAddress(*accountAddr))
	if err != nil {
		fmt.Printf("balanceOf failed,err:%s\n", err.Error())
		return err
	}
	fmt.Printf("%s balance:%s\n", *accountAddr, (*out).String())
	return nil
}

// 代币转移
func sendToken() error {
	if err := loadKeyfile(*keyFile); err != nil {
		fmt.Printf("loadKeyfile failed,err:%s\n", err.Error())
		return err
	}

	// RPC 拨号
	dial, err := rpc.Dial(*rpcAddr)
	if err != nil {
		fmt.Printf("rpc dail faield,err:%s\n", err.Error())
		return err
	}
	defer dial.Close()
	cli := ethclient.NewClient(dial)
	defer cli.Close()
	/*
		NewToken接口在token.go文件中，其生成命令如下
		./abigen --abi token.abi --pkg main --type Token --out token.go
		其中token.abi文件内容为合约的ABI json串
	*/
	// 操作合约对象
	tk, err := NewToken(common.HexToAddress(*contractAddr), cli)
	if err != nil {
		fmt.Printf("newToken failed,err:%s\n", err.Error())
		return err
	}

	// 根据秘钥文件和密码准备签名
	trOpts, err := bind.NewTransactor(bytes.NewReader(jsonKey), *privateKey)
	if err != nil {
		fmt.Printf("new transactor failed,err:%s\n", err.Error())
		return err
	}

	// 发送代币
	tr, err := tk.TokenTransactor.Transfer(trOpts, common.HexToAddress(*accountAddr), big.NewInt(*sendAmounts))
	if err != nil {
		fmt.Printf("transfer failed,err:%s\n", err.Error())
		return err
	}
	fmt.Printf("transaction hash:%s\n", tr.Hash().String())
	return nil
}

func main() {
	flag.Parse()
	var err error

	switch *cmd {
	case "tokenInfo":
		err = getTokenInfo()
	case "balanceOf":
		err = getBalanceOf()
	case "sendToken":
		err = sendToken()
	default:
		usage(os.Args[0])
	}
	if err != nil {
		fmt.Printf("cmd:%s,err:%s\n", *cmd, err.Error())
	}
}
