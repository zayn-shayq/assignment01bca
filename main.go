package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"github.com/olekukonko/tablewriter"
)

// Block structure containing information about a transaction.
type Block struct {
	Transaction  string // The details of the transaction.
	Nonce        int    // A unique number for the block.
	PreviousHash string // The hash of the previous block in the chain.
	Hash         string // The hash of the current block.
}

// A slice to store the blocks of the blockchain.
var blockchain []*Block
var mutex = &sync.Mutex{} // Mutex to handle concurrent access to the blockchain slice.

// NewBlock creates a new block with the given transaction details, nonce, and previous hash.
// It calculates the hash of the new block and adds it to the blockchain slice.
func NewBlock(transaction string, nonce int, previousHash string) *Block {
	block := &Block{transaction, nonce, previousHash, ""}
	block.Hash = CalculateHash(block)
	mutex.Lock()
	blockchain = append(blockchain, block)
	mutex.Unlock()
	DisplayBlocks()
	return block
}

// DisplayBlocks prints all the blocks in the blockchain in a tabular format.
func DisplayBlocks() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Transaction", "Nonce", "Previous Hash", "Hash"})

	for _, block := range blockchain {
		row := []string{block.Transaction, fmt.Sprintf("%d", block.Nonce), block.PreviousHash, block.Hash}
		table.Append(row)
	}

	table.Render() // Render the table to the standard output.
}

// ChangeBlock updates the transaction of a given block and recalculates its hash.
func ChangeBlock(index int, transaction string) {
	if index < 0 || index >= len(blockchain) {
		fmt.Println("Invalid block index")
		return
	}
	blockchain[index].Transaction = transaction
	blockchain[index].Hash = CalculateHash(blockchain[index])
	fmt.Println("Block changed:")
	DisplayBlocks()
}

// VerifyChain checks the integrity of the blockchain.
// It verifies if each block (except the first one) points to the previous block by comparing hashes.
func VerifyChain() {
	for i := 1; i < len(blockchain); i++ {
		if blockchain[i].PreviousHash != blockchain[i-1].Hash {
			fmt.Println("Blockchain is not valid!")
			return
		}
	}
	fmt.Println("Blockchain is valid.")
}

// CalculateHash generates the hash of a block using SHA-256.
func CalculateHash(block *Block) string {
	record := block.Transaction + string(block.Nonce) + block.PreviousHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// The main function demonstrates the creation, modification, and verification of a blockchain.
func main() {
	// Creating blocks with transactions.
	NewBlock("Ahmed to Fatima: 50 PKR", 1, "0")
	NewBlock("Fatima to Zain: 20 PKR", 2, blockchain[len(blockchain)-1].Hash)
	NewBlock("Zain to Sara: 30 PKR", 3, blockchain[len(blockchain)-1].Hash)

	// Changing a block's transaction.
	ChangeBlock(1, "Fatima to Ali: 25 PKR")

	// Verifying the integrity of the blockchain.
	VerifyChain()
}
