package blockchain

import "math/big"

func WeiToEth(wei *big.Int) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		big.NewFloat(1e18),
	)
}

// EthToWei converts eth to wei
func EthToWei(eth *big.Float) *big.Int {
	truncInt, _ := new(big.Float).Mul(
		eth,
		big.NewFloat(1e18),
	).Int(nil)
	return truncInt
}