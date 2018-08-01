package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	rpcAddr     = flag.String("rpc", "http://123.56.10.119:8545", "-rpc http://123.56.10.119:8545")
	accountAddr = flag.String("account", "", "-account 0x63fdb173af269faf42a85a6a5964bb72830b8151")
	sendAmounts = flag.String("amounts", "0xDE0B6B3A7640000", "-amounts 0xDE0B6B3A7640000")
	keyFile     = flag.String("keyfile", "", "-keyfile xxx.json")
	privateKey  = flag.String("privatekey", "", "-privatekey xxx")
	cmd         = flag.String("cmd", "", "-cmd balanceOf|sendETC")
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
	fmt.Printf("get ETC for one account\n")
	fmt.Printf("  %s -cmd balanceOf -rpc http://123.56.10.119:8545 -account 0x63fdb173af269faf42a85a6a5964bb72830b8151\n", prog)
	fmt.Printf("send ETC to one account\n")
	fmt.Printf("  %s -cmd sendETC -keyfile xxx.json -privatekey xxx -account 0xxxx -amounts 0xDE0B6B3A7640000\n", prog)
}

func getBalanceOf() error {
	c, err := rpc.Dial(*rpcAddr)
	if err != nil {
		fmt.Printf("rpc dial:%s failed,err:%s\n", *rpcAddr, err.Error())
		return err
	}
	defer c.Close()

	client := ethclient.NewClient(c)
	defer client.Close()

	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(*accountAddr), nil)
	if err != nil {
		fmt.Printf("get balance for %s failed,err:%s\n", *accountAddr, err.Error())
		return err
	}
	fmt.Printf("balance for %s : %s Wei\n", *accountAddr, (*balance).String())
	return nil
}

func sendETC() error {
	err := loadKeyfile(*keyFile)
	if err != nil {
		fmt.Printf("load keyfile:%s failed,err:%s\n", *keyFile, err.Error())
		return err
	}
	key, err := keystore.DecryptKey(jsonKey, *privateKey)
	if err != nil {
		fmt.Printf("DecryptKey failed,err:%s\n", err.Error())
		return err
	}

	c, err := rpc.Dial(*rpcAddr)
	if err != nil {
		fmt.Printf("rpc dial:%s failed,err:%s\n", *rpcAddr, err.Error())
		return err
	}
	defer c.Close()

	client := ethclient.NewClient(c)
	defer client.Close()

	nonce, err := client.NonceAt(context.Background(), key.Address, nil)
	if err != nil {
		fmt.Printf("get nonce by:%s failed,err:%s\n", key.Address.String(), err.Error())
		return err
	}
	fmt.Printf("nonce for:%s, is:%d\n", key.Address.String(), nonce)
	amount, _ := hexutil.DecodeBig(*sendAmounts)
	t := types.NewTransaction(nonce, common.HexToAddress(*accountAddr), amount, 1000000, big.NewInt(1), []byte{})
	trans, err := types.SignTx(t, types.HomesteadSigner{}, key.PrivateKey)
	if err != nil {
		fmt.Printf("trans WithSignature failed,err:%s\n", err.Error())
		return err
	}

	fmt.Printf("send transaction:%s\n", trans.Hash().String())
	return client.SendTransaction(context.Background(), trans)
}

func main() {
	flag.Parse()
	var err error
	switch *cmd {
	case "balanceOf":
		err = getBalanceOf()
	case "sendETC":
		err = sendETC()
	default:
		usage(os.Args[0])
	}
	if err != nil {
		fmt.Printf("cmd:%s failed,err:%s\n", *cmd, err.Error())
	}
}
