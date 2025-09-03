package block

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tjfoc/gmsm/sm3"
)

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash string, timestamp int64, transactions []*Transaction) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    timestamp,
		transactions: transactions,
	}
}

func (b *Block) String() string {
	transactionsStr := ""
	for _, tx := range b.transactions {
		transactionsStr += tx.String() + "\n"
	}
	return fmt.Sprintf("{\nNonce: %d,\nPreviousHash: %s,\nTimestamp: %d,\nTransactions: \n",
		b.nonce, b.previousHash, b.timestamp) + transactionsStr + "}"
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash string) *Block {
	timestamp := time.Now().UnixNano()
	b := NewBlock(nonce, previousHash, timestamp, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{} // Clear the transaction pool after creating a block
	return b
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

func (b *Block) Hash() []byte {
	blockJson, _ := json.Marshal(b)
	h := sm3.New()
	h.Write(blockJson)
	return h.Sum(nil)
}

func ByteToString(b []byte) string {
	ret := ""
	for _, value := range b {
		ret += fmt.Sprintf("%02x", value)
	}
	return ret
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}
	return true
}
