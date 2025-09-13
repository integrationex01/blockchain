package block

import (
	"blockchain/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	BLOCKCHAIN_PORT_RANGE_START      = 5000
	BLOCKCHAIN_PORT_RANGE_END        = 5003
	NEIBOUR_IP_RANGE_START           = 0
	NEIBOUR_IP_RANGE_END             = 0
	BLOCKCHIAN_NEIBOUR_SYNC_INTERVAL = 10 // seconds
)

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchianAddress string
	port              uint16
	mux               sync.Mutex

	neibours    []string
	muxNeibours sync.Mutex
}

func NewBlockchain(blockchianAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, ByteToString(b.Hash()))
	bc.blockchianAddress = blockchianAddress
	bc.port = port
	return bc
}

func (bc *Blockchain) SetNeibours() {
	bc.neibours = utils.FindNeibours(utils.GetHost(), bc.port, BLOCKCHAIN_PORT_RANGE_START, BLOCKCHAIN_PORT_RANGE_END, NEIBOUR_IP_RANGE_START, NEIBOUR_IP_RANGE_END)
	log.Printf("action= SetNeibours, neibours= %v\n", bc.neibours)
}

func (bc *Blockchain) SyncNeigbours() {
	bc.muxNeibours.Lock()
	defer bc.muxNeibours.Unlock()
	bc.SetNeibours()
}

func (bc *Blockchain) StratSyncNeibours() {
	bc.SyncNeigbours()
	_ = time.AfterFunc(BLOCKCHIAN_NEIBOUR_SYNC_INTERVAL*time.Second, func() { bc.StratSyncNeibours() })
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) String() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n%s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25), block.String())
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))

}

func (bc *Blockchain) LastBlock() *Block {
	if len(bc.chain) == 0 {
		return nil
	}
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) CreateTransaction(sender, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isCreated := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)

	// TODO
	// sync

	return isCreated
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINE_OWNER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransaction(senderPublicKey, s, t) {
		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("[ERROR] Not enough balance for transaction")
		// 	return false
		// }
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Fatalln("[ERROR] Invalid transaction signature")
		return false
	}

}

func (bc *Blockchain) VerifyTransaction(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := t.MarshalJSON()
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	copiedPool := make([]*Transaction, len(bc.transactionPool))
	copy(copiedPool, bc.transactionPool)
	return copiedPool
}

type AmountResponse struct {
	Amount *float32 `json:"amount"`
}

func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{
		Amount: *ar.Amount,
	})
}
