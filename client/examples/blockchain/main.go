package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	// for creating a buffer for ouputting the blockchain
	"bytes"

	// for finding the type of an object for testing purposes
	// "reflect"

	// Example of reflect:
	//   tst := "string"
  //   tst2 := 10
  //   tst3 := 1.2
	//
  //   fmt.Println(reflect.TypeOf(tst))
  //   fmt.Println(reflect.TypeOf(tst2))
  //   fmt.Println(reflect.TypeOf(tst3))


	// TODO: importing other files to work seamlessly below
	// "github.com/"
)

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// ----------------------------------BLOCK--------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

// The Block type
type Block struct {
	Transactions  []string
	PublicKey     string
	PrevPublicKey string
	Index         int
}

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// --------------------------------BLOCKCHAIN-----------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------


// BlockChain is a slice of type *Block
type BlockChain struct {
	Blocks          []*Block
}

// "addBlock" Adds a block to the front/head of the Blockchain
func (blockChain *BlockChain) addBlock(block *Block) {
	// The parameter is a pointer because we are storing a list of pointers
	blockChain.Blocks = append(blockChain.Blocks, block)
}

// TODO: FIX THIS LATER. DON'T KNOW HOW TO PROPERLY HANDLE IF blockchain IS EMPTY
// getLastBlock returns a pointer to the last block in the blockchain
func (blockChain *BlockChain) getLastBlock() (*Block, error) {
	return blockChain.Blocks[len(blockChain.Blocks)-1], nil
}

// TODO: ToString function:
func (blockChain BlockChain) String() string{
	// blockChain.Blocks is of type []*Block

  var buffer bytes.Buffer

buffer.WriteString("{\n")

  for _, block := range blockChain.Blocks{
		buffer.WriteString("[ ")
		for _, transaction := range block.Transactions{
			buffer.WriteString(transaction)
			buffer.WriteString(", ")
		}
		buffer.WriteString(" ]\n")
  }
	
	buffer.WriteString("}")
  return buffer.String()
}


// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------GENERAL FUNCTIONS-------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

// Calculates the Hash for the block inserted
// This function assumes that all fields of the block have been instantiated
func calculateHash(b Block) string {
	record := string(b.Index) + strings.Join(b.Transactions, ",") + b.PrevPublicKey
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// "createBlock" is a function which creates a new block. This function does
// not add a block to the BlockChain - done by addBlock.
func createBlock(oldBlock *Block, transactions []string) (*Block, error){

	newBlock := &Block{
		Index:        oldBlock.Index + 1,
		Transactions: transactions,
		PrevPublicKey: oldBlock.PublicKey,
	}

	newBlock.PublicKey = calculateHash(*newBlock)

	return newBlock, nil
}

// Checks if the block that is created is valid to be placed on the blockChain
// TODO: Ensure that all fields on the block given are instantiated in this function as well?
func isBlockValid(newBlock, oldBlock Block) bool {
	if newBlock.PrevPublicKey != oldBlock.PublicKey {
		return false
	}

	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if calculateHash(newBlock) != newBlock.PublicKey {
		return false
	}

	return true
}

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// ------------------------------MAIN FUNCTION----------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

func main() {
	initialTransactions := []string{"hey", "hi", "how are you?"}
	secondaryTransactions := []string{"yo", "yo"}

	// A block contains the Transactions, PublicKey, PrevPublicKey, and Index fields
	// respectively

	genesisBlock := &Block{
		Transactions:  initialTransactions,
		PublicKey:     "",
		PrevPublicKey: "",
		Index:         0,
	}

	blockChain := &BlockChain{
		Blocks:          []*Block{},
	}

	blockChain.addBlock(genesisBlock)

	recentBlock, err := blockChain.getLastBlock()
	// recentBlock is now a pointer to the last block in the block chain now

	if err != nil {
		log.Fatal(err)
	}

	// Generates a new block and puts it at the front of the BlockChain
	newBlock, err := createBlock(recentBlock, secondaryTransactions)

	if err != nil {
		// fmt.Println("Hey there is an error in generating this new block")
		log.Fatal(err)
	}

	blockChain.addBlock(newBlock)

	// TODO: Create a toString method for outputting the blockchain for now.
	fmt.Println("result is:\n", blockChain)
}
