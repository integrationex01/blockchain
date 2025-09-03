package block

import (
	"encoding/json"
	"fmt"
)

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(sender, recipient string, value float32) *Transaction {
	return &Transaction{
		senderBlockchainAddress:    sender,
		recipientBlockchainAddress: recipient,
		value:                      value,
	}
}

func (t *Transaction) String() string {
	return fmt.Sprintf("{Sender: %s, Recipient: %s, Value: %.2f}",
		t.senderBlockchainAddress, t.recipientBlockchainAddress, t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func (bc *Blockchain) CalculateTotalAmount(senderAddress string) float32 {
	total := float32(0)
	for _, block := range bc.chain {
		for _, tx := range block.transactions {
			if tx.senderBlockchainAddress == senderAddress {
				total -= tx.value
			}
			if tx.recipientBlockchainAddress == senderAddress {
				total += tx.value
			}
		}
	}
	return total
}
