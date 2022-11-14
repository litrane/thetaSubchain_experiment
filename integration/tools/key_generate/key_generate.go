package main

import (
	"fmt"
	"io/ioutil"
	"log"

	ks "github.com/thetatoken/theta/wallet/softwallet/keystore"
)

type PriKeyFile struct {
	Address string                 `json:"address"`
	Crypto  map[string]interface{} `json:"crypto"`
	Id      string                 `json:"id"`
	Version int                    `json:"version"`
}

func getPrivateKeys() {

	// password := ""
	keystore, err := ks.NewKeystoreEncrypted("/Users/lipengze/subchainclientkeys", ks.StandardScryptN, ks.StandardScryptP)
	if err != nil {
		log.Fatalf("Failed to create key store: %v", err)
	}
	keyAddresses, err := keystore.ListKeyAddresses()
	if err != nil {
		log.Fatalf("Failed to get key address: %v", err)
	}
	outputStr := "["
	outputPrivateKey := "["
	for _, nodeAddrss := range keyAddresses {
		outputStr += "\"" + nodeAddrss.String() + "\"" + ","
		// nodeKey, err := keystore.GetKey(nodeAddrss, password)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		// nodePrivKey := nodeKey.PrivateKey
		// fmt.Println(nodeAddrss, ": ", nodePrivKey)
		// nodePrivKey.SaveToFile(fmt.Sprintf("/Users/lipengze/subchainclientkeys/%v", nodeAddrss))
		input, err := ioutil.ReadFile(fmt.Sprintf("/Users/lipengze/subchainclientkeys/%v", nodeAddrss))
		if err != nil {
			fmt.Println(err)
			return
		}
		outputPrivateKey += "\"" + string(input) + "\"" + ", "
	}
	outputStr += "]"
	outputPrivateKey += "]"
	fmt.Println(outputStr)
	fmt.Println(outputPrivateKey)

}

func main() {

	getPrivateKeys()

}
