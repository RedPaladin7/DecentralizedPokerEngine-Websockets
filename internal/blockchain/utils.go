package blockchain

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateGameID(creater common.Address, timestamp int64, buyIn *big.Int) [32]byte {
	data := append(creater.Bytes(), big.NewInt(timestamp).Bytes()...)
	data = append(data, buyIn.Bytes()...)
	return crypto.Keccak256Hash(data)
}

func BytesToGameID(b []byte) ([32]byte, error) {
	var gameID [32]byte 
	if len(b) != 32 {
		return gameID, fmt.Errorf("invalid game ID length: expected 32, got %d", len(b))
	}
	copy(gameID[:], b)
	return gameID, nil
}

func HexToGameID(hexStr string)([32]byte, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	b, err := hex.DecodeString(hexStr) // turns hex string to byte value
	if err != nil {
		return [32]byte{}, fmt.Errorf("invalid hex string: %w", err)
	}
	return BytesToGameID(b)
}

func GameIDToHex(gameID [32]byte) string {
	return "0x" + hex.EncodeToString(gameID[:])
}

func IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

func FormatAddress(address common.Address) string {
	return address.Hex()
}

func ParseAddress(addressStr string) (common.Address, error) {
	if !IsValidAddress(addressStr) {
		return common.Address{}, fmt.Errorf("invalid address: %s", addressStr)
	}
	return common.HexToAddress(addressStr), nil 
	// converts address in string format to ethereum address
	// type: common.Address
}

func ConvertToWei(amount float64) *big.Int {
	return EthToWei(big.NewFloat(amount))
}

func ConvertFromWei(wei *big.Int) float64 {
	eth := WeiToEth(wei)
	result, _ := eth.Float64()
	return result
}

func FormatWei(wei *big.Int) string {
	eth := WeiToEth(wei)
	return fmt.Sprintf("%.6f ETH", eth)
}

func CalculatePlatformFee(pot *big.Int, feePercent int) *big.Int {
	fee := new(big.Int).Mul(pot, big.NewInt(int64(feePercent)))
	return new(big.Int).Div(fee, big.NewInt(100))
}

func CalculateNetPot(pot *big.Int, feePercent int) *big.Int {
	fee := CalculatePlatformFee(pot, feePercent)
	return new(big.Int).Sub(pot, fee)
}

func SplitPot(pot *big.Int, numWinners int) []*big.Int {
	if numWinners <= 0 {
		return []*big.Int{}
	}

	share := new(big.Int).Div(pot, big.NewInt(int64(numWinners)))
	remainder := new(big.Int).Mod(pot, big.NewInt(int64(numWinners)))

	shares := make([]*big.Int, numWinners)
	for i := 0; i < numWinners; i++ {
		shares[i] = new(big.Int).Set(share)
		if i == 0 {
			shares[i].Add(shares[i], remainder)
		}
	}

	return shares
}

func ValidateBuyIn(buyIn *big.Int, minBuyIn, maxBuyIn *big.Int) error {
	if buyIn.Cmp(minBuyIn) < 0 {
		return fmt.Errorf("buy-in %s is less than minumum %s", FormatWei(buyIn), FormatWei(minBuyIn))
	}
	if buyIn.Cmp(maxBuyIn) > 0 {
		return fmt.Errorf("buy-in %s is greater than maximum %s", FormatWei(buyIn), FormatWei(maxBuyIn))
	}
	return nil
}

func CalculateGasCost(gasUsed uint64, gasPrice *big.Int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(gasUsed)), gasPrice)
}

func FormatGasCost(gasUsed uint64, gasPrice *big.Int) string {
	cost := CalculateGasCost(gasUsed, gasPrice)
	return FormatWei(cost)
}

func HashMessage(message []byte) common.Hash {
	return crypto.Keccak256Hash(message)
}

func AddressToString(addresses []common.Address) []string {
	result := make([]string, len(addresses))
	for i, addr := range addresses {
		result[i] = addr.Hex()
	}
	return result
}

func StringToAddress(addressStrings []string) ([]common.Address, error) {
	result := make([]common.Address, len(addressStrings))
	for i, addrStr := range addressStrings {
		if !IsValidAddress(addrStr){
			return nil, fmt.Errorf("invalid address at index %d: %s", i, addrStr)
		}
		result[i] = common.HexToAddress(addrStr)
	}
	return result, nil
}