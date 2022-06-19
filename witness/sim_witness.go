package witness

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	scom "github.com/thetatoken/thetasubchain/common"
	score "github.com/thetatoken/thetasubchain/core"

	"github.com/thetatoken/theta/common"
)

const mainchainBlockIntervalMilliseconds int64 = 2000 // millseconds

// SimulatedMainchainWitness is a simulated mainchain witness for end-to-end testing
type SimulatedMainchainWitness struct {
	subchainID        string
	witnessedDynasty  *big.Int
	validatorSetCache map[string]*score.ValidatorSet
	updateTicker      *time.Ticker
	startingTime      time.Time

	// Life cycle
	wg                   *sync.WaitGroup
	ctx                  context.Context
	cancel               context.CancelFunc
	crossChainEventCache *score.InterChainEventCache

	lastSimEventNonce    map[score.InterChainMessageEventType]*big.Int
	hasTransferredTNT721 bool
	testId               int
}

// NewSimulatedMainchainWitness creates a new SimulatedMainchainWitness
func NewSimulatedMainchainWitness(
	ethClientAddress string,
	subchainID string,
	registerContractAddr common.Address,
	ercContractAddr common.Address,
	crossChainEventCache *score.InterChainEventCache,
	testId int,
) *SimulatedMainchainWitness {
	mw := &SimulatedMainchainWitness{
		subchainID:           subchainID,
		witnessedDynasty:     nil, // will be updated in the first update() call
		validatorSetCache:    make(map[string]*score.ValidatorSet),
		startingTime:         time.Now(),
		crossChainEventCache: crossChainEventCache,
		wg:                   &sync.WaitGroup{},
		// lastSimEventNonce:    big.NewInt(2),
		lastSimEventNonce:    make(map[score.InterChainMessageEventType]*big.Int),
		hasTransferredTNT721: false,
		testId:               testId,
	}
	mw.lastSimEventNonce[score.IMCEventTypeCrossChainTFuelTransfer] = common.Big1
	mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT20Transfer] = common.Big1
	mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT721Transfer] = common.Big1

	// amount := big.NewInt(88888)
	// nonce := big.NewInt(1)
	// mainchainBlockNumber := big.NewInt(0)
	// event := mw.generateInterChainEventForTFuelTransfer(amount, nonce, mainchainBlockNumber)

	// err := mw.crossChainEventCache.Insert(event)
	// if err != nil {
	// 	logger.Panicf("Insert Fail!! %v", err)
	// }
	return mw
}

func (mw *SimulatedMainchainWitness) Start(ctx context.Context) {
	c, cancel := context.WithCancel(ctx)
	mw.ctx = c
	mw.cancel = cancel

	mw.wg.Add(1)
	go mw.mainloop(ctx)
}

func (mw *SimulatedMainchainWitness) Stop() {
	if mw.updateTicker != nil {
		mw.updateTicker.Stop()
	}
	mw.cancel()
}

func (mw *SimulatedMainchainWitness) Wait() {
	mw.wg.Wait()
}

func (mw *SimulatedMainchainWitness) GetMainchainBlockNumber() (*big.Int, error) {
	blockNumber := int64((time.Since(mw.startingTime)).Milliseconds()) / mainchainBlockIntervalMilliseconds
	return big.NewInt(int64(blockNumber)), nil
}

func (mw *SimulatedMainchainWitness) GetMainchainBlockNumberUint() (uint64, error) {
	blockNumber, _ := mw.GetMainchainBlockNumber()
	return blockNumber.Uint64(), nil
}

func (mw *SimulatedMainchainWitness) GetValidatorSetByDynasty(dynasty *big.Int) (*score.ValidatorSet, error) {
	validatorSet, ok := mw.validatorSetCache[dynasty.String()]
	if ok && validatorSet != nil && validatorSet.Dynasty() == dynasty {
		return validatorSet, nil
	}

	var err error
	validatorSet, err = mw.updateValidatorSetCache(dynasty) // cache lazy update
	if err != nil {
		return nil, err
	}

	return validatorSet, nil
}

func (mw *SimulatedMainchainWitness) mainloop(ctx context.Context) {
	mw.updateTicker = time.NewTicker(time.Duration(1000) * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return
		case <-mw.updateTicker.C:
			mw.update()
		}
	}
}

func (mw *SimulatedMainchainWitness) update() {
	mainchainBlockNumber, err := mw.GetMainchainBlockNumber()
	logger.Debugf("witnessed main chain block height: %v", mainchainBlockNumber)

	if err != nil {
		logger.Warnf("failed to get the mainchain block number %v", err)
		return
	}

	dynasty := scom.CalculateDynasty(mainchainBlockNumber)
	if mw.witnessedDynasty == nil || dynasty.Cmp(mw.witnessedDynasty) > 0 { // needs to update the cache
		mw.updateValidatorSetCache(dynasty)
		mw.witnessedDynasty = dynasty
		logger.Infof("updated the witnessed dynasty to %v", dynasty)
	}
	if mw.testId == 0 {
		// TFuel cross-chain transfers
		for i := 0; i < 1; i++ {
			amount, ok := big.NewInt(0).SetString("82000000000000000000", 10) // 82 TFuel
			if !ok {
				logger.Panicf("failed to set amount %v", err)
			}
			event := mw.generateInterChainEventForTFuelTransfer(amount, big.NewInt(10), mainchainBlockNumber)
			mw.crossChainEventCache.Insert(event)
			// logger.Infof("Inserted Event %v", event)
		}
		return
	}

	// TFuel cross-chain transfers
	for i := 0; i < 3; i++ {
		amount, ok := big.NewInt(0).SetString("88000000000000000000", 10) // 88 TFuel
		if !ok {
			logger.Panicf("failed to set amount %v", err)
		}
		event := mw.generateInterChainEventForTFuelTransfer(amount, mw.lastSimEventNonce[score.IMCEventTypeCrossChainTFuelTransfer], mainchainBlockNumber)
		mw.crossChainEventCache.Insert(event)
		// logger.Infof("Inserted Event %v", event)
		mw.lastSimEventNonce[score.IMCEventTypeCrossChainTFuelTransfer] = new(big.Int).Add(mw.lastSimEventNonce[score.IMCEventTypeCrossChainTFuelTransfer], big.NewInt(1))
	}

	// TNT20 cross-chain transfers
	amount, ok := big.NewInt(0).SetString("66000000000000000000", 10) // 66 TDROP
	if !ok {
		logger.Panicf("failed to set amount %v", err)
	}
	event := mw.generateInterChainEventForTNT20Transfer(
		common.HexToAddress("0x1336739B05C7Ab8a526D40DCC0d04a826b5f8B03"), "TDrop", "TDROP", 18, amount,
		mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT20Transfer], mainchainBlockNumber)
	mw.crossChainEventCache.Insert(event)
	mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT20Transfer] = new(big.Int).Add(mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT20Transfer], big.NewInt(1))

	amount, ok = big.NewInt(0).SetString("9999999", 10)
	if !ok {
		logger.Panicf("failed to set amount %v", err)
	}
	event = mw.generateInterChainEventForTNT20Transfer(
		common.HexToAddress("0x15cc4c3f21417c392119054c8fe5895146e1a493"), "Random Token", "RTK", 6, amount,
		mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT20Transfer], mainchainBlockNumber)
	mw.crossChainEventCache.Insert(event)
	mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT20Transfer] = new(big.Int).Add(mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT20Transfer], big.NewInt(1))

	// TNT721 cross-chain transfers
	if !mw.hasTransferredTNT721 {
		event = mw.generateInterChainEventForTNT721Transfer(
			common.HexToAddress("0x0480c1097197831a1e4e9d64574f0048f8e35628"),
			"American Idol 20th Season Finalists",
			"AI20",
			big.NewInt(2076),
			"https://api.thetadrop.com/type/type_qyh516vms3hz4b24n8x8wcq3pgf.json?nft_id=nft_n5gr1291uge56ydf1cv0kvguaxzw",
			mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT721Transfer],
			mainchainBlockNumber,
		)
		mw.crossChainEventCache.Insert(event)
		mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT721Transfer] = new(big.Int).Add(mw.lastSimEventNonce[score.IMCEventTypeCrossChainTNT721Transfer], big.NewInt(1))
		mw.hasTransferredTNT721 = true
	}
}

func (mw *SimulatedMainchainWitness) updateValidatorSetCache(dynasty *big.Int) (*score.ValidatorSet, error) {
	// Simulate validator set updates
	validatorAddrList := []string{
		"0x2E833968E5bB786Ae419c4d13189fB081Cc43bab",
	}
	validatorSet := score.NewValidatorSet(dynasty)
	stake := big.NewInt(100000000)
	v := score.NewValidator(validatorAddrList[0], stake)
	validatorSet.AddValidator(v)
	// validatorAddrList := []string{
	// 	"0x9F1233798E905E173560071255140b4A8aBd3Ec6",
	// 	"0x2E833968E5bB786Ae419c4d13189fB081Cc43bab",
	// 	"0xC15E24083152dD76Ae6FC2aEb5269FF23d70330B",
	// 	"0x7631958d57Cf6a5605635a5F06Aa2ae2e000820e",
	// }
	// validatorSet := score.NewValidatorSet(dynasty)
	// stake1 := big.NewInt(100000000)
	// stake2 := big.NewInt(100000000)
	// stake3 := big.NewInt(100000000)
	// stake4 := big.NewInt(100000000)
	// v1 := score.NewValidator(validatorAddrList[0], stake1)
	// v2 := score.NewValidator(validatorAddrList[1], stake2)
	// v3 := score.NewValidator(validatorAddrList[2], stake3)
	// v4 := score.NewValidator(validatorAddrList[3], stake4)
	// validatorSet.AddValidator(v1)
	// validatorSet.AddValidator(v2)
	// validatorSet.AddValidator(v3)
	// validatorSet.AddValidator(v4)

	mw.validatorSetCache[dynasty.String()] = validatorSet

	logger.Infof("Witnessed validator set for dynasty %v", dynasty)
	for _, v := range validatorSet.Validators() {
		logger.Infof("Validator: %v", v)
	}

	return validatorSet, nil
}

func (mw *SimulatedMainchainWitness) GetInterChainEventCache() *score.InterChainEventCache {
	return mw.crossChainEventCache
}

func (mw *SimulatedMainchainWitness) generateInterChainEventForTFuelTransfer(amount *big.Int, nonce *big.Int, mainchainBlockNumber *big.Int) *score.InterChainMessageEvent {
	tfuelDenom := score.TFuelDenom(score.MainnetChainID)
	data, err := rlp.EncodeToBytes(score.TfuelTransferMetaData{
		TargetChainID:         big.NewInt(3610001),
		MainchainTokenSender:  common.HexToAddress("0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"),
		SubchainTokenReceiver: common.HexToAddress("0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"),
		Denom:                 tfuelDenom,
		LockedAmount:          amount,
		Nonce:                 nonce,
	})
	if err != nil {
		logger.Panicf("failed to get encode inter-chain message event data for TFuel transfer: %v", err)
	}

	event := score.NewInterChainMessageEvent(
		score.IMCEventTypeCrossChainTFuelTransfer,
		score.MainnetChainID,
		mw.subchainID,
		common.HexToAddress("0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"),
		common.HexToAddress("0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"),
		data,
		nonce,
		mainchainBlockNumber)

	return event
}

func (mw *SimulatedMainchainWitness) generateInterChainEventForTNT20Transfer(tokenSourceAddress common.Address,
	tokenName string, tokenSymbol string, tokenDecimals uint8, tokenAmount *big.Int, nonce *big.Int, mainchainBlockNumber *big.Int) *score.InterChainMessageEvent {
	tnt20Denom := score.TNT20Denom(score.MainnetChainID, tokenSourceAddress)
	targetChainId, _ := new(big.Int).SetString(mw.subchainID, 10)
	data, err := rlp.EncodeToBytes(score.TNT20TransferMetaData{
		TargetChainID:         targetChainId,
		MainchainTokenSender:  common.HexToAddress("0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"),
		SubchainTokenReceiver: common.HexToAddress("0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"),
		TNT20Contract:         tokenSourceAddress,
		Denom:                 tnt20Denom,
		Name:                  tokenName,
		Symbol:                tokenSymbol,
		Decimal:               tokenDecimals,
		Nonce:                 nonce,
		LockedAmount:          tokenAmount,
	})
	if err != nil {
		logger.Panicf("failed to get encode inter-chain message event data for TNT20 token transfer: %v", err)
	}

	event := score.NewInterChainMessageEvent(
		score.IMCEventTypeCrossChainTNT20Transfer,
		score.MainnetChainID,
		mw.subchainID,
		common.HexToAddress("0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"),
		common.HexToAddress("0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"),
		data,
		nonce,
		mainchainBlockNumber)

	return event
}

func (mw *SimulatedMainchainWitness) generateInterChainEventForTNT721Transfer(tokenSourceAddress common.Address,
	tokenName string, tokenSymbol string, tokenID *big.Int, tokenURI string, nonce *big.Int, mainchainBlockNumber *big.Int) *score.InterChainMessageEvent {
	tnt721Denom := score.TNT721Denom(score.MainnetChainID, tokenSourceAddress)
	data, err := rlp.EncodeToBytes(score.TNT721TransferMetaData{
		Denom:    tnt721Denom,
		Name:     tokenName,
		Symbol:   tokenSymbol,
		TokenID:  tokenID,
		TokenURI: tokenURI,
	})
	if err != nil {
		logger.Panicf("failed to get encode inter-chain message event data for TNT721 token transfer: %v", err)
	}

	event := score.NewInterChainMessageEvent(
		score.IMCEventTypeCrossChainTNT721Transfer,
		score.MainnetChainID,
		mw.subchainID,
		common.HexToAddress("0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"),
		common.HexToAddress("0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"),
		data,
		nonce,
		mainchainBlockNumber)

	return event
}
