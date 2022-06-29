package core

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"

	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/rlp"
	scta "github.com/thetatoken/thetasubchain/contracts/accessors"
	"github.com/thetatoken/thetasubchain/eth/abi"
)

// var logger *log.Entry = log.WithFields(log.Fields{"prefix": "core"})

type InterChainMessageEventType uint64

const (
	IMCEventTypeUnknown InterChainMessageEventType = 0
	// 1 - 9999 reserved for future use

	IMCEventTypeCrossChainTokenLock       InterChainMessageEventType = 10000
	IMCEventTypeCrossChainTokenLockTFuel  InterChainMessageEventType = 10001
	IMCEventTypeCrossChainTokenLockTNT20  InterChainMessageEventType = 10002
	IMCEventTypeCrossChainTokenLockTNT721 InterChainMessageEventType = 10003

	IMCEventTypeCrossChainVoucherBurn       InterChainMessageEventType = 20000
	IMCEventTypeCrossChainVoucherBurnTFuel  InterChainMessageEventType = 20001
	IMCEventTypeCrossChainVoucherBurnTNT20  InterChainMessageEventType = 20002
	IMCEventTypeCrossChainVoucherBurnTNT721 InterChainMessageEventType = 20003

	IMCEventTypeCrossChainTokenUnlock       InterChainMessageEventType = 40000
	IMCEventTypeCrossChainTokenUnlockTFuel  InterChainMessageEventType = 40001
	IMCEventTypeCrossChainTokenUnlockTNT20  InterChainMessageEventType = 40002
	IMCEventTypeCrossChainTokenUnlockTNT721 InterChainMessageEventType = 40003
)

// InterChainMessageEvent represents an inter-chain messaging event.
type InterChainMessageEvent struct {
	Type          InterChainMessageEventType
	SourceChainID *big.Int
	TargetChainID *big.Int
	Sender        common.Address // sender of the message on the source chain
	Receiver      common.Address // receiver of the msssage on the target chain
	Data          common.Bytes   // generic data field that can be used to encode arbitrary data for inter-chain messaging
	Nonce         *big.Int
	BlockHeight   *big.Int
}

// NewInterChainMessageEvent creates a new inter-chain messaging event instance.
func NewInterChainMessageEvent(eventType InterChainMessageEventType, sourceChainID *big.Int, targetChainID *big.Int, sender common.Address, receiver common.Address,
	data common.Bytes, nonce *big.Int, blockHeight *big.Int) *InterChainMessageEvent {
	return &InterChainMessageEvent{eventType, sourceChainID, targetChainID, sender, receiver, data, nonce, blockHeight}
}

// ID returns the ID of the inter-chain messaging event.
func (c *InterChainMessageEvent) ID() string {
	return strconv.FormatUint(uint64(c.Type), 10) + "/" + c.Nonce.String()
}

// Equals checks whether an inter-chain messaging event is identical to the other
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
	if c.BlockHeight.Cmp(x.BlockHeight) != 0 {
		return false
	}
	return true
}

// String represents the string representation of the event
func (c *InterChainMessageEvent) String() string {
	return fmt.Sprintf("{ID: %v, Type: %v, SourceChainID: %v, TargetChainID: %v, Sender: %v, Receiver: %v,  Data: %v, Nonce: %v, BlockHeight: %v}",
		c.ID(), c.Type, c.SourceChainID, c.TargetChainID, c.Sender.Hex(), c.Receiver.Hex(), string(c.Data), c.Nonce.String(), c.BlockHeight.String())
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
		c.BlockHeight,
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

	sourceChainID := big.NewInt(0)
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

	blockHeight := big.NewInt(0)
	err = stream.Decode(blockHeight)
	if err != nil {
		return err
	}
	c.BlockHeight = blockHeight

	return stream.ListEnd()
}

// ------------------------------------ Cross-Chain: Token Lock --------------------------------------------

// Cross-Chain TFuel Lock

type CrossChainTFuelTokenLockedEvent struct { // corresponding to the "TFuelTokenLocked" event
	TargetChainID              *big.Int // targetChain: the chain to send the token to (i.e. on which vouchers will be minted)
	Denom                      string
	SourceChainTokenSender     common.Address
	TargetChainVoucherReceiver common.Address
	LockedAmount               *big.Int
	TokenLockNonce             *big.Int
}

func ParseToCrossChainTFuelTokenLockedEvent(icme *InterChainMessageEvent) (*CrossChainTFuelTokenLockedEvent, error) {
	if icme.Type != IMCEventTypeCrossChainTokenLockTFuel {
		return nil, fmt.Errorf("invalid inter-chain message event type: %v", icme.Type)
	}

	var event CrossChainTFuelTokenLockedEvent
	contractAbi, err := abi.JSON(strings.NewReader(string(scta.MainchainTFuelTokenBankABI)))
	if err != nil {
		return nil, err
	}
	contractAbi.UnpackIntoInterface(&event, "TFuelTokenLocked", icme.Data)
	if err := ValidateDenom(event.Denom); err != nil {
		return nil, err
	}
	originatedChainID, err := ExtractOriginatedChainIDFromDenom(event.Denom)
	if err != nil {
		return nil, err
	}
	if icme.SourceChainID.Cmp(originatedChainID) != 0 {
		// Token Lock events can only happen on the chain where the authenic token was deployed. Thus, the "source chain", i.e. where the token is sending from
		// needs to be the same as the "originated chain".
		return nil, fmt.Errorf("source chain ID mismatch for TFuel lock: %v vs %v", icme.SourceChainID, originatedChainID)
	}

	return &event, nil
}

// Cross-Chain TNT20 Lock

type CrossChainTNT20TokenLockedEvent struct { // corresponding to the "TNT20TokenLocked" event
	TargetChainID              *big.Int // targetChain: the chain to send the token to (i.e. on which vouchers will be minted)
	Denom                      string
	SourceChainTokenSender     common.Address
	TargetChainVoucherReceiver common.Address
	LockedAmount               *big.Int
	Name                       string
	Symbol                     string
	Decimal                    uint8
	TokenLockNonce             *big.Int
}

func ParseToCrossChainTNT20TokenLockedEvent(icme *InterChainMessageEvent) (*CrossChainTNT20TokenLockedEvent, error) {
	if icme.Type != IMCEventTypeCrossChainTokenLockTNT20 {
		return nil, fmt.Errorf("invalid inter-chain message event type: %v", icme.Type)
	}

	var event CrossChainTNT20TokenLockedEvent
	contractAbi, err := abi.JSON(strings.NewReader(string(scta.MainchainTNT20TokenBankABI)))
	if err != nil {
		return nil, err
	}
	contractAbi.UnpackIntoInterface(&event, "TNT20TokenLocked", icme.Data)
	event.Denom = strings.ToLower(event.Denom)
	if err := ValidateDenom(event.Denom); err != nil {
		return nil, err
	}
	originatedChainID, err := ExtractOriginatedChainIDFromDenom(event.Denom)
	if err != nil {
		return nil, err
	}
	if icme.SourceChainID.Cmp(originatedChainID) != 0 {
		// Token Lock events can only happen on the chain where the authenic token was deployed. Thus, the "source chain", i.e. where the token is sending from
		// needs to be the same as the "originated chain".
		return nil, fmt.Errorf("source chain ID mismatch for TNT20 lock: %v vs %v", icme.SourceChainID, originatedChainID)
	}
	return &event, nil
}

// Cross-Chain TNT721 Lock

// type CrossChainTNT721TokenLockedEvent struct { // corresponding to the "TNT721TokenLocked" event
// 	TargetChainID              *big.Int // targetChain: the chain to send the token to (i.e. on which vouchers will be minted)
// 	Denom                      string
// 	SourceChainTokenSender     common.Address
// 	TargetChainVoucherReceiver common.Address
// 	TokenID                    *big.Int
// 	TokenURI                   string
// 	Name                       string
// 	Symbol                     string
// 	TokenLockNonce             *big.Int
// }

// func ParseToCrossChainTNT721TokenLockedEvent(icme *InterChainMessageEvent) (*CrossChainTNT721TokenLockedEvent, error) {
// 	if icme.Type != IMCEventTypeCrossChainTokenLockTNT721 {
// 		return nil, fmt.Errorf("invalid inter-chain message event type: %v", icme.Type)
// 	}

// 	var event CrossChainTNT721TokenLockedEvent
// 	contractAbi, err := abi.JSON(strings.NewReader(string(scta.MainchainTNT721TokenBankABI)))
// 	if err != nil {
// 		return nil, err
// 	}
// 	contractAbi.UnpackIntoInterface(&event, "TNT721TokenLocked", icme.Data)
// 	event.Denom = strings.ToLower(event.Denom)
// 	if err := ValidateDenom(event.Denom); err != nil {
// 		return nil, err
// 	}
// 	originatedChainID, err := ExtractOriginatedChainIDFromDenom(event.Denom)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if icme.SourceChainID.Cmp(originatedChainID) != 0 {
// 		// Token Lock events can only happen on the chain where the authenic token was deployed. Thus, the "source chain", i.e. where the token is sending from
// 		// needs to be the same as the "originated chain".
// 		return nil, fmt.Errorf("source chain ID mismatch for TNT721 lock: %v vs %v", icme.SourceChainID, originatedChainID)
// 	}

// 	return &event, nil
// }

// ------------------------------------ Cross-Chain: Token Unlock --------------------------------------------

// Cross-Chain TFuel Unlock

type CrossChainTFuelUnlockedEvent struct { // corresponding to the "TFuelTokenUnlocked" event
	SourceChainID               *big.Int // sourceChain: the chain on which the voucher burn happens
	Denom                       string
	SourceChainTokenSender      common.Address
	UnlockedAmount              *big.Int
	sourceChainVoucherBurnNonce *big.Int
	tokenUnlockNonce            *big.Int
}

func ParseToCrossChainTFuelTokenUnlockedEvent(icme *InterChainMessageEvent) (*CrossChainTFuelUnlockedEvent, error) {
	if icme.Type != IMCEventTypeCrossChainTokenUnlockTFuel {
		return nil, fmt.Errorf("invalid inter-chain message event type: %v", icme.Type)
	}

	var event CrossChainTFuelUnlockedEvent
	contractAbi, err := abi.JSON(strings.NewReader(string(scta.MainchainTFuelTokenBankABI)))
	if err != nil {
		return nil, err
	}
	contractAbi.UnpackIntoInterface(&event, "TFuelTokenUnlocked", icme.Data)
	if err := ValidateDenom(event.Denom); err != nil {
		return nil, err
	}
	originatedChainID, err := ExtractOriginatedChainIDFromDenom(event.Denom)
	if err != nil {
		return nil, err
	}
	if icme.TargetChainID.Cmp(originatedChainID) != 0 {
		return nil, fmt.Errorf("source chain ID mismatch for TFuel unlock: %v vs %v", icme.TargetChainID, originatedChainID)
	}

	return &event, nil
}

// Cross-Chain TNT20 Unlock

type CrossChainTNT20UnlockedEvent struct { // corresponding to the "TNT20TokenUnlocked" event
	SourceChainID              *big.Int // sourceChain: the chain on which the voucher burn happens
	Denom                      string
	SourceChainTokenSender     common.Address
	TargetChainVoucherReceiver common.Address
	LockedAmount               *big.Int
	Name                       string
	Symbol                     string
	Decimal                    uint8
	TokenLockNonce             *big.Int
}

func ParseToCrossChainTNT20TokenUnlockedEvent(icme *InterChainMessageEvent) (*CrossChainTNT20UnlockedEvent, error) {
	if icme.Type != IMCEventTypeCrossChainTokenUnlockTNT20 {
		return nil, fmt.Errorf("invalid inter-chain message event type: %v", icme.Type)
	}

	var event CrossChainTNT20UnlockedEvent
	contractAbi, err := abi.JSON(strings.NewReader(string(scta.MainchainTNT20TokenBankABI)))
	if err != nil {
		return nil, err
	}
	contractAbi.UnpackIntoInterface(&event, "TNT20TokenUnlocked", icme.Data)
	event.Denom = strings.ToLower(event.Denom)
	if err := ValidateDenom(event.Denom); err != nil {
		return nil, err
	}
	originatedChainID, err := ExtractOriginatedChainIDFromDenom(event.Denom)
	if err != nil {
		return nil, err
	}
	if icme.TargetChainID.Cmp(originatedChainID) != 0 {
		return nil, fmt.Errorf("source chain ID mismatch for TNT20 unlock: %v vs %v", icme.TargetChainID, originatedChainID)
	}
	return &event, nil
}

// ------------------------------------ Cross-Chain: Voucher Burn --------------------------------------------

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

type CrossChainTFuelVoucherBurnedEvent struct { // corresponding to the "TFuelVoucherBurned" event
	TargetChainID            *big.Int // targetChain: the chain on which authentic token will be unlocked after the voucher burn
	Denom                    string
	SourceChainVoucherOwner  common.Address
	TargetChainTokenReceiver common.Address
	Amount                   *big.Int
	VoucherBurnNonce         *big.Int
}

func ParseToCrossChainTFuelVoucherBurnedEvent(icme *InterChainMessageEvent) (*CrossChainTFuelVoucherBurnedEvent, error) {
	if icme.Type != IMCEventTypeCrossChainVoucherBurnTFuel {
		return nil, fmt.Errorf("invalid inter-chain message event type: %v", icme.Type)
	}

	var event CrossChainTFuelVoucherBurnedEvent
	contractAbi, err := abi.JSON(strings.NewReader(string(scta.MainchainTFuelTokenBankABI)))
	if err != nil {
		return nil, err
	}
	contractAbi.UnpackIntoInterface(&event, "TFuelVoucherBurned", icme.Data)
	if err := ValidateDenom(event.Denom); err != nil {
		return nil, err
	}
	originatedChainID, err := ExtractOriginatedChainIDFromDenom(event.Denom)
	if err != nil {
		return nil, err
	}
	if icme.TargetChainID.Cmp(originatedChainID) != 0 {
		return nil, fmt.Errorf("source chain ID mismatch for TFuel voucher burn: %v vs %v", icme.TargetChainID, originatedChainID)
	}

	return &event, nil
}

type CrossChainTNT20VoucherBurnedEvent struct { // corresponding to the "TNT20VoucherBurned" event
	SourceChainID *big.Int
	TargetChainID *big.Int
	TxHash        common.Hash
	Denom         string
	Amount        *big.Int
	TokenID       *big.Int
}

func ParseToCrossChainTNT20VoucherBurnedEvent(icme *InterChainMessageEvent) (*CrossChainTNT20VoucherBurnedEvent, error) {
	if icme.Type != IMCEventTypeCrossChainVoucherBurnTNT20 {
		return nil, fmt.Errorf("invalid inter-chain message event type: %v", icme.Type)
	}

	var event CrossChainTNT20VoucherBurnedEvent
	contractAbi, err := abi.JSON(strings.NewReader(string(scta.MainchainTNT20TokenBankABI)))
	if err != nil {
		return nil, err
	}
	contractAbi.UnpackIntoInterface(&event, "TNT20VoucherBurned", icme.Data)
	if err := ValidateDenom(event.Denom); err != nil {
		return nil, err
	}
	originatedChainID, err := ExtractOriginatedChainIDFromDenom(event.Denom)
	if err != nil {
		return nil, err
	}
	if icme.TargetChainID.Cmp(originatedChainID) != 0 {
		return nil, fmt.Errorf("source chain ID mismatch for TNT20 voucher burn: %v vs %v", icme.TargetChainID, originatedChainID)
	}

	return &event, nil
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

// originatedChainID: the chainID of the chain that the token was originated
func TFuelDenom(originatedChainID string) string {
	return strings.ToLower(fmt.Sprintf("%v/%v/%v", originatedChainID, CrossChainTokenTypeTFuel, tfuelAddressPlaceholder)) // normalize to lower case to prevent duplication
}

func TNT20Denom(originatedChainID string, contractAddress common.Address) string {
	return strings.ToLower(fmt.Sprintf("%v/%v/%v", originatedChainID, CrossChainTokenTypeTNT20, contractAddress.Hex())) // normalize to lower case to prevent duplication
}

func TNT721Denom(originatedChainID string, contractAddress common.Address) string {
	return strings.ToLower(fmt.Sprintf("%v/%v/%v", originatedChainID, CrossChainTokenTypeTNT721, contractAddress.Hex())) // normalize to lower case to prevent duplication
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

func ExtractOriginatedChainIDFromDenom(denom string) (*big.Int, error) {
	parts := strings.Split(denom, "/")
	if len(parts) != 3 {
		return big.NewInt(0), fmt.Errorf("invalid denom: %v", denom)
	}

	chainID, ok := big.NewInt(0).SetString(parts[0], 10)
	if !ok {
		return big.NewInt(0), fmt.Errorf("invalid denom: %v", denom)
	}

	return chainID, nil
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
