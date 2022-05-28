package core

import (
	// "bytes"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"

	// "sort"

	// log "github.com/sirupsen/logrus"

	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/crypto"
	"github.com/thetatoken/theta/rlp"
	ts "github.com/thetatoken/theta/store"
	"github.com/thetatoken/theta/store/database"
	"github.com/thetatoken/theta/store/kvstore"

	scta "github.com/thetatoken/thetasubchain/contracts/accessors"
	"github.com/thetatoken/thetasubchain/eth/abi"
)

// var logger *log.Entry = log.WithFields(log.Fields{"prefix": "core"})

type InterChainMessageEventType uint64

const (
	IMCEventTypeUnknown                  InterChainMessageEventType = 0
	IMCEventTypeCrossChainTransfer       InterChainMessageEventType = 1
	IMCEventTypeCrossChainTFuelTransfer  InterChainMessageEventType = 10001
	IMCEventTypeCrossChainTNT20Transfer  InterChainMessageEventType = 10002
	IMCEventTypeCrossChainTNT721Transfer InterChainMessageEventType = 10003

	IMCEventTypeVoucherBurn       InterChainMessageEventType = 2
	IMCEventTypeVoucherBurnTFuel  InterChainMessageEventType = 20001
	IMCEventTypeVoucherBurnTNT20  InterChainMessageEventType = 20002
	IMCEventTypeVoucherBurnTNT721 InterChainMessageEventType = 20003

	IMCEventLock   InterChainMessageEventType = 3
	IMCEventUnLock InterChainMessageEventType = 4
)

// InterChainMessageEvent contains the public information of a crosschain transfer event.
type InterChainMessageEvent struct {
	Type          InterChainMessageEventType
	SourceChainID string
	TargetChainID string
	Sender        common.Address
	Receiver      common.Address
	Data          common.Bytes // generic data field that can be used to encode arbitrary data for inter-chain messaging
	Nonce         *big.Int
	BlockNumber   *big.Int
}

// NewInterChainMessageEvent creates a new crosschain transfer event instance.
func NewInterChainMessageEvent(eventType InterChainMessageEventType, sourceChainID string, targetChainID string, sender common.Address, receiver common.Address,
	data common.Bytes, nonce *big.Int, blockNumber *big.Int) *InterChainMessageEvent {
	return &InterChainMessageEvent{eventType, sourceChainID, targetChainID, sender, receiver, data, nonce, blockNumber}
}

// ID returns the ID of the crosschain transfer event, which is the string representation of its address.
func (c *InterChainMessageEvent) ID() *big.Int {
	return c.Nonce
}

func (c *InterChainMessageEvent) NewID() string {
	return strconv.FormatUint(uint64(c.Type), 10) + "/" + c.Nonce.String()
}

// Equals checks whether the crosschain transfer event is the same as another crosschain transfer event
func (c *InterChainMessageEvent) Equals(x *InterChainMessageEvent) bool {
	if c.Type != x.Type {
		return false
	}
	if c.SourceChainID != x.SourceChainID {
		return false
	}
	if c.TargetChainID != x.TargetChainID {
		return false
	}
	if c.Nonce.Cmp(x.Nonce) != 0 {
		return false
	}
	if c.Sender.Hex() != x.Sender.Hex() {
		return false
	}
	if c.Receiver.Hex() != x.Receiver.Hex() {
		return false
	}
	if !bytes.Equal(c.Data, x.Data) {
		return false
	}
	if c.BlockNumber.Cmp(x.BlockNumber) != 0 {
		return false
	}
	return true
}

// String represents the string representation of the validator
func (c *InterChainMessageEvent) String() string {
	return fmt.Sprintf("{ID: %v, Type: %v, SourceChainID: %v, TargetChainID: %v, Sender: %v, Receiver: %v,  Data: %v, Nonce: %v, BlockNumber: %v}",
		c.ID(), c.Type, c.SourceChainID, c.TargetChainID, c.Sender.Hex(), c.Receiver.Hex(), string(c.Data), c.Nonce.String(), c.BlockNumber.String())
}

// ByID implements sort.Interface for InterChainMessageEvent based on ID (Nonce).
type ICMEByID []InterChainMessageEvent

func (b ICMEByID) Len() int           { return len(b) }
func (b ICMEByID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ICMEByID) Less(i, j int) bool { return b[i].Nonce.Cmp(b[j].Nonce) < 0 }

var _ rlp.Encoder = (*InterChainMessageEvent)(nil)

// EncodeRLP implements RLP Encoder interface.
func (c *InterChainMessageEvent) EncodeRLP(w io.Writer) error {
	if c == nil {
		return rlp.Encode(w, &InterChainMessageEvent{})
	}
	return rlp.Encode(w, []interface{}{
		c.Type,
		c.SourceChainID,
		c.TargetChainID,
		c.Sender,
		c.Receiver,
		c.Data,
		c.Nonce,
		c.BlockNumber,
	})
}

// DecodeRLP implements RLP Decoder interface.
func (c *InterChainMessageEvent) DecodeRLP(stream *rlp.Stream) error {
	_, err := stream.List()
	if err != nil {
		return err
	}

	eventType := InterChainMessageEventType(0)
	err = stream.Decode(&eventType)
	if err != nil {
		return err
	}
	c.Type = eventType

	sourceChainID := ""
	err = stream.Decode(&sourceChainID)
	if err != nil {
		return err
	}
	c.SourceChainID = sourceChainID

	targetChainID := ""
	err = stream.Decode(&targetChainID)
	if err != nil {
		return err
	}

	sender := &common.Address{}
	err = stream.Decode(sender)
	if err != nil {
		return err
	}
	c.Sender = *sender

	receiver := &common.Address{}
	err = stream.Decode(receiver)
	if err != nil {
		return err
	}
	c.Receiver = *receiver

	data := common.Bytes{}
	err = stream.Decode(&data)
	if err != nil {
		return err
	}
	c.Data = data

	nonce := big.NewInt(0)
	err = stream.Decode(nonce)
	if err != nil {
		return err
	}
	c.Nonce = nonce

	blockNumber := big.NewInt(0)
	err = stream.Decode(blockNumber)
	if err != nil {
		return err
	}
	c.BlockNumber = blockNumber

	return stream.ListEnd()
}

// ------------------------------------ Inter-Chain Event Cache ----------------------------------------------

var (
	// ErrInterChainMessageEventrNotFound for ID is not found in crosschain transfer event set.
	ErrInterChainMessageEventNotFound      = errors.New("InterChainMessageEventNotFound")
	ErrInterChainMessageEventExisted       = errors.New("InterChainMessageEventrExisted")
	ErrInterChainMessageEventPersistFailed = errors.New("InterChainMessageEventPersistFailed")
)

// InterChainEventIndexKey constructs the DB key for the given block hash.
func InterChainEventIndexKey(ICMEtype InterChainMessageEventType, nonce *big.Int) common.Bytes {
	return common.Bytes("ice/" + strconv.FormatUint(uint64(ICMEtype), 10) + "/" + nonce.String())
}

func InterChainTransferEventNextNonceKey(ICMEtype InterChainMessageEventType) common.Bytes {
	return common.Bytes("ictenn/" + strconv.FormatUint(uint64(ICMEtype), 10))
}

func LastQueryedHeightKey(ICMEtype InterChainMessageEventType) common.Bytes {
	switch ICMEtype {
	case IMCEventTypeCrossChainTransfer:
		return common.Bytes("tlq")
	case IMCEventTypeVoucherBurn:
		return common.Bytes("vblq")
	default:
		return nil
	}
}

func LastProcessedUnfinalizedVoucherBurnNonceKey(ICMEtype InterChainMessageEventType) common.Bytes {
	return common.Bytes("vblp" + strconv.FormatUint(uint64(ICMEtype), 10))
}

func VoucherBurnStatusInfoKey(ICMEtype InterChainMessageEventType, nonce *big.Int) common.Bytes {
	return common.Bytes("vbsi" + strconv.FormatUint(uint64(ICMEtype), 10) + "/" + nonce.String())
}

type InterChainEventCache struct {
	mutex            *sync.Mutex // mutex to for concurrency protection, e.g., the witness thread and consensus thread may access it concurrently
	db               database.Database
	EventSelectors   map[InterChainMessageEventType]string
	TransferTypes    []InterChainMessageEventType
	VoucherBurnTypes []InterChainMessageEventType
}

// NewInterChainEventCache creates a new crosschain transfer event cache instance.
func NewInterChainEventCache(db database.Database) *InterChainEventCache {
	cache := &InterChainEventCache{
		mutex:            &sync.Mutex{},
		db:               db,
		EventSelectors:   make(map[InterChainMessageEventType]string),
		TransferTypes:    []InterChainMessageEventType{IMCEventTypeCrossChainTFuelTransfer, IMCEventTypeCrossChainTNT20Transfer, IMCEventTypeCrossChainTNT721Transfer},
		VoucherBurnTypes: []InterChainMessageEventType{IMCEventTypeVoucherBurnTFuel, IMCEventTypeVoucherBurnTNT20, IMCEventTypeVoucherBurnTNT721},
	}
	cache.createEventSelectors()
	return cache
}

func (c *InterChainEventCache) createEventSelectors() {
	c.EventSelectors[IMCEventTypeCrossChainTFuelTransfer] = crypto.Keccak256Hash([]byte("TFeulTokenLocked(uint256,address,address,uint256,uint256,string)")).Hex()
	c.EventSelectors[IMCEventTypeCrossChainTNT20Transfer] = crypto.Keccak256Hash([]byte("TNT20TokenLocked(uint256,address,address,uint256,address,string,string,uint8,uint256,string)")).Hex()
	// c.EventSelectors[IMCEventTypeCrossChainTNT721Transfer] = crypto.Keccak256Hash([]byte("")).Hex()
	c.EventSelectors[IMCEventTypeVoucherBurnTFuel] = crypto.Keccak256Hash([]byte("BurnTFuelVouchers(address,address,uint256,uint256)")).Hex()
	c.EventSelectors[IMCEventTypeVoucherBurnTNT20] = crypto.Keccak256Hash([]byte("BurnTNT20Vouchers(string,address,address,uint256,uint256)")).Hex()
	c.EventSelectors[IMCEventTypeVoucherBurnTNT721] = crypto.Keccak256Hash([]byte("BurnTNT721Vouchers(string,address,address,uint256,uint256)")).Hex()
}

func (c *InterChainEventCache) Insert(event *InterChainMessageEvent) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	store := kvstore.NewKVStore(c.db)
	err := store.Put(InterChainEventIndexKey(event.Type, event.Nonce), event)
	return err // the caller should handle the error
}

func (c *InterChainEventCache) Delete(IMCEtype InterChainMessageEventType, nonce *big.Int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	store := kvstore.NewKVStore(c.db)
	err := store.Delete(InterChainEventIndexKey(IMCEtype, nonce))
	return err // the caller should handle the error
}

func (c *InterChainEventCache) Get(IMCEtype InterChainMessageEventType, nonce *big.Int) (*InterChainMessageEvent, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	event := InterChainMessageEvent{}
	store := kvstore.NewKVStore(c.db)
	err := store.Get(InterChainEventIndexKey(IMCEtype, nonce), &event)
	return &event, err // the caller should handle the error
}

func (c *InterChainEventCache) Exists(IMCEtype InterChainMessageEventType, nonce *big.Int) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	event := InterChainMessageEvent{}
	store := kvstore.NewKVStore(c.db)
	err := store.Get(InterChainEventIndexKey(IMCEtype, nonce), &event)
	if err == nil {
		return true, nil
	}

	if err == ts.ErrKeyNotFound {
		return false, nil
	}

	return false, err // the caller should handle the error
}

// RPC related

type LogData struct {
	LogIndex         string   `json:"logIndex"`
	TransactionIndex string   `json:"transactionIndex"`
	TransactionHash  string   `json:"transactionHash"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Address          string   `json:"address"`
	Data             string   `json:"data"`
	Topics           []string `json:"topics"`
	Type             string   `json:"type"`
}

type RPCResult struct {
	Jsonrpc string    `json:"jsonrpc"`
	Id      int64     `json:"id"`
	Result  []LogData `json:"result"`
}

type TransferEvent struct {
	Denom  string
	Amount *big.Int
	Nonce  *big.Int
}

func (c *InterChainEventCache) RpcEventLogQuery(fromBlock *big.Int, toBlock *big.Int, contractAddr common.Address, IMCEType InterChainMessageEventType) {
	url := "http://127.0.0.1:18888/rpc"
	queryStr := fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"eth_getLogs",
		"params":[{"fromBlock":"%v","toBlock":"%v", "address":"%v","topics":["%v"]}],
		"id":74
	}`, fmt.Sprintf("%x", fromBlock), fmt.Sprintf("%x", toBlock), contractAddr.Hex(), c.EventSelectors[IMCEType])
	var jsonData = []byte(queryStr)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Fatal(err)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	var rpcres RPCResult
	err = json.Unmarshal(body, &rpcres)
	if err != nil {
		logger.Fatal(err)
	}

	for _, logData := range rpcres.Result {
		logData := logData
		data, _ := hex.DecodeString(logData.Data[2:])
		switch IMCEType {
		case IMCEventTypeCrossChainTFuelTransfer:
			var tma TfuelTransferMetaData
			contractAbi, _ := abi.JSON(strings.NewReader(string(scta.MainchainTFuelTokenBankABI)))
			contractAbi.UnpackIntoInterface(&tma, "TFeulTokenLocked", data)
			sourceChainID, _ := ExtractSourceChainIDFromDenom(tma.Denom)
			blockNumber, _ := new(big.Int).SetString(logData.BlockNumber, 10)
			event := &InterChainMessageEvent{
				Type:          IMCEventTypeCrossChainTFuelTransfer,
				SourceChainID: sourceChainID,
				TargetChainID: tma.TargetChainID.String(),
				Sender:        tma.MainchainTokenSender,
				Receiver:      tma.SubchainTokenReceiver,
				Data:          data,
				Nonce:         tma.Nonce,
				BlockNumber:   blockNumber,
			}
			c.Insert(event)
		case IMCEventTypeCrossChainTNT20Transfer:
			var tma TNT20TransferMetaData
			contractAbi, _ := abi.JSON(strings.NewReader(string(scta.MainchainTNT20TokenBankABI)))
			contractAbi.UnpackIntoInterface(&tma, "TNT20TokenLocked", data)
			sourceChainID, _ := ExtractSourceChainIDFromDenom(tma.Denom)
			blockNumber, _ := new(big.Int).SetString(logData.BlockNumber, 10)
			event := &InterChainMessageEvent{
				Type:          IMCEventTypeCrossChainTFuelTransfer,
				SourceChainID: sourceChainID,
				TargetChainID: tma.TargetChainID.String(),
				Sender:        tma.MainchainTokenSender,
				Receiver:      tma.SubchainTokenReceiver,
				Data:          data,
				Nonce:         tma.Nonce,
				BlockNumber:   blockNumber,
			}
			c.Insert(event)
		case IMCEventTypeCrossChainTNT721Transfer:
		default:
		}
	}
}

// ------------------------------------ Getters and Setters for utility values --------------------------------------------

func (c *InterChainEventCache) GetLastQueryedHeightForType(ICMEtype InterChainMessageEventType) (*big.Int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	height := big.NewInt(0)
	store := kvstore.NewKVStore(c.db)
	err := store.Get(LastQueryedHeightKey(ICMEtype), &height)
	if err == nil {
		return height, nil
	}
	return big.NewInt(0), err
}

func (c *InterChainEventCache) SetLastQueryedHeightForType(ICMEtype InterChainMessageEventType, height *big.Int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	store := kvstore.NewKVStore(c.db)
	err := store.Put(LastQueryedHeightKey(ICMEtype), height)
	return err // the caller should handle the error
}

func (c *InterChainEventCache) GetLastProcessedUnfinalizedVoucherBurnNonce(ICMEtype InterChainMessageEventType) (*big.Int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	nonce := big.NewInt(0)
	store := kvstore.NewKVStore(c.db)
	err := store.Get(LastProcessedUnfinalizedVoucherBurnNonceKey(ICMEtype), &nonce)
	return nonce, err
}

func (c *InterChainEventCache) SetLastProcessedUnfinalizedVoucherBurnNonce(ICMEtype InterChainMessageEventType, nonce *big.Int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	store := kvstore.NewKVStore(c.db)
	err := store.Put(LastProcessedUnfinalizedVoucherBurnNonceKey(ICMEtype), nonce)
	return err // the caller should handle the error
}

func (c *InterChainEventCache) GetVoucherBurnStatus(ICMEtype InterChainMessageEventType, nonce *big.Int) (*VoucherBurnEventStatusInfo, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	statusInfo := VoucherBurnEventStatusInfo{}
	store := kvstore.NewKVStore(c.db)
	err := store.Get(VoucherBurnStatusInfoKey(ICMEtype, nonce), &statusInfo)
	return &statusInfo, err
}

func (c *InterChainEventCache) SetVoucherBurnStatus(statusInfo *VoucherBurnEventStatusInfo) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	store := kvstore.NewKVStore(c.db)
	err := store.Put(VoucherBurnStatusInfoKey(statusInfo.Type, statusInfo.Nonce), statusInfo)
	return err // the caller should handle the error
}

func (c *InterChainEventCache) VoucherBurnNonceExists(IMCEtype InterChainMessageEventType, nonce *big.Int) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	event := InterChainMessageEvent{}
	store := kvstore.NewKVStore(c.db)
	err := store.Get(VoucherBurnStatusInfoKey(IMCEtype, nonce), &event)
	if err == nil {
		return true, nil
	}

	if err == ts.ErrKeyNotFound {
		return false, nil
	}

	return false, err // the caller should handle the error
}

func (c *InterChainEventCache) GetNextTransferNonceForType(ICMEtype InterChainMessageEventType) (*big.Int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	nonce := big.NewInt(0)
	store := kvstore.NewKVStore(c.db)
	err := store.Get(InterChainTransferEventNextNonceKey(ICMEtype), &nonce)
	return nonce, err
}

func (c *InterChainEventCache) SetNextTransferNonceForType(ICMEtype InterChainMessageEventType, nonce *big.Int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	store := kvstore.NewKVStore(c.db)
	err := store.Put(InterChainTransferEventNextNonceKey(ICMEtype), nonce)
	return err // the caller should handle the error
}

// ------------------------------------ Cross-Chain Voucher Burn --------------------------------------------
type VoucherBurnEventStatus byte

const (
	VoucherBurnEventStatusPending VoucherBurnEventStatus = VoucherBurnEventStatus(iota)
	VoucherBurnEventStatusProcessed
	VoucherBurnEventStatusFinalized
	VoucherBurnEventStatusFailed
)

type VoucherBurnEventStatusInfo struct {
	Type                     InterChainMessageEventType
	Nonce                    *big.Int
	Status                   VoucherBurnEventStatus
	LastProcessedBlockHeight *big.Int
	RetriedTime              uint
}

type VoucherBurnData struct {
	SubchainID  *big.Int
	Dynasty     *big.Int
	TxHash      common.Hash
	Denom       string
	Receiver    common.Address
	Amount      *big.Int
	TFuelAmount *big.Int
}

type TfuelVoucherBurnMetaData struct {
	TxHash common.Hash
	Denom  string
	Amount *big.Int
}

type TNT20VoucherBurnMetaData struct {
	TxHash  common.Hash
	Denom   string
	Amount  *big.Int
	TokenId *big.Int
}

// ------------------------------------ Cross-Chain Asset Transfer --------------------------------------------

// Cross-Chain TFuel Transfer

type TfuelTransferMetaData struct {
	TargetChainID         *big.Int
	MainchainTokenSender  common.Address
	SubchainTokenReceiver common.Address
	LockedAmount          *big.Int
	Nonce                 *big.Int
	Denom                 string
}

type CrossChainTFuelTransferEvent struct {
	Sender      common.Address
	Receiver    common.Address
	Denom       string
	Amount      *big.Int
	Nonce       *big.Int
	BlockNumber *big.Int
}

func NewCrossChainTFuelTransferEvent(sender common.Address, receiver common.Address, denom string,
	amount *big.Int, nonce *big.Int, blockNumber *big.Int) *CrossChainTFuelTransferEvent {
	return &CrossChainTFuelTransferEvent{sender, receiver, denom, amount, nonce, blockNumber}
}

func ParseToCrossChainTFuelTransferEvent(icme *InterChainMessageEvent) (*CrossChainTFuelTransferEvent, error) {
	if icme.Type != IMCEventTypeCrossChainTFuelTransfer {
		return nil, fmt.Errorf("invalid inter-chain message event type: %v", icme.Type)
	}

	var tma TfuelTransferMetaData
	if err := rlp.DecodeBytes(icme.Data, &tma); err != nil {
		return nil, err
	}
	if err := ValidateDenom(tma.Denom); err != nil {
		return nil, err
	}
	extractedSourceChainID, err := ExtractSourceChainIDFromDenom(tma.Denom)
	if err != nil {
		return nil, err
	}
	if icme.SourceChainID != extractedSourceChainID {
		return nil, fmt.Errorf("source chain ID mismatch for TFuel transfer: %v vs %v", icme.SourceChainID, extractedSourceChainID)
	}

	ccatEvent := NewCrossChainTFuelTransferEvent(icme.Sender, icme.Receiver, tma.Denom, tma.LockedAmount, icme.Nonce, icme.BlockNumber)
	return ccatEvent, nil
}

func (cct *CrossChainTFuelTransferEvent) IsVoucherBurn(selfChainID string) (bool, error) {
	isVoucherBurn, err := isVoucherBurn(selfChainID, cct.Denom)
	return isVoucherBurn, err
}

// Cross-Chain TNT20 Transfer

type TNT20TransferMetaData struct {
	TargetChainID         *big.Int
	MainchainTokenSender  common.Address
	SubchainTokenReceiver common.Address
	LockedAmount          *big.Int
	TNT20Contract         common.Address
	Name                  string
	Symbol                string
	Decimal               uint8
	Nonce                 *big.Int
	Denom                 string
}

type CrossChainTNT20TransferEvent struct {
	Sender      common.Address
	Receiver    common.Address
	Denom       string
	Name        string
	Symbol      string
	Decimals    uint8
	Amount      *big.Int
	Nonce       *big.Int
	BlockNumber *big.Int
}

func NewCrossChainTNT20TransferEvent(sender common.Address, receiver common.Address, denom string,
	name string, symbol string, decimals uint8, amount *big.Int, nonce *big.Int, blockNumber *big.Int) *CrossChainTNT20TransferEvent {
	return &CrossChainTNT20TransferEvent{sender, receiver, denom, name, symbol, decimals, amount, nonce, blockNumber}
}

func ParseToCrossChainTNT20TransferEvent(icme *InterChainMessageEvent) (*CrossChainTNT20TransferEvent, error) {
	if icme.Type != IMCEventTypeCrossChainTNT20Transfer {
		return nil, fmt.Errorf("invalid inter-chain message event type: %v", icme.Type)
	}

	var tma TNT20TransferMetaData
	if err := rlp.DecodeBytes(icme.Data, &tma); err != nil {
		return nil, err
	}
	if err := ValidateDenom(tma.Denom); err != nil {
		return nil, err
	}
	extractedSourceChainID, err := ExtractSourceChainIDFromDenom(tma.Denom)
	if err != nil {
		return nil, err
	}
	if icme.SourceChainID != extractedSourceChainID {
		return nil, fmt.Errorf("source chain ID mismatch for TNT20 transfer: %v vs %v", icme.SourceChainID, extractedSourceChainID)
	}
	ccatEvent := NewCrossChainTNT20TransferEvent(icme.Sender, icme.Receiver, tma.Denom, tma.Name, tma.Symbol, tma.Decimal, tma.LockedAmount, icme.Nonce, icme.BlockNumber)
	return ccatEvent, nil
}

func (cct *CrossChainTNT20TransferEvent) IsVoucherBurn(selfChainID string) (bool, error) {
	isVoucherBurn, err := isVoucherBurn(selfChainID, cct.Denom)
	return isVoucherBurn, err
}

// Cross-Chain TNT721 Transfer

type TNT721TransferMetaData struct {
	Denom    string
	Name     string
	Symbol   string
	TokenID  *big.Int
	TokenURI string
}

type CrossChainTNT721TransferEvent struct {
	Sender      common.Address
	Receiver    common.Address
	Denom       string
	Name        string
	Symbol      string
	TokenID     *big.Int
	TokenURI    string
	Nonce       *big.Int
	BlockNumber *big.Int
}

func NewCrossChainTNT721TransferEvent(sender common.Address, receiver common.Address, denom string,
	name string, symbol string, tokenID *big.Int, tokenURI string, nonce *big.Int, blockNumber *big.Int) *CrossChainTNT721TransferEvent {
	return &CrossChainTNT721TransferEvent{sender, receiver, denom, name, symbol, tokenID, tokenURI, nonce, blockNumber}
}

func ParseToCrossChainTNT721TransferEvent(icme *InterChainMessageEvent) (*CrossChainTNT721TransferEvent, error) {
	if icme.Type != IMCEventTypeCrossChainTNT721Transfer {
		return nil, fmt.Errorf("invalid inter-chain message event type: %v", icme.Type)
	}

	var tma TNT721TransferMetaData
	if err := rlp.DecodeBytes(icme.Data, &tma); err != nil {
		return nil, err
	}
	if err := ValidateDenom(tma.Denom); err != nil {
		return nil, err
	}
	extractedSourceChainID, err := ExtractSourceChainIDFromDenom(tma.Denom)
	if err != nil {
		return nil, err
	}
	if icme.SourceChainID != extractedSourceChainID {
		return nil, fmt.Errorf("source chain ID mismatch for TNT721 transfer: %v vs %v", icme.SourceChainID, extractedSourceChainID)
	}

	ccatEvent := NewCrossChainTNT721TransferEvent(icme.Sender, icme.Receiver, tma.Denom, tma.Name, tma.Symbol, tma.TokenID, tma.TokenURI, icme.Nonce, icme.BlockNumber)
	return ccatEvent, nil
}

func (cct *CrossChainTNT721TransferEvent) IsVoucherBurn(selfChainID string) (bool, error) {
	isVoucherBurn, err := isVoucherBurn(selfChainID, cct.Denom)
	return isVoucherBurn, err
}

// ------------------------------------ Denom Utils ----------------------------------------------

type CrossChainTokenType int

const (
	CrossChainTokenTypeInvalid CrossChainTokenType = -1
	CrossChainTokenTypeTFuel   CrossChainTokenType = 0
	CrossChainTokenTypeTNT20   CrossChainTokenType = 20
	CrossChainTokenTypeTNT721  CrossChainTokenType = 721
)

const tfuelAddressPlaceholder = "0x0000000000000000000000000000000000000000"

// sourceChainID: the chainID of the chain that the token was originated
func TFuelDenom(sourceChainID string) string {
	return strings.ToLower(fmt.Sprintf("%v/%v/%v", sourceChainID, CrossChainTokenTypeTFuel, tfuelAddressPlaceholder)) // normalize to lower case to prevent duplication
}

func TNT20Denom(sourceChainID string, contractAddress common.Address) string {
	return strings.ToLower(fmt.Sprintf("%v/%v/%v", sourceChainID, CrossChainTokenTypeTNT20, contractAddress.Hex())) // normalize to lower case to prevent duplication
}

func TNT721Denom(sourceChainID string, contractAddress common.Address) string {
	return strings.ToLower(fmt.Sprintf("%v/%v/%v", sourceChainID, CrossChainTokenTypeTNT721, contractAddress.Hex())) // normalize to lower case to prevent duplication
}

func ValidateDenom(denom string) error {
	if !isLowerCase(denom) {
		return fmt.Errorf("invalid denom (must be lower case): %v", denom)
	}

	parts := strings.Split(denom, "/")
	if len(parts) != 3 {
		return fmt.Errorf("invalid denom (incorrect format): %v", denom)
	}

	tokenType, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid denom (failed to parse token type): %v, %v", denom, err)
	}

	switch CrossChainTokenType(tokenType) {
	case CrossChainTokenTypeTFuel:
		if parts[2] != tfuelAddressPlaceholder {
			return fmt.Errorf("invalid TFuel denom: %v", denom)
		}
	case CrossChainTokenTypeTNT20:
		if !common.IsHexAddress(parts[2]) {
			return fmt.Errorf("invalid TNT20 denom: %v", denom)
		}
	case CrossChainTokenTypeTNT721:
		if !common.IsHexAddress(parts[2]) {
			return fmt.Errorf("invalid TNT20 denom: %v", denom)
		}
	default:
		return fmt.Errorf("invalid denom (unknown token type): %v", denom)
	}

	return nil
}

func ExtractSourceChainIDFromDenom(denom string) (string, error) {
	parts := strings.Split(denom, "/")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid denom: %v", denom)
	}

	return parts[0], nil
}

func ExtractCrossChainTokenTypeFromDenom(denom string) (CrossChainTokenType, error) {
	parts := strings.Split(denom, "/")
	if len(parts) != 3 {
		return CrossChainTokenTypeInvalid, fmt.Errorf("invalid denom: %v", denom)
	}
	tokenType, err := strconv.Atoi(parts[1])
	if err != nil {
		return CrossChainTokenTypeInvalid, fmt.Errorf("invalid denom: %v", denom)
	}

	if (tokenType != int(CrossChainTokenTypeTFuel)) && (tokenType != int(CrossChainTokenTypeTNT20)) && (tokenType != int(CrossChainTokenTypeTNT721)) {
		return CrossChainTokenTypeInvalid, fmt.Errorf("invalid denom: %v", denom)
	}

	return CrossChainTokenType(tokenType), nil
}

func isLowerCase(str string) bool {
	return str == strings.ToLower(str)
}

func isVoucherBurn(selfChainID string, denom string) (bool, error) {
	extractedSourceChainID, err := ExtractSourceChainIDFromDenom(denom)
	if err != nil {
		return false, err
	}
	isVoucherBurn := (extractedSourceChainID != selfChainID)
	return isVoucherBurn, nil
}
