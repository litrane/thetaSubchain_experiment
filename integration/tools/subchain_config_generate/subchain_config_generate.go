package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ks "github.com/thetatoken/theta/wallet/softwallet/keystore"
)

type PriKeyFile struct {
	Address string                 `json:"address"`
	Crypto  map[string]interface{} `json:"crypto"`
	Id      string                 `json:"id"`
	Version int                    `json:"version"`
}

func createDirs(dirPath string, addresses []string, keyFiles []string, chainID int) {
	initValidatorSetJsonContent := fmt.Sprintf(initValidatorData, addresses[0], addresses[1], addresses[2], addresses[3])
	var genesisHash string
	var nodeDirPaths []string
	for i := 1; i <= 4; i++ {
		nodeDirPath := dirPath + fmt.Sprintf("node%v", i)
		nodeDirPaths = append(nodeDirPaths, nodeDirPath)
		os.MkdirAll(nodeDirPath+"/data", os.ModePerm)
		os.MkdirAll(nodeDirPath+"/key/encrypted", os.ModePerm)
		os.Chmod(nodeDirPath, os.ModePerm)
		// init_validator_set.json
		createFile(initValidatorSetJsonContent, nodeDirPath+"/data/init_validator_set.json")

		// key/encrypted
		input, err := ioutil.ReadFile(keyFiles[i-1])
		if err != nil {
			fmt.Println(err)
			return
		}
		keyName := strings.Split(keyFiles[i-1], "/")
		createFile(string(input), nodeDirPath+"/key/encrypted/"+keyName[len(keyName)-1])
		// key.plain
		keysDir := nodeDirPath + "/key"
		password := "qwe"
		keystore, err := ks.NewKeystoreEncrypted(keysDir, ks.StandardScryptN, ks.StandardScryptP)
		if err != nil {
			log.Fatalf("Failed to create key store: %v", err)
		}
		keyAddresses, err := keystore.ListKeyAddresses()
		if err != nil {
			log.Fatalf("Failed to get key address: %v", err)
		}
		nodeAddrss := keyAddresses[0]
		fmt.Println(nodeAddrss)
		nodeKey, err := keystore.GetKey(nodeAddrss, password)
		if err != nil {
			fmt.Println(err)
			return
		}

		nodePrivKey := nodeKey.PrivateKey
		nodePrivKey.SaveToFile(nodeDirPath + "/key.plain")
		if i == 1 {
			genesisHash = createGenesis(nodeDirPath, chainID)
		} else {
			// copy the genesis from node1
			input, err := ioutil.ReadFile(nodeDirPaths[0] + "/snapshot")
			if err != nil {
				fmt.Println(err)
				return
			}
			err = ioutil.WriteFile(nodeDirPath+"/snapshot", input, 0644)
			if err != nil {
				fmt.Println("Error creating snapshot")
				fmt.Println(err)
				return
			}
		}

		// config.yaml
		createFile(fmt.Sprintf(configYaml, genesisHash, 9*chainID, 9*chainID+1, 9*chainID+2, 9*chainID+3, i), nodeDirPath+"/config.yaml")
	}

}

func createFile(fileContent string, filePath string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(fileContent)
	write.Flush()
}

func createGenesis(filePath string, chainID int) string {
	/*
		cd $SUBCHAIN_HOME/integration/privatenet/node
		subchain_generate_genesis -mainchainID=privatenet -subchainID=tsub360777 -initValidatorSet=./data/init_validator_set.json -genesis=./snapshot
	*/
	chainIDStr := fmt.Sprintf("tsub360%03d", chainID)
	fmt.Println(chainIDStr)
	cmd := exec.Command("subchain_generate_genesis", "-mainchainID", "privatenet", "-subchainID", chainIDStr, "-initValidatorSet", "./data/init_validator_set.json", "-genesis", "./snapshot")
	cmd.Dir = filePath
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	// extractGenesisHash(string(out))

	intermediateGenesisHash := strings.Split(string(out), "Genesis block hash: ")
	genesisHash := intermediateGenesisHash[1][:66]
	return genesisHash
}

func main() {
	// createGenesis("/Users/lipengze/go/src/github.com/thetatoken/thetaSubchain_experiment/integration/allsubchains/DSN_360001/node1/")
	// return
	basePath := "/Users/lipengze/go/src/github.com/thetatoken/thetaSubchain_experiment/integration/allsubchains/DSN_360"
	for i := 1; i <= 3; i++ {
		fmt.Println(fmt.Sprintf("creating files for subchain 360%03d", i))
		var addresses []string
		var files []string

		dirPath := basePath + fmt.Sprintf("%03d", i) + "/"
		// fmt.Println(dirPath)

		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
		// fmt.Println(files[0])
		files = files[1:]
		for _, f := range files {
			fmt.Println(f)
		}

		for _, file := range files {
			// fmt.Println(file)
			jsonFile, err := os.Open(file)
			if err != nil {
				fmt.Println(err)
			}
			defer jsonFile.Close()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			// fmt.Println(string(byteValue))
			p := &PriKeyFile{}
			err = json.Unmarshal(byteValue, p)
			if err != nil {
				fmt.Println("解析数据失败", err)
				return
			}
			addresses = append(addresses, p.Address)
			// fmt.Println(p.Address)
		}
		createDirs(dirPath, addresses, files, i)
	}
}

const configYaml = `# Theta configuration
genesis:
  hash: "%v"
p2p:
  port: 12100
  seeds: 10.0.0.%v:12100, 10.0.0.%v:12100, 10.0.0.%v:12100, 10.0.0.%v:12100
  seedPeerOnlyOutbound: true
rpc:
  enable: true
log:
  levels: "*:info"
consensus:
  minBlockInterval: 1
subchain:
  mainchainEthRpcURL: "http://10.0.0.%v:18888/rpc"
  subchainEthRpcURL: "http://127.0.0.1:19888/rpc"
  chainRegistrarOnMainchain: "0x08425D9Df219f93d5763c3e85204cb5B4cE33aAa"
  mainchainTFuelTB: "0x7f1C87Bd3a22159b8a2E5D195B1a3283D10ea895"
  mainchainTNT20TB: "0x2Ce636d6240f8955d085a896e12429f8B3c7db26"
  mainchainTNT721TB: "0xEd8d61f42dC1E56aE992D333A4992C3796b22A74"
  mockTNT20: "0x47c5e40890bcE4a473A49D7501808b9633F29782"
  updateInterval: 200
`

const initValidatorData = `[
  {
    "address" : "%v",
    "stake"   : "100000000"
  },
  {
    "address" : "%v",
    "stake"   : "100000000"
  },
  {
    "address" : "%v",
    "stake"   : "100000000"
  },
  {
    "address" : "%v",
    "stake"   : "100000000"
  }
]
`
