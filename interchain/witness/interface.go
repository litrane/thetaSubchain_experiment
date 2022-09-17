package witness

import (
	"context"
	"math/big"

	score "github.com/thetatoken/thetasubchain/core"
	siu "github.com/thetatoken/thetasubchain/interchain/utils"
)

type ChainWitness interface {
	Start(ctx context.Context)
	Stop()
	Wait()
	GetMainchainBlockHeight() (*big.Int, error)
	GetValidatorSetByDynasty(dynasty *big.Int) (*score.ValidatorSet, error)
	GetValidatorSetByDynastyForChain(dynasty *big.Int, subchainID *big.Int) (*score.ValidatorSet, error)
	GetInterChainEventCache() *siu.InterChainEventCache
	GetInterSubchainChannelWatchList() []*big.Int
}
