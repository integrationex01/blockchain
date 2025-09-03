package server

import (
	"blockchain/block"
	"blockchain/utils"
	"blockchain/wallet"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{
		port: port,
	}
}

func (s *BlockchainServer) GetPort() uint16 {
	return s.port
}

func (s *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, exists := cache["blockchain"]
	if !exists {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minersWallet.GetBlockchainAddress(), s.GetPort())
		cache["blockchain"] = bc
		log.Printf("New blockchain created with address %s\n", minersWallet.GetBlockchainAddress())
		log.Printf("Private key %s\n", minersWallet.GetPrivateKey())
		log.Printf("Public key %s\n", minersWallet.GetPublicKey())
	}
	return bc
}

func Helloworld(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, BC World!")
}

func (s *BlockchainServer) GetChain(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		bc := s.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.Writer.Write(w, m[:])
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Method not allowed")
	}
}

func (s *BlockchainServer) Transaction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		bc := s.GetBlockchain()
		transaction := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transactions []*block.Transaction `json:"transactions"`
			Length       int                  `json:"length"`
		}{
			Transactions: transaction,
			Length:       len(transaction),
		})
		io.Writer.Write(w, m[:])
		log.Printf("Transaction, Pool: %v\n", transaction)

	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t block.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("Transaction, Error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Bad request: decoding JSON failed")
			return
		}
		if !t.Validate() {
			log.Printf("Transaction, Invalid transaction: %v\n", t)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Bad request: invalid transaction")
			return
		}
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)

		bc := s.GetBlockchain()
		isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Set("Content-Type", "application/json")
		var m []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonMessage("Creating a new transaction failed")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = utils.JsonMessage("Transaction created successfully")
		}
		io.Writer.Write(w, m[:])
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func (s *BlockchainServer) Start() {
	// http.HandleFunc("/", Helloworld)
	http.HandleFunc("/", s.GetChain)
	http.HandleFunc("/transaction", s.Transaction)
	log.Printf("Starting BC server at port %d\n", s.GetPort())
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(s.port)), nil)
}
