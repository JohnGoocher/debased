package main

import (
	"crypto/sha256"
	// "encoding/hex"
	//"fmt"
	"log"
	//"strings"
	"strconv"

	// for creating a buffer for ouputting the blockchain
	"bytes"
	"github.com/davecgh/go-spew/spew"
	// "math/big"
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

// // -----------------------------------------------------------------------------
// // -----------------------------------------------------------------------------
// // -----------------------------------------------------------------------------
// // -----------------------------------------------------------------------------
// // ----------------------------------BLOCK--------------------------------------
// // -----------------------------------------------------------------------------
// // -----------------------------------------------------------------------------
// // -----------------------------------------------------------------------------
// // -----------------------------------------------------------------------------

// The Block type
type Block struct {
	// TODO: Do a slice of bytes
	Transactions  [][]byte
	PublicKey     []byte
	PrevPublicKey []byte
	Index         int
}

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// --------------------------------CELL LOCATION-----------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

// type CellLocation struct {
// 	BlockNumber *big.Int
// 	//position 0 is the first transaction in a block
// 	Position *big.Int
// 	//byte position the cell begins at
// 	PostionInRecord *big.Int
// }

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// --------------------------------BLOCKCHAIN-----------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

// BlockChain is a slice of type *Block
type BlockChain struct {
	Blocks []*Block
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


// If the input slice(list) of blocks has a length that is greater
// than the existent block chain, the existent block chain becomes
// the new block chain.
func (blockChain *BlockChain) replaceChain(newBlocks []*Block) {
	if len(newBlocks) > len(blockChain.Blocks) {
		blockChain.Blocks = newBlocks
	}
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

// TODO: ensure that we can append byte slices together where record is.
func calculateHash(b Block) []byte {
	// buffer := bytes.Buffer
	// buffer.WriteString(strconv.Itoa(b.Index))
	// byte(strconv.Itoa(b.Index))
	// buffer.WriteString(string(bytes.Join(b.Transactions, []byte(""))))
	// bytes.Join(b.Transactions, []byte(""))
	// buffer.WriteString(string(b.PrevPublicKey))
	// b.PrevPublicKey
	h := sha256.New()
	h.Write([]byte(strconv.Itoa(b.Index)))
	h.Write(bytes.Join(b.Transactions, []byte("")))
	h.Write(b.PrevPublicKey)
	hashed := h.Sum(nil)
	return hashed
}

// "createBlock" is a function which creates a new block. This function does
// not add a block to the BlockChain - done by addBlock.
func createBlock(oldBlock *Block, transactions [][]byte) (*Block, error) {

	newBlock := &Block{
		Index:         oldBlock.Index + 1,
		Transactions:  transactions,
		PrevPublicKey: oldBlock.PublicKey,
	}

	newBlock.PublicKey = calculateHash(*newBlock)

	return newBlock, nil
}

// Checks if the block that is created is valid to be placed on the blockChain
// TODO: Ensure that all fields on the block given are instantiated in this function as well?
func isBlockValid(newBlock, oldBlock Block) bool {

	if bytes.Compare(newBlock.PrevPublicKey, oldBlock.PublicKey) != 0 {
		return false
	}

	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if bytes.Compare(calculateHash(newBlock), newBlock.PublicKey) != 0 {
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

	// Build the initial 2d slice
	initialBlockTransactions := [][]byte{}
	firstTransaction := []byte{'a','b', 'c'}
	secondTransaction := []byte{'d','e','f'}

	initialBlockTransactions = append(initialBlockTransactions, firstTransaction)
	initialBlockTransactions = append(initialBlockTransactions, secondTransaction)


	// A block contains the Transactions, PublicKey, PrevPublicKey, and Index fields
	// respectively

	genesisBlock := &Block{
		Transactions:  initialBlockTransactions,
		PublicKey:     []byte{},
		PrevPublicKey: []byte{},
		Index:         0,
	}

	blockChain := &BlockChain{
		Blocks: []*Block{},
	}

	blockChain.addBlock(genesisBlock)


	// Create a second block for testing purposes
	firstBlock, err := blockChain.getLastBlock()
	// firstBlock is now a pointer to the last block in the block chain now

	if err != nil {
		log.Fatal(err)
	}

	// Build the second 2d slice
	secondBlockTransactions := [][]byte{}
	thirdTransaction := []byte{'g','h','i'}
	fourthTransaction := []byte{'j','k','l'}

	secondBlockTransactions = append(secondBlockTransactions, thirdTransaction)
	secondBlockTransactions = append(secondBlockTransactions, fourthTransaction)

	// Generates a new block and puts it at the front of the BlockChain
	newBlock, err := createBlock(firstBlock, secondBlockTransactions)

	if err != nil {
		// fmt.Println("Hey there is an error in generating this new block")
		log.Fatal(err)
	}


	blockChain.addBlock(newBlock)


	// Create a third block for testing purposes

	secondBlock, err := blockChain.getLastBlock()

	thirdBlockTransactions := [][]byte{}
	fifthTransaction := []byte{'m','n','o'}
	sixthTransaction := []byte{'p','q','r'}

	thirdBlockTransactions = append(thirdBlockTransactions, fifthTransaction)
	thirdBlockTransactions = append(thirdBlockTransactions, sixthTransaction)

	newerBlock, err := createBlock(secondBlock, thirdBlockTransactions)

	blockChain.addBlock(newerBlock)

	spew.Dump(blockChain)
}
