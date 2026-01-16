package blockchain

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockchainClient struct {
	client *ethclient.Client 
	chainID *big.Int
	privateKey *ecdsa.PrivateKey
	publicAddress common.Address 
	pokerTableAddress common.Address
	potManagerAddress common.Address
	playerRegistryAddress common.Address
	disputeResolverAddres common.Address

	// pokerTable *PokerTable 
	// potManager *PotManager 
	// playerRegistry *PlayerRegistry 
	// disputeResolver *DisputeResolver 
}