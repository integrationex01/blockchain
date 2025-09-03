package main

import (
	"blockchain/server"
	"flag"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 5000, "tcp Port number for the BC server")
	trigger := flag.String("server", "bc", "start the blockchain server")
	flag.Parse()

	switch *trigger {
	case "bc":
		app := server.NewBlockchainServer(uint16(*port))
		app.Start()
		fmt.Printf("Starting the BC server, Port: %v\n", *port)
	case "wallet":
		walletApp := server.NewWalletServer(uint16(*port), "http://127.0.0.1:5000")
		walletApp.Start()
		fmt.Printf("Starting the Wallet server, Port: %v\n", *port)
	default:
		log.Fatal("Unknown server type. Use 'bc' or 'wallet'")
	}

}
