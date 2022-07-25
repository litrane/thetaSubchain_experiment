package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/thetasubchain/eth/ethclient"
	ct "github.com/thetatoken/thetasubchain/integration/tools/go_new/accessors"
)

func mainchainTNT20Locked() {
	AccountsInit()
	client, err := ethclient.Dial("http://localhost:18888/rpc")
	if err != nil {
		log.Fatal(err)
	}
	subchainID := big.NewInt(360777)
	var dec18 = new(big.Int)
	dec18.SetString("1000000000000000000", 10)
	user := accountList[1].fromAddress

	instanceTNT20VoucherContract, err := ct.NewTNT20VoucherContract(TNT20VoucherContractAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	instanceTNT20TokenBank, err := ct.NewTNT20TokenBank(TNT20TokenBankAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	authAccount0 := SelectAccount(client, 0)
	_, err = instanceTNT20VoucherContract.Mint(authAccount0, user, big.NewInt(30))
	if err != nil {
		log.Fatal(err)
	}
	authUser := SelectAccount(client, 1)
	_, err = instanceTNT20VoucherContract.Approve(authUser, TNT20TokenBankAddress, big.NewInt(30))
	if err != nil {
		log.Fatal(err)
	}

	authUser = SelectAccount(client, 1)
	LockTx, err := instanceTNT20TokenBank.LockTokens(authUser, subchainID, TNT20VoucherContractAddress, user, big.NewInt(25))
	if err != nil {
		log.Fatal(err)
	}

	subchainClient, _ := ethclient.Dial("http://localhost:19888/rpc")
	fromHeight, _ := subchainClient.BlockNumber(context.Background())
	receipt, err := client.TransactionReceipt(context.Background(), LockTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status != 1 {
		fmt.Println("lock error")
	}
	fmt.Println(LockTx.Hash())
	var subchainVoucherAddress common.Address
	for {
		time.Sleep(2 * time.Second)
		toHeight, _ := subchainClient.BlockNumber(context.Background())
		result := getMintlog(int(fromHeight), int(toHeight), SubchainTNT20TokenBankAddress, user)
		if result != nil {
			subchainVoucherAddress = *result
			break
		}
	}
	fmt.Println(subchainVoucherAddress)
}
func subchainTNT20Lock() {
	client, err := ethclient.Dial("http://localhost:19888/rpc")
	if err != nil {
		log.Fatal(err)
	}
	subchainTNT20VoucherAddress := common.HexToAddress("0x7D7e270b7E279C94b265A535CdbC00Eb62E6e68f")
	if err != nil {
		log.Fatal(err)
	}
	testSubchainTNT20TokenBankInstance, err := ct.NewTestSubchainTNT20TokenBank(SubchainTNT20TokenBankAddress, client)
	subchainTNT20VoucherAddressInstance, _ := ct.NewTNT20VoucherContract(subchainTNT20VoucherAddress, client)
	tx, err := subchainTNT20VoucherAddressInstance.BalanceOf(nil, accountList[1].fromAddress)
	fmt.Println("subchain_account_balance is", tx)
	//await TNT20Token.approve(TNT20TokenBank1.address, 20, { from: valGuarantor1 })
	authUser := SubchainSelectAccount(client, 1)
	tx2, err := subchainTNT20VoucherAddressInstance.Approve(authUser, SubchainTNT20TokenBankAddress, big.NewInt(20))

	authUser = SubchainSelectAccount(client, 1)
	tx2, err = testSubchainTNT20TokenBankInstance.LockTokens(authUser, big.NewInt(366), subchainTNT20VoucherAddress, accountList[6].fromAddress, big.NewInt(20))
	//tx, err := testSubchainTNT20TokenBankInstance.GetDenom(nil, subchainTNT20VoucherAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx2.Hash().Hex())

	receipt, err := client.TransactionReceipt(context.Background(), tx2.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status != 1 {
		fmt.Println("lock error")
	}
	mainchainClient, _ := ethclient.Dial("http://localhost:18888/rpc")
	fromHeight, _ := mainchainClient.BlockNumber(context.Background())
	var mainchainVoucherAddress common.Address
	for {
		tx, err = subchainTNT20VoucherAddressInstance.BalanceOf(nil, accountList[1].fromAddress)
		fmt.Println("subchain_account_balance is", tx)
		time.Sleep(2 * time.Second)
		toHeight, _ := mainchainClient.BlockNumber(context.Background())
		result := getMainchainTNT20Mintlog(int(fromHeight), int(toHeight), TNT20TokenBankAddress, accountList[6].fromAddress)
		if result != nil {
			mainchainVoucherAddress = *result
			break
		}
	}
	fmt.Println(mainchainVoucherAddress)
}
func mainchainTNT20Burn() {
	client, err := ethclient.Dial("http://localhost:18888/rpc")
	if err != nil {
		log.Fatal(err)
	}

	mainchainTNT20VoucherAddress := common.HexToAddress("0xb0DBBcba1Be5B71Dcb42aB1935773B3675e645e8")
	mainchainTNT20TokenBankInstance, err := ct.NewTNT20TokenBank(TNT20TokenBankAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	mainchainTNT20VoucherAddressInstance, _ := ct.NewTNT20VoucherContract(mainchainTNT20VoucherAddress, client)
	tx, err := mainchainTNT20VoucherAddressInstance.BalanceOf(nil, accountList[6].fromAddress)
	fmt.Println("mainchain account 6 tnt20 balance is", tx)
	//await TNT20Token.approve(TNT20TokenBank1.address, 20, { from: valGuarantor1 })
	authUser := SelectAccount(client, 6)
	tx1, err := mainchainTNT20VoucherAddressInstance.Approve(authUser, TNT20TokenBankAddress, big.NewInt(11))

	authUser = SelectAccount(client, 6)
	tx1, err = mainchainTNT20TokenBankInstance.BurnVouchers(authUser, mainchainTNT20VoucherAddress, accountList[1].fromAddress, big.NewInt(10))
	//tx, err := testSubchainTNT20TokenBankInstance.GetDenom(nil, subchainTNT20VoucherAddress)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(tx)
	fmt.Println("burn tx is", tx1.Hash().Hex())
	receipt, err := client.TransactionReceipt(context.Background(), tx1.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status != 1 {
		fmt.Println("lock error")
	}
	time.Sleep(2 * time.Second)
	tx, err = mainchainTNT20VoucherAddressInstance.BalanceOf(nil, accountList[6].fromAddress)
	fmt.Println("mainchain account 6 tnt20 balance is", tx)
	subchainClient, _ := ethclient.Dial("http://localhost:19888/rpc")
	subchainTNT20VoucherAddress := common.HexToAddress("0x7D7e270b7E279C94b265A535CdbC00Eb62E6e68f")
	subchainTNT20VoucherAddressInstance, _ := ct.NewTNT20VoucherContract(subchainTNT20VoucherAddress, subchainClient)
	tx, err = subchainTNT20VoucherAddressInstance.BalanceOf(nil, accountList[1].fromAddress)
	fmt.Println("subchain_account_balance is", tx)
	for {
		time.Sleep(2 * time.Second)
		tx_subchain, _ := subchainTNT20VoucherAddressInstance.BalanceOf(nil, accountList[1].fromAddress)
		if tx.Cmp(tx_subchain) != 0 {
			tx = tx_subchain
			break
		}
	}
	fmt.Println("subchain_account_balance is", tx)
}
func subchainTNT20Burn() {
	client, err := ethclient.Dial("http://localhost:19888/rpc")
	if err != nil {
		log.Fatal(err)
	}

	subchainTNT20VoucherAddress := common.HexToAddress("0x7D7e270b7E279C94b265A535CdbC00Eb62E6e68f")
	if err != nil {
		log.Fatal(err)
	}
	testSubchainTNT20TokenBankInstance, err := ct.NewTestSubchainTNT20TokenBank(SubchainTNT20TokenBankAddress, client)
	subchainTNT20VoucherAddressInstance, _ := ct.NewTNT20VoucherContract(subchainTNT20VoucherAddress, client)
	tx, err := subchainTNT20VoucherAddressInstance.BalanceOf(nil, accountList[1].fromAddress)
	fmt.Println("subchain_account_balance is", tx)
	//await TNT20Token.approve(TNT20TokenBank1.address, 20, { from: valGuarantor1 })
	authUser := SubchainSelectAccount(client, 1)
	tx1, err := subchainTNT20VoucherAddressInstance.Approve(authUser, SubchainTNT20TokenBankAddress, big.NewInt(11))

	authUser = SubchainSelectAccount(client, 1)
	tx1, err = testSubchainTNT20TokenBankInstance.BurnVouchers(authUser, subchainTNT20VoucherAddress, accountList[1].fromAddress, big.NewInt(10))
	//tx, err := testSubchainTNT20TokenBankInstance.GetDenom(nil, subchainTNT20VoucherAddress)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(tx)
	fmt.Println("burn tx is", tx1.Hash().Hex())
	receipt, err := client.TransactionReceipt(context.Background(), tx1.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status != 1 {
		fmt.Println("lock error")
		return
	}
	time.Sleep(2 * time.Second)
	tx, err = subchainTNT20VoucherAddressInstance.BalanceOf(nil, accountList[1].fromAddress)
	fmt.Println("subchain_account_balance is", tx)
	mainchainClient, _ := ethclient.Dial("http://localhost:18888/rpc")
	mainchainTNT20VoucherAddressInstance, _ := ct.NewTNT20VoucherContract(TNT20VoucherContractAddress, mainchainClient)
	tx, err = mainchainTNT20VoucherAddressInstance.BalanceOf(nil, accountList[1].fromAddress)
	fmt.Println("mainchain_account_balance is", tx)
	for {
		time.Sleep(2 * time.Second)
		tx_subchain, _ := mainchainTNT20VoucherAddressInstance.BalanceOf(nil, accountList[1].fromAddress)
		fmt.Println("mainchain_account_balance is", tx_subchain)
		if tx.Cmp(tx_subchain) != 0 {
			tx = tx_subchain
			break
		}
	}
	fmt.Println("mainchain_account_balance is", tx)
}