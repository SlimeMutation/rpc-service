package client

import (
	"fmt"
	"testing"
)

func TestGetSupportCoins(t *testing.T) {
	client := NewWalletClient("http://127.0.0.1:8970")
	support, err := client.GetSupportCoins("Ethereum", "MainNet")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("support: ", support)
}

func TestGetWalletAddress(t *testing.T) {
	client := NewWalletClient("http://127.0.0.1:8970")
	address, err := client.GetWalletAddress("Ethereum", "MainNet")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("address: ", address)
}
