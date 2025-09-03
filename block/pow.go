package block

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 1
	MINE_OWNER        = "THE B.C."
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 20
)

func (bc *Blockchain) VaildProof(nonce int, previousHash string, transactions []*Transaction, difficulty int) bool {
	if len(bc.chain) == 0 {
		return false
	}
	block := NewBlock(nonce, previousHash, 0, bc.transactionPool)
	hash := block.Hash()
	prefix := strings.Repeat("0", difficulty)
	if string(hash[:difficulty]) == prefix {
		fmt.Printf("Hash: %x\n", hash)
		fmt.Printf("len: %d\n", len(hash[:difficulty]))
		fmt.Printf("Hash1: %s\n", ByteToString(hash))
		fmt.Printf("Hash2: %s\n", string(hash[:difficulty]))
		fmt.Printf("Hash3: %s\n", ByteToString(hash[:difficulty]))
		fmt.Printf("Prefix: %s\n", prefix)
	}
	return string(hash[:difficulty]) == prefix

}

func (bc *Blockchain) ProofOfWork() int {
	trans := bc.CopyTransactionPool()
	previousHash := ByteToString(bc.LastBlock().Hash())
	nonce := 0
	for !bc.VaildProof(nonce, previousHash, trans, MINING_DIFFICULTY) {
		nonce++
	}
	fmt.Printf("Proof of Work found: %d\n", nonce)
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	if len(bc.transactionPool) == 0 {
		return false
	}
	bc.AddTransaction(MINE_OWNER, bc.blockchianAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := ByteToString(bc.LastBlock().Hash())
	bc.CreateBlock(nonce, previousHash)
	log.Println("action= Mining, status= success, blockchianAddress=", bc.blockchianAddress)
	return true
}

func (bc *Blockchain) StratMining() {
	bc.Mining()
	_ = time.AfterFunc(MINING_TIMER_SEC*time.Second, bc.StratMining)

}
