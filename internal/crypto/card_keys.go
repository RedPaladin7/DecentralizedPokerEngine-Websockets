package crypto

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

type SerializedKeys struct {
	EncKey string `json:"enc_key"`
	DecKey string `json:"dec_key"`
	Prime  string `json:"prime"`
}

func (ck *CardKeys) Serialize() SerializedKeys {
	// .text function has configurable base
	// .string is base 10 by default
	return SerializedKeys{
		EncKey: ck.EncKey.Text(16),
		DecKey: ck.DecKey.Text(16),
		Prime:  ck.Prime.Text(16),
	}
}

func DeserializeKeys(sk SerializedKeys) (*CardKeys, error) {
	encKey := new(big.Int)
	if _, ok := encKey.SetString(sk.EncKey, 16); !ok {
		return nil, fmt.Errorf("invalid encryption format")
	}
	decKey := new(big.Int)
	if _, ok := decKey.SetString(sk.DecKey, 16); !ok {
		return nil, fmt.Errorf("invalid decryption format")
	}
	prime := new(big.Int)
	if _, ok := prime.SetString(sk.Prime, 16); !ok {
		return nil, fmt.Errorf("invalid prime format")
	}
	return &CardKeys{
		EncKey: encKey,
		DecKey: decKey,
		Prime:  prime,
	}, nil
}

func ToHex(data []byte) string {
	return hex.EncodeToString(data)
}

func FromHex(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

func EncryptDeck(deck [][]byte, keys *CardKeys) [][]byte {
	encrypted := make([][]byte, len(deck))
	for i, card := range deck {
		encrypted[i] = keys.Encrypt(card)
	}
	return encrypted
}

func DecryptDeck(deck [][]byte, keys *CardKeys) [][]byte {
	decrypted := make([][]byte, len(deck))
	for i, card := range deck {
		decrypted[i] = keys.Decrypt(card)
	}
	return decrypted
}

func DecryptSpeceficCards(deck [][]byte, indices []int, keys *CardKeys) map[int][]byte {
	decrypted := make(map[int][]byte)
	for _, idx := range indices {
		decrypted[idx] = keys.Decrypt(deck[idx])
	}
	return decrypted
}

func VerifyDecryption(original []byte, keys *CardKeys) bool {
	encrypted := keys.Encrypt(original)
	decrypted := keys.Decrypt(encrypted)
	if len(original) != len(decrypted) {
		return false 
	}
	for i := range original {
		if original[i] != decrypted[i] {
			return false
		}
	}
	return true 
}

func CombineDecryption(data []byte, keyList []*CardKeys) []byte {
	result := data 
	for _, keys := range keyList {
		result = keys.Decrypt(result)
	}
	return result
}

func CombineEncryption(data []byte, keyList []*CardKeys) []byte {
	result := data 
	for _, keys := range keyList {
		result = keys.Encrypt(result)
	}
	return result
}