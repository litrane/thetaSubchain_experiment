package tools

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/thetatoken/theta/crypto"
	"github.com/thetatoken/theta/crypto/sha3"

	"github.com/thetatoken/theta/common"
	scom "github.com/thetatoken/thetasubchain/common"
	"github.com/thetatoken/thetasubchain/eth/abi/bind"
	"github.com/thetatoken/thetasubchain/eth/ethclient"

	// rg "chainRegistrarOnMainchain" // for demo
	ct "github.com/thetatoken/thetasubchain/interchain/contracts/accessors"
)

type accounts struct {
	priKey      string
	privateKey  *ecdsa.PrivateKey
	fromAddress common.Address
}

var wthetaAddress common.Address
var registrarOnMainchainAddress common.Address
var governanceTokenAddress common.Address
var mainchainTNT20TokenBankAddress common.Address
var subchainTNT20TokenBankAddress common.Address
var accountList []accounts
var mainchainTNT721TokenBankAddress common.Address
var mainchainTFuelTokenBankAddress common.Address
var subchainTFuelTokenBankAddress common.Address
var subchainTNT721TokenBankAddress common.Address
var subchainID *big.Int

func keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

func pubkeyToAddress(p ecdsa.PublicKey) common.Address {
	pubBytes := crypto.FromECDSAPub(&p)
	return common.BytesToAddress(keccak256(pubBytes[1:])[12:])
}

func DeployTokens() {
	// fmt.Printf("Deploying Tokens to the mainchain...\n\n")
	// mainchainClient, err := ethclient.Dial("http://localhost:18888/rpc")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// deployMockTNT20(mainchainClient)
	// deployMockTNT721(mainchainClient)
	// deployMockTNT1155(mainchainClient)
	// fmt.Println("")

	fmt.Printf("Deploying Tokens to the subchain...\n\n")
	subchainClient, err := ethclient.Dial("http://localhost:19888/rpc")
	if err != nil {
		log.Fatal(err)
	}
	deployMockTNT20(subchainClient)
	deployMockTNT721(subchainClient)
	deployMockTNT1155(subchainClient)
}

func deployMockTNT20(ethClient *ethclient.Client) common.Address {
	fmt.Printf("Deploying mock TNT20 token...\n")
	auth := subchainSelectAccount(ethClient, 1)
	address, _, _, err := ct.DeployMockTNT20(auth, ethClient)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mock TNT20 token deployed, Address:", address)
	// fmt.Println("Deployment tx:", tx.Hash().Hex())
	fmt.Println("")
	time.Sleep(6 * time.Second)
	return address
}

func deployMockTNT721(ethClient *ethclient.Client) common.Address {
	fmt.Printf("Deploying mock TNT721 token...\n")
	auth := subchainSelectAccount(ethClient, 1)
	address, _, _, err := ct.DeployMockTNT721(auth, ethClient)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mock TNT721 token deployed, Address:", address)
	// fmt.Println("Deployment tx:", tx.Hash().Hex())
	fmt.Println("")
	time.Sleep(6 * time.Second)
	return address
}

func deployMockTNT1155(ethClient *ethclient.Client) common.Address {
	fmt.Printf("Deploying mock TNT1155 token...\n")
	auth := subchainSelectAccount(ethClient, 1)
	address, _, _, err := ct.DeployMockTNT1155(auth, ethClient)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mock TNT1155 token deployed, Address:", address)
	// fmt.Println("Deployment tx:", tx.Hash().Hex())
	fmt.Println("")
	time.Sleep(6 * time.Second)
	return address
}

// func deploy_contracts() {
// 	subchainClient, err := ethclient.Dial("http://localhost:19888/rpc")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	auth := subchainSelectAccount(subchainClient, 1)
// 	address, tx, _, err := ct.DeployMockTNT20(auth, subchainClient)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("tnt20", address)
// 	fmt.Println(tx.Hash().Hex())

// 	subchainTNT20TokenAddress = address
// 	auth = subchainSelectAccount(subchainClient, 1)
// 	address, tx, _, err = ct.DeployMockTNT721(auth, subchainClient)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("tnt721", address)
// 	fmt.Println(tx.Hash().Hex())

// 	subchainTNT721TokenAddress = address
// 	auth = subchainSelectAccount(subchainClient, 1)
// 	address, tx, _, err = ct.DeployMockTNT1155(auth, subchainClient)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("tnt1155", address)
// 	fmt.Println(tx.Hash().Hex())

// 	subchainTNT1155TokenAddress = address
// }

func init() {
	subchainID = big.NewInt(360777)
	wthetaAddress = common.HexToAddress("0x7d73424a8256C0b2BA245e5d5a3De8820E45F390")
	registrarOnMainchainAddress = common.HexToAddress("0x08425D9Df219f93d5763c3e85204cb5B4cE33aAa")
	governanceTokenAddress = common.HexToAddress("0x6E05f58eEddA592f34DD9105b1827f252c509De0")

	mainchainTFuelTokenBankAddress = common.HexToAddress("0x7f1C87Bd3a22159b8a2E5D195B1a3283D10ea895")
	subchainTFuelTokenBankAddress = common.HexToAddress("0x5a443704dd4B594B382c22a083e2BD3090A6feF3")

	mainchainTNT20TokenBankAddress = common.HexToAddress("0x2Ce636d6240f8955d085a896e12429f8B3c7db26")
	subchainTNT20TokenBankAddress = common.HexToAddress("0x47e9Fbef8C83A1714F1951F142132E6e90F5fa5D")

	mainchainTNT721TokenBankAddress = common.HexToAddress("0xEd8d61f42dC1E56aE992D333A4992C3796b22A74")
	subchainTNT721TokenBankAddress = common.HexToAddress("0x8Be503bcdEd90ED42Eff31f56199399B2b0154CA")

	var map1 []string
	map1 = append(map1, "1111111111111111111111111111111111111111111111111111111111111111")
	map1 = append(map1, "93a90ea508331dfdf27fb79757d4250b4e84954927ba0073cd67454ac432c737")
	map1 = append(map1, "3333333333333333333333333333333333333333333333333333333333333333")
	map1 = append(map1, "4444444444444444444444444444444444444444444444444444444444444444")
	map1 = append(map1, "5555555555555555555555555555555555555555555555555555555555555555")
	map1 = append(map1, "6666666666666666666666666666666666666666666666666666666666666666")
	map1 = append(map1, "7777777777777777777777777777777777777777777777777777777777777777")
	map1 = append(map1, "8888888888888888888888888888888888888888888888888888888888888888")
	map1 = append(map1, "9999999999999999999999999999999999999999999999999999999999999999")
	map1 = append(map1, "1000000000000000000000000000000000000000000000000000000000000000")
	map1 = append(map1, "a249a82c42a282e87b2ddef63404d9dfcf6ea501dcaf5d447761765bd74f666d") //10
	map1 = append(map1, "d0d53ac0b4cd47d0ce0060dddc179d04145fea2ee2e0b66c3ee1699c6b492013") //11
	map1 = append(map1, "83f0bb8655139cef4657f90db64a7bb57847038a9bd0ccd87c9b0828e9cbf76d")

	// fmt.Println("-------------------------------------------------------- List of Accounts -------------------------------------------------------")
	for _, value := range map1 {

		privateKey, err := crypto.HexToECDSA(value)

		// privateKey, err := crypto.HexToECDSA("2dad160420b1e9b6fc152cd691a686a7080a0cee41b98754597a2ce57cc5dab1")
		if err != nil {
			log.Fatal(err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		fromAddress := pubkeyToAddress(*publicKeyECDSA)
		// fmt.Println("Private key:", value, "address:", fromAddress)
		accountList = append(accountList, accounts{priKey: value, privateKey: privateKey, fromAddress: fromAddress})
	}
	// fmt.Println("---------------------------------------------------------------------------------------------------------------------------------")
	// fmt.Println("")
	// deploy_contracts()
}

func mainchainSelectAccount(client *ethclient.Client, id int) *bind.TransactOpts {
	time.Sleep(1 * time.Second)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	//chainID := big.NewInt(360777)
	fromAddress := accountList[id].fromAddress
	privateKey := accountList[id].privateKey
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(crypto.ECDSAToPrivKey(privateKey), chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	// auth.Value = big.NewInt(20000000000000000000) // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	return auth
}

func subchainSelectAccount(client *ethclient.Client, id int) *bind.TransactOpts {
	time.Sleep(1 * time.Second)
	// chainID, err := client.ChainID(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	chainID := big.NewInt(360777)
	fromAddress := accountList[id].fromAddress
	privateKey := accountList[id].privateKey
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(crypto.ECDSAToPrivKey(privateKey), chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	// auth.Value = big.NewInt(20000000000000000000) // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	return auth
}

func subchainSelectAccountForchain(client *ethclient.Client, id int, chainid int64) *bind.TransactOpts {
	time.Sleep(1 * time.Second)
	// chainID, err := client.ChainID(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	chainID := big.NewInt(chainid)
	fromAddress := accountList[id].fromAddress
	privateKey := accountList[id].privateKey
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(crypto.ECDSAToPrivKey(privateKey), chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	// auth.Value = big.NewInt(20000000000000000000) // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	return auth
}

func OneAccountRegister(selected_subchainID *big.Int) {
	client, err := ethclient.Dial("http://localhost:18888/rpc")
	if err != nil {
		log.Fatal(err)
	}

	chainGuarantor := accountList[7].fromAddress

	instanceWrappedTheta, err := ct.NewMockWrappedTheta(wthetaAddress, client)
	if err != nil {
		log.Fatal("Failed to instantiate the wTHETA contract", err)
	}
	instanceChainRegistrar, err := ct.NewChainRegistrarOnMainchain(registrarOnMainchainAddress, client)
	if err != nil {
		log.Fatal("Failed to instantiate the ChainRegistrar contract", err)
	}
	var dec18 = new(big.Int)
	dec18.SetString("1000000000000000000", 10)
	amount := new(big.Int).Mul(dec18, big.NewInt(200000))

	auth := mainchainSelectAccount(client, 7) //chainGuarantor
	tx, err := instanceWrappedTheta.Mint(auth, chainGuarantor, amount)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Mint wTHETA tx hash (Mainchain):", tx.Hash().Hex())
	approveAmount := new(big.Int).Mul(dec18, big.NewInt(50000))
	authchainGuarantor := mainchainSelectAccount(client, 7)
	wThetaBalanceOfGuarantor, err := instanceWrappedTheta.BalanceOf(nil, chainGuarantor)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wTheta balance of %v: %v\n", chainGuarantor.Hex(), wThetaBalanceOfGuarantor)

	tx, err = instanceWrappedTheta.Approve(authchainGuarantor, registrarOnMainchainAddress, approveAmount)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Approve wTHETA tx hash (Mainchain):", tx.Hash().Hex())

	allChainIDs, _ := instanceChainRegistrar.GetAllSubchainIDs(nil)
	fmt.Printf("All subchain IDs before subchain registration: %v\n", allChainIDs)
	fmt.Printf("Registering subchain %v\n", selected_subchainID)
	collateralAmount := new(big.Int).Mul(dec18, big.NewInt(40000))
	authchainGuarantor = mainchainSelectAccount(client, 7)
	dummyGenesisHash := "0x012345679abcdef"
	tx, err = instanceChainRegistrar.RegisterSubchain(authchainGuarantor, selected_subchainID, governanceTokenAddress, collateralAmount, dummyGenesisHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Register Subchain tx hash (Mainchain):", tx.Hash().Hex())

	time.Sleep(12 * time.Second)
	allChainIDs, _ = instanceChainRegistrar.GetAllSubchainIDs(nil)
	fmt.Printf("Subchain registered, all subchain IDs: %v\n", allChainIDs)
}

func DSNRegister() {
	var DSNsubchains []*big.Int
	// DSNsubchains := [...]*big.Int{big.NewInt(360777), big.NewInt(360888)}
	for i := 1; i <= 16; i++ {
		chainIDBI, _ := new(big.Int).SetString(fmt.Sprintf("360%03d", i), 10)
		DSNsubchains = append(DSNsubchains, chainIDBI)
	}
	DSNsubchains = append(DSNsubchains, big.NewInt(360777))
	DSNsubchains = append(DSNsubchains, big.NewInt(360888))

	client, err := ethclient.Dial("http://localhost:18888/rpc")
	if err != nil {
		log.Fatal(err)
	}

	chainGuarantor := accountList[7].fromAddress

	instanceWrappedTheta, err := ct.NewMockWrappedTheta(wthetaAddress, client)
	if err != nil {
		log.Fatal("Failed to instantiate the wTHETA contract", err)
	}
	instanceChainRegistrar, err := ct.NewChainRegistrarOnMainchain(registrarOnMainchainAddress, client)
	if err != nil {
		log.Fatal("Failed to instantiate the ChainRegistrar contract", err)
	}
	var dec18 = new(big.Int)
	dec18.SetString("1000000000000000000", 10)
	amount := new(big.Int).Mul(dec18, big.NewInt(2000000000000))

	auth := mainchainSelectAccount(client, 7) //chainGuarantor
	tx, err := instanceWrappedTheta.Mint(auth, chainGuarantor, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mint: ", tx.Hash().Hex())
	// fmt.Println("Mint wTHETA tx hash (Mainchain):", tx.Hash().Hex())
	// approveAmount := new(big.Int).Mul(dec18, big.NewInt(50000))
	authchainGuarantor := mainchainSelectAccount(client, 7)
	wThetaBalanceOfGuarantor, err := instanceWrappedTheta.BalanceOf(nil, chainGuarantor)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wTheta balance of %v: %v\n", chainGuarantor.Hex(), wThetaBalanceOfGuarantor)

	tx, err = instanceWrappedTheta.Approve(authchainGuarantor, registrarOnMainchainAddress, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Approve: ", tx.Hash().Hex())
	// fmt.Println("Approve wTHETA tx hash (Mainchain):", tx.Hash().Hex())

	allChainIDs, _ := instanceChainRegistrar.GetAllSubchainIDs(nil)
	fmt.Printf("All subchain IDs before subchain registration: %v\n", allChainIDs)
	for _, sid := range DSNsubchains {
		fmt.Printf("Registering subchain %v\n", sid.String())
		collateralAmount := new(big.Int).Mul(dec18, big.NewInt(10001))
		authchainGuarantor = mainchainSelectAccount(client, 7)
		dummyGenesisHash := "0x012345679abcdef"
		tx, err := instanceChainRegistrar.RegisterSubchain(authchainGuarantor, sid, governanceTokenAddress, collateralAmount, dummyGenesisHash)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Register Subchain tx hash (Mainchain):", tx.Hash().Hex())
		time.Sleep(1 * time.Second)
	}
	time.Sleep(12 * time.Second)
	allChainIDs, _ = instanceChainRegistrar.GetAllSubchainIDs(nil)
	fmt.Printf("Subchain registered, all subchain IDs: %v\n", allChainIDs)
}

func prepareDSNMappings() map[string][]string {
	res := make(map[string][]string)
	//360001
	vs360001 := [...]string{"0x089514B8554FB9f3fd52047254f7cd5F06f45B27", "0x2778640bAC1eC93387EcEC9a40a034CF070768AC", "0xBE9B0fd263C33B934446EA5cA04dcc4cE01006e8", "0xF7A85ab52240D00a362DEF2A43BD968CA775A9E9"}
	res[big.NewInt(360001).String()] = vs360001[:]

	//360002
	vs360002 := [...]string{"0x0a74e562F8cF94Aa2B5675120b64f82a2d6F1e18", "0x65aa380A2277A3EC6cB2Ee86Cb8B81e407207BC4", "0x92e6De28a4bDcfd6C9C5Cc28b4FBE7C7e4ac276d", "0xaECbbcBe48d7cC7222494B3Dd7F6a570d7dd593c"}
	res[big.NewInt(360002).String()] = vs360002[:]

	//360003
	vs360003 := [...]string{"0x27E6ebAE0Fd53bcD9594Be2CBbEd51327dF4e9C9", "0xa09b0BDD35e4ef0412248a36eA9934830EfbD5E5", "0xf69A85c0f52164FA5E78B30E38C346099F495a07", "0xfba2581921f41B29d2FcC966Da6A886BDc6F5f84"}
	res[big.NewInt(360003).String()] = vs360003[:]

	//360004
	vs360004 := [...]string{"0x6c25030EB2c794308C03BA5f06070DC102a238C5", "0x8b17Bf7307e86C95A34941C9aF8aC500BEC55CdA", "0x920840b8AfcaABCe27De3234ad3B2cD0032e43C4", "0xE6Bf339BeC5a7605d7b6bC2306F42646a0155FC8"}
	res[big.NewInt(360004).String()] = vs360004[:]

	//360005
	vs360005 := [...]string{"0x7Ae1ceD9701b4D22f28C263A2F3565354277Ac76", "0xC62F63Ad0aE589B74771869999Aa4288170f19c3", "0xE30d5191f887d7d8964889c596912115E2CE366B", "0xfCb47e565a7E2B7Dfe4F01C52d8544c4032920b9"}
	res[big.NewInt(360005).String()] = vs360005[:]

	//360006
	vs360006 := [...]string{"0x0EcA08f33B5fd5869857b2b27fAF30C9A92CC2E6", "0x1EdE8fA5a8578888fd354c82F817c9c86ae293e8", "0x2eDb52449B96ebd1aa37036F8C5FFC4D7b811829", "0x71486d1799156eE3f937E8A04588a9F0aa6354ad"}
	res[big.NewInt(360006).String()] = vs360006[:]

	//360007
	vs360007 := [...]string{"0x29C4B88870dd5857D86595c9E28f90b2365A0B7f", "0x39530b33C7c9AbC77857ed3CC28F0d4C23Bc7ec6", "0x3BD98e246C5e48C9ec63d23873746Df151098104", "0x3F2A9af76d601Bf9bDd11FD3965c6E229bfDd665"}
	res[big.NewInt(360007).String()] = vs360007[:]

	//360008
	vs360008 := [...]string{"0x04FF6F40a0e91413031559adaf680fF620773fB7", "0x0eCfc6bCC5F1cb8643e669dda85592b938E96300", "0x8791882EB7dC0EbC7F51d01f95F067e54B85a65a", "0xbE89AbE4B208A8B796a232851b317B1646e16B05"}
	res[big.NewInt(360008).String()] = vs360008[:]

	//360009
	vs360009 := [...]string{"0x200282f2A9ae13A133f3Ab49A5607F67A0592c3d", "0x8A81cf2bD8e942db5E1D5164fa4d8105b72F9491", "0xAf743822E4adC26831EfBD361cA554831bF34D0e", "0xfd3755C3b66D2f8030a1e5e80e1f0ef2091DF52b"}
	res[big.NewInt(360009).String()] = vs360009[:]

	//360010
	vs360010 := [...]string{"0x022bdfB042611fb882D0710a4F0aE9f5cD9dEE84", "0x1F336cC04183850bb986B9396f852f1fce7EB644", "0x35EBb18F275330f4511036e5cbC4A3A51AF4809f", "0x9fe3827Ae8461f5d4e2905bB8CA574799Fec9161"}
	res[big.NewInt(360010).String()] = vs360010[:]

	//360011
	vs360011 := [...]string{"0x67081eD50352eC3524E7e3acd6791511dad5393e", "0x961bBC5806e034D4853a932EF8eaAfFeC530b2A4", "0xBeF9B766a89321D6806AA4D599CE629a6fbBab0C", "0xCB10b02f80a4B7Cf4ff919Ad2494b8190547DAA3"}
	res[big.NewInt(360011).String()] = vs360011[:]

	//360012
	vs360012 := [...]string{"0x4ABF9E41aDE70B0515EbdC64432E42ACdee1097F", "0xBC9F39161AeDf65cF66384C80fEf8b577Ae6beD2", "0xBcE748d50E794407F24452A55a225E86C0FBBEc4", "0xf6993FE260ee5004C06251acd80a03E7D6E27099"}
	res[big.NewInt(360012).String()] = vs360012[:]

	//360013
	vs360013 := [...]string{"0x0F898D64298c81d34741f29D9Bd59F452a9783ae", "0x7910967FCFAd74bBB3cF873dAC28aD8fC1675C94", "0xF209b73EAa566d9Fd7050BBB518427B511A92865", "0xF6113F645c0D6E52A7bCE5cDF6011264779277Fe"}
	res[big.NewInt(360013).String()] = vs360013[:]

	//360014
	vs360014 := [...]string{"0x02570A9338E3E08Af2d1d6F1A0D1993F479e709a", "0x28BE920062e90464b0F7a72A8DbdfbC1911CB9a6", "0x90688BE5b65B28C9B1d59276DE79399689b56B21", "0x9B1d3DF7D408a7FE194bEdB0a108255958fe9C89"}
	res[big.NewInt(360014).String()] = vs360014[:]

	//360015
	vs360015 := [...]string{"0x24dB4C8b2566896B2B3e0591867329CA0987c9Ec", "0x39673Af6d6f86FB9826EC7e5a3052EB6D50d0Fcf", "0x43A41Eb0BbA76500497A446aC1dc9A2A5Cb4f89E", "0xb1E6a392E83e1A24A22303E25cF9344BAA58e806"}
	res[big.NewInt(360015).String()] = vs360015[:]

	//360016
	vs360016 := [...]string{"0x6b942982537C08e997453ec1f0535461FEEa2fA0", "0x859BBe732710996360276a5A14F21E9940c239e5", "0x9A788c18D2A425De71d5c7c3a98f8F3697050C30", "0xdDDa75a5a880b64b90f4D3c57498f09db955DD12"}
	res[big.NewInt(360016).String()] = vs360016[:]

	//360777
	vs360777 := [...]string{"9220995f674b67d05f8ebc3643833aeafd421ec0", "33a027c2ac93b66987b3c8ea2bf5bd9f19e2a004", "9f99c71fd2cc01e748f96b74513f23050f56f564", "bf244894b0c6d9c18139cdc9a86ec501bdff6a26"}
	res[big.NewInt(360777).String()] = vs360777[:]
	//360888
	vs360888 := [...]string{"9220995f674b67d05f8ebc3643833aeafd421ec0", "33a027c2ac93b66987b3c8ea2bf5bd9f19e2a004", "9f99c71fd2cc01e748f96b74513f23050f56f564", "bf244894b0c6d9c18139cdc9a86ec501bdff6a26"}
	res[big.NewInt(360888).String()] = vs360888[:]
	return res
}

/*
func DSNStake() {
	validatorsInSubchains := prepareDSNMappings()

	client, err := ethclient.Dial("http://localhost:18888/rpc")
	if err != nil {
		log.Fatal(err)
	}
	var dec18 = new(big.Int)
	dec18.SetString("1000000000000000000", 10)
	instanceWrappedTheta, err := ct.NewMockWrappedTheta(wthetaAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	instanceGovernanceToken, err := ct.NewSubchainGovernanceToken(governanceTokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	instanceChainRegistrar, err := ct.NewChainRegistrarOnMainchain(registrarOnMainchainAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	validatorCollateralManagerAddr, _ := instanceChainRegistrar.Vcm(nil)
	validatorStakeManagerAddr, _ := instanceChainRegistrar.Vsm(nil)

	//
	// The guarantor deposits wTHETA collateral for the validator
	//

	fmt.Println("Prepare for validator collateral deposit...")

	validatorCollateral := new(big.Int).Mul(dec18, big.NewInt(2000))
	guarantor := mainchainSelectAccount(client, 1)
	tx, err := instanceWrappedTheta.Mint(guarantor, guarantor.From, new(big.Int).Mul(dec18, big.NewInt(200000000000000000)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Collateral mint tx hash (Mainchain): %v\n", tx.Hash().Hex())
	guarantor = mainchainSelectAccount(client, 1)
	tx, err = instanceWrappedTheta.Approve(guarantor, validatorCollateralManagerAddr, new(big.Int).Mul(dec18, big.NewInt(200000000000000000)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Collateral approve tx hash (Mainchain): %v\n", tx.Hash().Hex())

	for sid, addrs := range validatorsInSubchains {
		for _, addr := range addrs {
			guarantor = mainchainSelectAccount(client, 1)
			validator := common.HexToAddress(addr)
			bisid, ok := new(big.Int).SetString(sid, 10)
			if !ok {
				log.Fatal("Failed to parse bisid: %v")
			}
			tx, err = instanceChainRegistrar.DepositCollateral(guarantor, bisid, validator, validatorCollateral)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Collateral deposit tx hash (Mainchain): %v to address %v\n", tx.Hash().Hex(), validator.String())
			time.Sleep(1 * time.Second)
		}
	}
	time.Sleep(6 * time.Second)

	//
	// The staker deposits Gov token stakes for the validator
	//
	fmt.Println("Prepare for validator stake deposit...")
	staker := mainchainSelectAccount(client, 1)
	validatorStakingAmount := new(big.Int).Mul(dec18, big.NewInt(100000))
	validatorStakingAmountMint := new(big.Int).Mul(dec18, big.NewInt(10))
	// validatorStakingAmountMint.Mul(validatorStakingAmount, big.NewInt(10000000))
	authGovTokenInitDistrWallet := mainchainSelectAccount(client, 6)
	tx, err = instanceGovernanceToken.Transfer(authGovTokenInitDistrWallet, staker.From, validatorStakingAmountMint)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stake transfer tx hash (Mainchain): %v\n", tx.Hash().Hex())
	time.Sleep(6 * time.Second)

	staker = mainchainSelectAccount(client, 1)
	tx, err = instanceGovernanceToken.Approve(staker, validatorStakeManagerAddr, validatorStakingAmountMint)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stake approve tx hash (Mainchain): %v\n", tx.Hash().Hex())
	time.Sleep(12 * time.Second)

	minInitFeeRequired := new(big.Int).Mul(dec18, big.NewInt(100000)) // 100,000 TFuel
	for sid, addrs := range validatorsInSubchains {
		for _, addr := range addrs {
			staker = mainchainSelectAccount(client, 1)
			staker.Value.Set(minInitFeeRequired)
			validator := common.HexToAddress(addr)
			tx, err := instanceChainRegistrar.DepositStake(staker, sid, validator, validatorStakingAmount)
			staker.Value.Set(common.Big0)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Deposit %v Wei Gov Tokens as stake to subchain %v validator %v\n", validatorStakingAmount, sid, validator)
			fmt.Printf("Stake deposit tx hash (Mainchain): %v\n", tx.Hash().Hex())
			time.Sleep(1 * time.Second)
		}
		// time.Sleep(12 * time.Second)
		mainchainHeight, _ := client.BlockNumber(context.Background())
		dynasty := scom.CalculateDynasty(big.NewInt(int64(mainchainHeight)))
		fmt.Printf("Maichain block height: %v, dynasty: %v\n", mainchainHeight, dynasty)

		valset, _ := instanceChainRegistrar.GetValidatorSet(nil, sid, dynasty)
		fmt.Printf("Validator Set for subchain %v during the current dynasty %v: %v\n", sid, dynasty, valset)
		nextDynasty := big.NewInt(0).Add(dynasty, big.NewInt(1))
		valsetNext, _ := instanceChainRegistrar.GetValidatorSet(nil, sid, nextDynasty)
		fmt.Printf("Validator Set for subchain %v during the next dynasty    %v: %v\n", sid, nextDynasty, valsetNext)
	}

}
*/
func DSNStake() {
	id := 1
	validatorsInSubchains := prepareDSNMappings()
	// validator := common.HexToAddress(validatorAddrStr)
	validator := common.Address{}
	client, err := ethclient.Dial("http://localhost:18888/rpc")
	if err != nil {
		log.Fatal(err)
	}
	var dec18 = new(big.Int)
	dec18.SetString("1000000000000000000", 10)
	instanceWrappedTheta, err := ct.NewMockWrappedTheta(wthetaAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	instanceGovernanceToken, err := ct.NewSubchainGovernanceToken(governanceTokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	instanceChainRegistrar, err := ct.NewChainRegistrarOnMainchain(registrarOnMainchainAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	validatorCollateralManagerAddr, _ := instanceChainRegistrar.Vcm(nil)
	validatorStakeManagerAddr, _ := instanceChainRegistrar.Vsm(nil)

	//
	// The guarantor deposits wTHETA collateral for the validator
	//

	fmt.Println("Prepare for validator collateral deposit...")
	validatorCollateral := new(big.Int).Mul(dec18, big.NewInt(2000))
	guarantor := mainchainSelectAccount(client, 1)
	tx, err := instanceWrappedTheta.Mint(guarantor, guarantor.From, new(big.Int).Mul(dec18, big.NewInt(200000000000000000)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Collateral mint tx hash (Mainchain): %v\n", tx.Hash().Hex())
	guarantor = mainchainSelectAccount(client, 1)
	tx, err = instanceWrappedTheta.Approve(guarantor, validatorCollateralManagerAddr, new(big.Int).Mul(dec18, big.NewInt(200000000000000000)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Collateral approve tx hash (Mainchain): %v\n", tx.Hash().Hex())

	for sid, addrs := range validatorsInSubchains {
		for _, addr := range addrs {
			guarantor = mainchainSelectAccount(client, 1)
			validator := common.HexToAddress(addr)
			bisid, ok := new(big.Int).SetString(sid, 10)
			if !ok {
				log.Fatal("Failed to parse bisid: %v")
			}
			tx, err = instanceChainRegistrar.DepositCollateral(guarantor, bisid, validator, validatorCollateral)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Collateral deposit tx hash (Mainchain): %v to address %v\n", tx.Hash().Hex(), validator.String())
			time.Sleep(1 * time.Second)
		}
	}
	/*
		validatorCollateral := new(big.Int).Mul(dec18, big.NewInt(2000))
		guarantor := mainchainSelectAccount(client, id)
		tx, err := instanceWrappedTheta.Mint(guarantor, guarantor.From, big.NewInt(200000000000000000))
		if err != nil {
			log.Fatal(err)
		}
		guarantor = mainchainSelectAccount(client, id)
		tx, err = instanceWrappedTheta.Approve(guarantor, validatorCollateralManagerAddr, big.NewInt(200000000000000000))
		if err != nil {
			log.Fatal(err)
		}

		for sid, addrs := range validatorsInSubchains {
			for _, addr := range addrs {
				guarantor = mainchainSelectAccount(client, id)
				validator = common.HexToAddress(addr)
				selected_subchainID, ok := new(big.Int).SetString(sid, 10)
				if !ok {
					log.Fatal("Failed to parse bisid")
				}
				tx, err = instanceChainRegistrar.DepositCollateral(guarantor, selected_subchainID, validator, validatorCollateral)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Collateral deposit tx hash (Mainchain): %v to address %v on chain %v\n", tx.Hash().Hex(), validator.String(), selected_subchainID)
				time.Sleep(1 * time.Second)
			}
		}
	*/

	fmt.Println("Prepare for validator stake deposit...")

	//
	// The staker deposits Gov token stakes for the validator
	//

	staker := mainchainSelectAccount(client, id)
	validatorStakingAmount := new(big.Int).Mul(dec18, big.NewInt(100000))
	validatorStakingAmountMint := new(big.Int)
	validatorStakingAmountMint.Mul(validatorStakingAmount, big.NewInt(100))

	authGovTokenInitDistrWallet := mainchainSelectAccount(client, 6)
	tx, err = instanceGovernanceToken.Transfer(authGovTokenInitDistrWallet, staker.From, validatorStakingAmountMint)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(6 * time.Second)

	staker = mainchainSelectAccount(client, id)
	tx, err = instanceGovernanceToken.Approve(staker, validatorStakeManagerAddr, validatorStakingAmountMint)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(12 * time.Second)

	for sid, addrs := range validatorsInSubchains {
		selected_subchainID, ok := new(big.Int).SetString(sid, 10)
		if !ok {
			log.Fatal("Failed to parse bisid")
		}
		for _, addr := range addrs {
			staker = mainchainSelectAccount(client, id)
			minInitFeeRequired := new(big.Int).Mul(dec18, big.NewInt(100000)) // 100,000 TFuel
			staker.Value.Set(minInitFeeRequired)
			validator = common.HexToAddress(addr)
			tx, err = instanceChainRegistrar.DepositStake(staker, selected_subchainID, validator, validatorStakingAmount)
			// staker.Value.Set(common.Big0)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Deposit %v Wei Gov Tokens as stake to subchain %v validator %v\n", validatorStakingAmount, selected_subchainID, validator)
			fmt.Printf("Stake deposit tx hash (Mainchain): %v\n", tx.Hash().Hex())
			time.Sleep(1 * time.Second)
		}
	}

	time.Sleep(12 * time.Second)
	for sid, _ := range validatorsInSubchains {
		selected_subchainID, ok := new(big.Int).SetString(sid, 10)
		if !ok {
			log.Fatal("Failed to parse bisid")
		}
		mainchainHeight, _ := client.BlockNumber(context.Background())
		dynasty := scom.CalculateDynasty(big.NewInt(int64(mainchainHeight)))
		fmt.Printf("Maichain block height: %v, dynasty: %v\n", mainchainHeight, dynasty)

		valset, _ := instanceChainRegistrar.GetValidatorSet(nil, selected_subchainID, dynasty)
		fmt.Printf("Validator Set for subchain %v during the current dynasty %v: %v\n", selected_subchainID, dynasty, valset)

		nextDynasty := big.NewInt(0).Add(dynasty, big.NewInt(1))
		valsetNext, _ := instanceChainRegistrar.GetValidatorSet(nil, selected_subchainID, nextDynasty)
		fmt.Printf("Validator Set for subchain %v during the next dynasty    %v: %v\n", selected_subchainID, nextDynasty, valsetNext)
	}

}
func StakeToValidatorFromAccount(id int, validatorAddrStr string, selected_subchainID *big.Int) {
	validator := common.HexToAddress(validatorAddrStr)
	client, err := ethclient.Dial("http://localhost:18888/rpc")
	if err != nil {
		log.Fatal(err)
	}
	var dec18 = new(big.Int)
	dec18.SetString("1000000000000000000", 10)
	instanceWrappedTheta, err := ct.NewMockWrappedTheta(wthetaAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	instanceGovernanceToken, err := ct.NewSubchainGovernanceToken(governanceTokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	instanceChainRegistrar, err := ct.NewChainRegistrarOnMainchain(registrarOnMainchainAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	validatorCollateralManagerAddr, _ := instanceChainRegistrar.Vcm(nil)
	validatorStakeManagerAddr, _ := instanceChainRegistrar.Vsm(nil)

	//
	// The guarantor deposits wTHETA collateral for the validator
	//

	fmt.Println("Prepare for validator collateral deposit...")

	validatorCollateral := new(big.Int).Mul(dec18, big.NewInt(2000))
	guarantor := mainchainSelectAccount(client, id)
	tx, err := instanceWrappedTheta.Mint(guarantor, guarantor.From, validatorCollateral)
	if err != nil {
		log.Fatal(err)
	}
	guarantor = mainchainSelectAccount(client, id)
	tx, err = instanceWrappedTheta.Approve(guarantor, validatorCollateralManagerAddr, validatorCollateral)
	if err != nil {
		log.Fatal(err)
	}
	guarantor = mainchainSelectAccount(client, id)
	tx, err = instanceChainRegistrar.DepositCollateral(guarantor, selected_subchainID, validator, validatorCollateral)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Prepare for validator stake deposit...")

	//
	// The staker deposits Gov token stakes for the validator
	//

	staker := mainchainSelectAccount(client, id)
	validatorStakingAmount := new(big.Int).Mul(dec18, big.NewInt(100000))
	validatorStakingAmountMint := new(big.Int)
	validatorStakingAmountMint.Mul(validatorStakingAmount, big.NewInt(10))

	authGovTokenInitDistrWallet := mainchainSelectAccount(client, 6)
	tx, err = instanceGovernanceToken.Transfer(authGovTokenInitDistrWallet, staker.From, validatorStakingAmountMint)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(6 * time.Second)

	staker = mainchainSelectAccount(client, id)
	tx, err = instanceGovernanceToken.Approve(staker, validatorStakeManagerAddr, validatorStakingAmountMint)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(12 * time.Second)

	staker = mainchainSelectAccount(client, id)
	minInitFeeRequired := new(big.Int).Mul(dec18, big.NewInt(100000)) // 100,000 TFuel
	staker.Value.Set(minInitFeeRequired)
	tx, err = instanceChainRegistrar.DepositStake(staker, selected_subchainID, validator, validatorStakingAmount)
	staker.Value.Set(common.Big0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deposit %v Wei Gov Tokens as stake to subchain %v validator %v\n", validatorStakingAmount, selected_subchainID, validator)
	fmt.Printf("Stake deposit tx hash (Mainchain): %v\n", tx.Hash().Hex())

	time.Sleep(12 * time.Second)
	mainchainHeight, _ := client.BlockNumber(context.Background())
	dynasty := scom.CalculateDynasty(big.NewInt(int64(mainchainHeight)))
	fmt.Printf("Maichain block height: %v, dynasty: %v\n", mainchainHeight, dynasty)

	valset, _ := instanceChainRegistrar.GetValidatorSet(nil, selected_subchainID, dynasty)
	fmt.Printf("Validator Set for subchain %v during the current dynasty %v: %v\n", selected_subchainID, dynasty, valset)

	nextDynasty := big.NewInt(0).Add(dynasty, big.NewInt(1))
	valsetNext, _ := instanceChainRegistrar.GetValidatorSet(nil, selected_subchainID, nextDynasty)
	fmt.Printf("Validator Set for subchain %v during the next dynasty    %v: %v\n", selected_subchainID, nextDynasty, valsetNext)
}

func claimStake() {
	client, err := ethclient.Dial("http://localhost:18888/rpc")
	if err != nil {
		log.Fatal(err)
	}
	instanceChainRegistrar, err := ct.NewChainRegistrarOnMainchain(registrarOnMainchainAddress, client)
	if err != nil {
		log.Fatal("hhh", err)
	}
	var dec18 = new(big.Int)
	dec18.SetString("1000000000000000000", 10)
	var validatorStakingAmount = new(big.Int)

	//validatorStakingAmount := new(big.Int).Mul(dec18, big.NewInt(100000))
	// authGovTokenInitDistrWallet := mainchainSelectAccount(client, 6)
	// validator1 := accountList[12].fromAddress
	// instanceGovernanceToken, err := ct.NewSubchainGovernanceToken(governanceTokenAddress, client)
	if err != nil {
		log.Fatal("hhh", err)
	}
	// tx, err := instanceGovernanceToken.Transfer(authGovTokenInitDistrWallet, validator1, validatorStakingAmount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//validatorStakeManagerAddr, _ := instanceChainRegistrar.Vsm(nil)
	// authValidator1 := mainchainSelectAccount(client, 12) //Validator1
	// tx, err = instanceGovernanceToken.Approve(authValidator1, validatorStakeManagerAddr, validatorStakingAmount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// authValidator1 = mainchainSelectAccount(client, 12) //Validator1
	// tx, err = instanceChainRegistrar.DepositStake(authValidator1, subchainID, validator1, validatorStakingAmount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	height, _ := client.BlockNumber(context.Background())
	fmt.Println(big.NewInt(int64(height)))
	dynasty := big.NewInt(int64(height/100 + 1))
	re, _ := instanceChainRegistrar.GetValidatorSet(nil, subchainID, dynasty)
	fmt.Println(re)
	//validatorStakingAmount.SetString("1621710387562809105166", 10)
	validatorStakingAmount = re.ShareAmounts[0]
	authValidator1 := mainchainSelectAccount(client, 10) //Validator1
	tx, err := instanceChainRegistrar.WithdrawStake(authValidator1, subchainID, accountList[10].fromAddress, validatorStakingAmount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx.Hash().Hex())
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status != 1 {
		fmt.Println("error")
	}
	// validatorStakingAmount = re.ShareAmounts[1]
	// authValidator1 = mainchainSelectAccount(client, 10) //Validator1
	// tx, err = instanceChainRegistrar.WithdrawStake(authValidator1, subchainID, accountList[10].fromAddress, validatorStakingAmount)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(tx.Hash().Hex())
	// receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if receipt.Status != 1 {
	// 	fmt.Println("error")
	// }
}
