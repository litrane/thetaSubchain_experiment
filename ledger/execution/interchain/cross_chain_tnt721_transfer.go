package interchain

import (
	"fmt"

	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/ledger/types"
	"github.com/thetatoken/thetasubchain/contracts/predeployed"
	score "github.com/thetatoken/thetasubchain/core"
	slst "github.com/thetatoken/thetasubchain/ledger/state"
)

func ConstructMintTNT721VoucherProxySctx(blockProposer common.Address, view *slst.StoreView, icme *score.InterChainMessageEvent) (*types.SmartContractTx, error) {
	ccte, err := score.ParseToCrossChainTNT721TransferEvent(icme)
	if err != nil {
		return nil, err
	}

	err = score.ValidateDenom(ccte.Denom)
	if err != nil {
		return nil, err
	}

	tokenType, err := score.ExtractCrossChainTokenTypeFromDenom(ccte.Denom)
	if err != nil {
		return nil, fmt.Errorf("failed to get token type from denom %v: %v", ccte.Denom, err)
	}
	if tokenType != score.CrossChainTokenTypeTNT721 {
		return nil, fmt.Errorf("denom %v is not a TNT721 token", ccte.Denom)
	}

	tokenBank := predeployed.NewTNT721TokenBank()
	proxySctx, err := tokenBank.GenerateMintVouchersProxySctx(blockProposer, view, ccte)
	if err != nil {
		return nil, err
	}

	return proxySctx, nil
}