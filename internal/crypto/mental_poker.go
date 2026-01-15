package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type CardKeys struct {
	EncKey *big.Int
	DecKey *big.Int
	Prime  *big.Int
}

func GenerateCardKeys() (*CardKeys, error) {
	sharedPrim, success := new(big.Int).SetString("C7970CEDCC5226685694605929849D3D", 16)
	if !success {
		return nil, fmt.Errorf("failed to set shared prime")
	}
	return GenerateCardKeysWithPrime(sharedPrim)
}

func GenerateCardKeysWithPrime(prime *big.Int) (*CardKeys, error) {
	encKey, err := generateRandomKey(prime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}
	phiN := new(big.Int).Sub(prime, big.NewInt(1))
	decKey := new(big.Int).ModInverse(encKey, phiN)
	if decKey == nil {
		return nil, fmt.Errorf("failed to calculate decryption key")
	}
	return &CardKeys{
		EncKey: encKey,
		DecKey: decKey,
		Prime:  prime,
	}, nil
}

func generateRandomKey(prime *big.Int) (*big.Int, error) {
	phiN := new(big.Int).Sub(prime, big.NewInt(1))
	maxAttemps := 1000 

	for i := 0; i < maxAttemps; i++ {
		key, err := rand.Int(rand.Reader, new(big.Int).Sub(prime, big.NewInt(2)))
		if err != nil {
			return nil, err
		}
		key.Add(key, big.NewInt(2))
		gcd := new(big.Int).GCD(nil, nil, key, phiN)
		if gcd.Cmp(big.NewInt(1)) == 0 {
			return key, nil
		}
	}
	return nil, fmt.Errorf("failed to generate comprime key")
}

func (ck *CardKeys) Encrypt(data []byte) []byte {
	// convert bytes array to single large integer
	plaintext := new(big.Int).SetBytes(data) 
	// plaintext ^ encKey mod prime 
	ciphertext := new(big.Int).Exp(plaintext, ck.EncKey, ck.Prime)
	return ciphertext.Bytes()
}

func (ck *CardKeys) Decrypt(data []byte) []byte {
	ciphertext := new(big.Int).SetBytes(data)
	plaintext := new(big.Int).Exp(ciphertext, ck.DecKey, ck.Prime)
	return plaintext.Bytes()
}

func (ck *CardKeys) EncryptMultiple(dataList [][]byte) [][]byte {
	encrypted := make([][]byte, len(dataList))
	for i, data := range dataList {
		encrypted[i] = ck.Encrypt(data)
	}
	return encrypted
}

func (ck *CardKeys) DecryptMultiple(dataList [][]byte) [][]byte {
	decrypted := make([][]byte, len(dataList))
	for i, data := range dataList {
		decrypted[i] = ck.Decrypt(data)
	}
	return decrypted
}

func (ck *CardKeys) Clone() *CardKeys {
	return &CardKeys{
		EncKey: new(big.Int).Set(ck.EncKey),
		DecKey: new(big.Int).Set(ck.DecKey),
		Prime:  new(big.Int).Set(ck.Prime),
	}
}

func (ck *CardKeys) Validate() error {
	if ck.EncKey == nil || ck.DecKey == nil || ck.Prime == nil {
		return fmt.Errorf("keys cannot be nil")
	}
	phiN := new(big.Int).Sub(ck.Prime, big.NewInt(1))
	product := new(big.Int).Mul(ck.EncKey, ck.DecKey)
	modResult := new(big.Int).Mod(product, phiN)
	if modResult.Cmp(big.NewInt(1)) != 0 {
		return fmt.Errorf("invalid key pair")
	}
	return nil
}