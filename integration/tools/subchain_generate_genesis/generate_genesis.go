package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/ledger/types"
	"github.com/thetatoken/theta/rlp"
	"github.com/thetatoken/theta/store/database/backend"
	"github.com/thetatoken/theta/store/trie"

	score "github.com/thetatoken/thetasubchain/core"
	slst "github.com/thetatoken/thetasubchain/ledger/state"
)

var logger *log.Entry = log.WithFields(log.Fields{"prefix": "genesis"})

const (
	GenBlockHashMode int = iota
	GenGenesisFileMode
)

type Validator struct {
	Address string `json:"address"`
	Stake   string `json:"stake"`
}

//
// Example:
// cd $SUBCHAIN_HOME/integration/privatenet/node
// subchain_generate_genesis -chainID=private_subchain -initValidatorSet=./data/init_validator_set.json -genesis=./genesis
//
func main() {
	chainID, initValidatorSetPath, genesisSnapshotFilePath := parseArguments()

	sv, metadata, err := generateGenesisSnapshot(chainID, initValidatorSetPath, genesisSnapshotFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate genesis snapshot: %v", err))
	}

	err = sanityChecks(sv)
	if err != nil {
		panic(fmt.Sprintf("Sanity checks failed: %v", err))
	} else {
		logger.Infof("Sanity checks all passed.")
	}

	err = writeGenesisSnapshot(sv, metadata, genesisSnapshotFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to write genesis snapshot: %v", err))
	}

	genesisBlockHeader := metadata.TailTrio.Second.Header
	genesisBlockHash := genesisBlockHeader.Hash()

	fmt.Println("")
	fmt.Printf("-----------------------------------------------------------------------------------------\n")
	fmt.Printf("Genesis block hash: %v\n", genesisBlockHash.Hex())
	fmt.Printf("-----------------------------------------------------------------------------------------\n")
	fmt.Println("")
}

func parseArguments() (chainID, initValidatorSetPath, genesisSnapshotFilePath string) {
	chainIDPtr := flag.String("chainID", "private_subchain", "the ID of the chain")
	initValidatorSetPathPtr := flag.String("initValidatorSet", "./init_validator_set.json", "the initial validator set")
	genesisSnapshotFilePathPtr := flag.String("genesis", "./genesis", "the genesis snapshot")
	flag.Parse()

	chainID = *chainIDPtr
	initValidatorSetPath = *initValidatorSetPathPtr
	genesisSnapshotFilePath = *genesisSnapshotFilePathPtr

	return
}

// generateGenesisSnapshot generates the genesis snapshot.
func generateGenesisSnapshot(chainID, initValidatorSetFilePath, genesisSnapshotFilePath string) (*slst.StoreView, *score.SnapshotMetadata, error) {
	metadata := &score.SnapshotMetadata{}
	genesisHeight := score.GenesisBlockHeight

	sv := slst.NewStoreView(0, common.Hash{}, backend.NewMemDatabase())
	setInitialValidatorSet(initValidatorSetFilePath, genesisHeight, sv)
	setInitalEventNonce(sv)

	stateHash := sv.Hash()

	genesisBlock := score.NewBlock()
	genesisBlock.ChainID = chainID
	genesisBlock.Height = genesisHeight
	genesisBlock.Epoch = genesisBlock.Height
	genesisBlock.Parent = common.Hash{}
	genesisBlock.StateHash = stateHash
	genesisBlock.Timestamp = big.NewInt(time.Now().Unix())

	metadata.TailTrio = score.SnapshotBlockTrio{
		First:  score.SnapshotFirstBlock{},
		Second: score.SnapshotSecondBlock{Header: genesisBlock.BlockHeader},
		Third:  score.SnapshotThirdBlock{},
	}

	return sv, metadata, nil
}

func setInitialValidatorSet(initValidatorSetFilePath string, genesisHeight uint64, sv *slst.StoreView) *score.ValidatorSet {
	var validators []Validator
	initValidatorSetFile, err := os.Open(initValidatorSetFilePath)
	if err != nil {
		panic(fmt.Sprintf("failed to open initial stake deposit file: %v", err))
	}
	initValidatorSetByteValue, err := ioutil.ReadAll(initValidatorSetFile)
	if err != nil {
		panic(fmt.Sprintf("failed to read initial stake deposit file: %v", err))
	}

	json.Unmarshal(initValidatorSetByteValue, &validators)
	initialDynasty := big.NewInt(0)
	validatorSet := score.NewValidatorSet(initialDynasty)
	for _, v := range validators {
		stake, ok := big.NewInt(0).SetString(v.Stake, 10)
		if !ok {
			panic(fmt.Sprintf("failed to parse stake: %v", v.Stake))
		}
		validator := score.NewValidator(v.Address, stake)
		validatorSet.AddValidator(validator)
	}

	sv.UpdateValidatorSet(validatorSet)

	hl := &types.HeightList{}
	hl.Append(genesisHeight)
	sv.UpdateValidatorSetUpdateTxHeightList(hl)
	sv.Save()

	return validatorSet
}

func setInitalEventNonce(sv *slst.StoreView) *big.Int {
	sv.UpdateLastProcessedEventNonce(big.NewInt(0))
	sv.Save()

	return big.NewInt(0)
}

func proveValidatorSet(sv *slst.StoreView) (*score.ValidatorSetProof, error) {
	vp := &score.ValidatorSetProof{}
	vsKey := slst.ValidatorSetKey()
	err := sv.ProveValidatorSet(vsKey, vp)
	return vp, err
}

// writeGenesisSnapshot writes genesis snapshot to file system.
func writeGenesisSnapshot(sv *slst.StoreView, metadata *score.SnapshotMetadata, genesisSnapshotFilePath string) error {
	file, err := os.Create(genesisSnapshotFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	err = score.WriteMetadata(writer, metadata)
	if err != nil {
		return err
	}
	writeStoreView(sv, true, writer)
	return err
}

func writeStoreView(sv *slst.StoreView, needAccountStorage bool, writer *bufio.Writer) {
	height := score.Itobytes(sv.Height())
	err := score.WriteRecord(writer, []byte{score.SVStart}, height)
	if err != nil {
		panic(err)
	}
	sv.GetStore().Traverse(nil, func(k, v common.Bytes) bool {
		err = score.WriteRecord(writer, k, v)
		if err != nil {
			panic(err)
		}
		return true
	})
	err = score.WriteRecord(writer, []byte{score.SVEnd}, height)
	if err != nil {
		panic(err)
	}
	writer.Flush()
}

func sanityChecks(sv *slst.StoreView) error {
	vsAnalyzed := false
	sv.GetStore().Traverse(nil, func(key, val common.Bytes) bool {
		if bytes.Equal(key, slst.ValidatorSetKey()) {
			var vs score.ValidatorSet
			err := rlp.DecodeBytes(val, &vs)
			if err != nil {
				panic(fmt.Sprintf("Failed to decode validator set: %v", err))
			}
			for _, validator := range vs.Validators() {
				logger.Infof("------------------------------------------------------------------------------")
				logger.Infof("Initial validator: %v, stake: %v", validator.Address.Hex(), validator.Stake)
				logger.Infof("------------------------------------------------------------------------------")
			}
			vsAnalyzed = true
		} else if bytes.Equal(key, slst.ValidatorSetUpdateTxHeightListKey()) {
			var hl types.HeightList
			err := rlp.DecodeBytes(val, &hl)
			if err != nil {
				panic(fmt.Sprintf("Failed to decode Height List: %v", err))
			}
			if len(hl.Heights) != 1 {
				panic(fmt.Sprintf("The genesis height list should contain only one height: %v", hl.Heights))
			}
			if hl.Heights[0] != uint64(0) {
				panic("Only height 0 should be in the genesis height list")
			}
		}
		return true
	})

	// Sanity checks for the initial validator set
	vsProof, err := proveValidatorSet(sv)
	if err != nil {
		panic("Failed to get VS proof from storeview")
	}
	_, _, err = trie.VerifyProof(sv.Hash(), slst.ValidatorSetKey(), vsProof)
	if err != nil {
		panic("Failed to verify VS proof in storeview")
	}
	if !vsAnalyzed {
		return fmt.Errorf("initial validator set not detected in the genesis file")
	}

	return nil
}

