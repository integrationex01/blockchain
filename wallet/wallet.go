package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {

	// 1. Create ECDSA private key (32 bytes) public key (64 bytes)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	// 2. Perform SHA256 hashing on the public key (32 bytes)
	h2 := sha256.New()
	h2.Write(privateKey.PublicKey.X.Bytes())
	h2.Write(privateKey.PublicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// 3. Perform RIPEMD160 hashing on the SHA256 hash (20 bytes)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	// 4. Add version byte in front of the RIPEMD160 hash (20 bytes + 1 byte = 21 bytes)
	vd4 := append([]byte{0x00}, digest3...)
	// 5. Perform SHA256 hashing on the versioned hash (32 bytes)
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	// 6. Perform SHA256 hashing on the previous SHA256 hash (32 bytes)
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7. Take the first 4 bytes of the second SHA256 hash for checkSum (4 bytes)
	checksum := digest6[:4]
	// 8. Add the checksum to the end of the versioned hash (21 bytes + 4 bytes = 25 bytes)
	dc8 := append(vd4, checksum...)
	// 9. Convert the address to a Base58 string (not implemented here, but would be done in practice)
	address := base58.Encode(dc8)
	return &Wallet{
		privateKey:        privateKey,
		publicKey:         privateKey.Public().(*ecdsa.PublicKey),
		blockchainAddress: address,
	}
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey:        w.PrivateKeyToString(),
		PublicKey:         w.PublicKeyToString(),
		BlockchainAddress: w.blockchainAddress,
	})
}

func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) GetPublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PrivateKeyToString() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKeyToString() string {
	return fmt.Sprintf("%064x%064x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) GetBlockchainAddress() string {
	return w.blockchainAddress
}

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	Value                      *string `json:"value"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPrivateKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.Value == nil {
		return false
	}
	return true
}

func (tr *TransactionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderPrivateKey           string `json:"sender_private_key"`
		SenderBlockchainAddress    string `json:"sender_blockchain_address"`
		SenderPublicKey            string `json:"sender_public_key"`
		RecipientBlockchainAddress string `json:"recipient_blockchain_address"`
		Value                      string `json:"value"`
	}{
		SenderPrivateKey:           *tr.SenderPrivateKey,
		SenderBlockchainAddress:    *tr.SenderBlockchainAddress,
		SenderPublicKey:            *tr.SenderPublicKey,
		RecipientBlockchainAddress: *tr.RecipientBlockchainAddress,
		Value:                      *tr.Value,
	})
}
