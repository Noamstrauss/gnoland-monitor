package main

import "encoding/json"

//nolint:govet // Using existing API field names
type Transaction struct {
	Hash        string               `json:"hash"`
	BlockHeight int                  `json:"block_height"`
	Success     bool                 `json:"success"`
	Messages    []TransactionMessage `json:"messages"`
}

//nolint:revive // Using existing API field names
type TransactionMessage struct {
	TypeUrl string          `json:"typeUrl"`
	Route   string          `json:"route"`
	Value   json.RawMessage `json:"value"`
}

type MsgAddPackage struct {
	Creator string     `json:"creator"`
	Package MemPackage `json:"package"`
}

type MemPackage struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type BankMsgSend struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
}
