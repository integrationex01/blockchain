package server

import (
	"blockchain/block"
	"blockchain/utils"
	"blockchain/wallet"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"
)

const tempDir = "server/templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{
		port:    port,
		gateway: gateway,
	}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "wallet_index.html"))
		t.Execute(w, "")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}

}

func (ws *WalletServer) WalletCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		wallet := wallet.NewWallet()
		m, _ := wallet.MarshalJSON()
		io.Writer.Write(w, m[:])
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Method not allowed"))
	}
}

func (ws *WalletServer) CreateTrasaction(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateTransaction, Method: %s\n", r.Method)
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t wallet.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("CreateTransaction ERR:%s", err.Error())
			io.Writer.Write(w, utils.JsonMessage("Create Transaction Error"))
			return
		}
		if !t.Validate() {
			log.Println("ERR: missing field(s)")
			io.Writer.Write(w, utils.JsonMessage("Missing field(s)"))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Printf("ERR: invalid value (%s)", *t.Value)
			io.Writer.Write(w, utils.JsonMessage("Invalid value"))
			return
		}
		w.Header().Set("Content-Type", "application/json")

		transaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress, float32(value))
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		value32 := float32(value)

		bt := &block.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}
		btBytes, _ := json.Marshal(bt)
		buf := bytes.NewBuffer(btBytes)
		log.Printf("CreateTransaction, Send to Gateway %s\n", ws.Gateway()+"/transaction")
		resp, _ := http.Post(ws.Gateway()+"/transaction", "application/json", buf)
		log.Printf("CreateTransaction, Resp: %v\n", resp)
		if resp.StatusCode == 201 {
			io.Writer.Write(w, utils.JsonMessage("Transaction created successfully"))
			return
		} else {
			respBytes, _ := io.ReadAll(resp.Body)
			io.Writer.Write(w, respBytes)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Method not allowed"))
	}
}

func (ws *WalletServer) Start() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet/create", ws.WalletCreate)
	http.HandleFunc("/transaction", ws.CreateTrasaction)
	log.Printf("Starting Wallet server at port %d\n", ws.Port())
	log.Printf("Wallet server Gateway: %s\n", ws.Gateway())
	log.Printf("Wallet server's filePath %s\n", path.Join(tempDir, "wallet_index.html"))
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
