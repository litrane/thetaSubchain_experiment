package tools

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/thetasubchain/eth/ethclient"
	ct "github.com/thetatoken/thetasubchain/interchain/contracts/accessors"
)

func SubchainChannelRegister(targetChainID *big.Int, IP string) {
	subchainClient, err := ethclient.Dial("http://localhost:19888/rpc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Preparing for RegisterChannelOnSubchain for chainID %v and IP is %v...\n", targetChainID.String(), IP)
	subchainRegisterAddr := common.HexToAddress("0xBd770416a3345F91E4B34576cb804a576fa48EB1")
	subchainRegisterInstance, _ := ct.NewChainRegistrarOnSubchain(subchainRegisterAddr, subchainClient)

	authUser := subchainSelectAccount(subchainClient, 1)
	regitserTx, err := subchainRegisterInstance.RegisterSubchainChannel(authUser, targetChainID, IP)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("registerTX", regitserTx.Hash().Hex())
	receipt, err := subchainClient.TransactionReceipt(context.Background(), regitserTx.Hash())
	time.Sleep(6 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status != 1 {
		log.Fatal("register error")
	}

}

func GetMaxProcessedNonceFromRegistrar() {
	subchainClient, err := ethclient.Dial("http://localhost:19888/rpc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Preparing for QueryMaxNonceForSubchainChannel...\n")
	subchainRegisterAddr := common.HexToAddress("0xBd770416a3345F91E4B34576cb804a576fa48EB1")
	subchainRegisterInstance, _ := ct.NewChainRegistrarOnSubchain(subchainRegisterAddr, subchainClient)

	maxProcessedSubchainRegisteredNonce, err := subchainRegisterInstance.GetMaxProcessedNonce(nil)
	if err != nil {
		log.Fatal(err)
		return // ignore
	}
	log.Println("Max nonce : ", maxProcessedSubchainRegisteredNonce)

}

func GetCrossChainFeeFromRegistrar() {
	subchainClient, err := ethclient.Dial("http://localhost:19888/rpc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Preparing for Chain Fee...\n")
	subchainRegisterAddr := common.HexToAddress("0xBd770416a3345F91E4B34576cb804a576fa48EB1")
	subchainRegisterInstance, _ := ct.NewChainRegistrarOnSubchain(subchainRegisterAddr, subchainClient)

	maxProcessedSubchainRegisteredNonce, err := subchainRegisterInstance.GetCrossChainFee(nil)
	if err != nil {
		log.Fatal(err)
		return // ignore
	}
	log.Println("Chain Fee : ", maxProcessedSubchainRegisteredNonce)

}

func GetChannelStatusFromRegistrar(targetChainID *big.Int, targetChainEthRpcClientURL string) {
	subchainClient, err := ethclient.Dial(targetChainEthRpcClientURL)
	if err != nil {
		log.Fatal(err)
	}
	localChainID, err := subchainClient.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Preparing for Query Channel Status...\n")
	subchainRegisterAddr := common.HexToAddress("0xBd770416a3345F91E4B34576cb804a576fa48EB1")
	subchainRegisterInstance, _ := ct.NewChainRegistrarOnSubchain(subchainRegisterAddr, subchainClient)

	channelStatus, err := subchainRegisterInstance.IsAnActiveChannel(nil, targetChainID)
	if err != nil {
		log.Fatal(err)
		return // ignore
	}
	log.Printf("The Channel from %v to %v is active? : %v", localChainID.String(), targetChainID.String(), channelStatus)
}

func VerifyChannel(targetChainID *big.Int, targetChainEthRpcClientURL string) {
	subchainClient, err := ethclient.Dial(targetChainEthRpcClientURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Preparing for Verify Query Channel for chainID %v...\n", targetChainID.String())
	subchainRegisterAddr := common.HexToAddress("0xBd770416a3345F91E4B34576cb804a576fa48EB1")
	subchainRegisterInstance, _ := ct.NewChainRegistrarOnSubchain(subchainRegisterAddr, subchainClient)
	authUser := subchainSelectAccount(subchainClient, 1)
	tx, err := subchainRegisterInstance.UpdateSubchainChannelStatus(authUser, targetChainID, true, big.NewInt(2))
	if err != nil {
		log.Fatal(err)
		return // ignore
	}
	fmt.Println(tx.Hash().Hex())
}
