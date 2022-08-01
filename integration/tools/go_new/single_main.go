package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/crypto"
)

const CHAIN_ID_OFFSET int64 = 360

func main() {
	// mainchainTfuelLock(big.NewInt(33))
	// mainchainTNT721Lock(big.NewInt(33))
	// subchainTNT721Burn(big.NewInt(33))
	// client, err := ethclient.Dial("http://localhost:19888/rpc")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// TfuelTokenBankInstance, err := ct.NewTFuelTokenBank(subchainTfuelTokenBank, client)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// //auth := subchainSelectAccount(client, 1)
	// tx, _ := TfuelTokenBankInstance.Id(nil)
	// fmt.Println(tx)
	// chainID, err := client.ChainID(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(chainID)
	chainIDStr := "tsub_360777"
	chainIDWithoutOffset := new(big.Int).Abs(crypto.Keccak256Hash(common.Bytes(chainIDStr)).Big())
	chainID1 := big.NewInt(1).Add(big.NewInt(CHAIN_ID_OFFSET), chainIDWithoutOffset)
	//chainID1.SetString("tsub_360777", 16)
	fmt.Println(hex.EncodeToString(chainID1.Bytes()))
}
