package main

import (
	"bytes"
	"crypto/sha256"
	"log"
	"strconv"

	// "fmt"
	"github.com/davecgh/go-spew/spew"
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

func (block *Block) findLineNumber(transaction []byte) int {
	for index, val := range block.Transactions {
		if bytes.Compare(val, transaction) == 0 {
			return index
		}
	}
	return -1
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
	Blocks []*Block
}

func (blockChain *BlockChain) addGenesisBlock(initialBlockTransactions [][]byte) {
	genesisBlock := &Block{
		Transactions:  initialBlockTransactions,
		PublicKey:     []byte{},
		PrevPublicKey: []byte{},
		Index:         0,
	}
	blockChain.Blocks = append(blockChain.Blocks, genesisBlock)
}

// "addBlock" Adds a block to the front/head of the Blockchain
func (blockChain *BlockChain) addBlock(newBlock *Block) {
	// The parameter is a pointer because we are storing a list of pointers
	latestBlock := blockChain.getLatestBlock()

	if isBlockValid(newBlock, latestBlock) {
		newBlockchain := append(blockChain.Blocks, newBlock)
		blockChain.replaceChain(newBlockchain)
	}
}

// getLatestBlock returns a pointer to the last block in the blockchain
func (blockChain *BlockChain) getLatestBlock() *Block {
	return blockChain.Blocks[len(blockChain.Blocks)-1]
}

// If the input slice(list) of blocks has a length that is greater
// than the existent block chain, the existent block chain becomes
// the new block chain.
func (blockChain *BlockChain) replaceChain(newBlocks []*Block) {
	if len(newBlocks) > len(blockChain.Blocks) {
		blockChain.Blocks = newBlocks
	}
}

// The 'write' function writes data on the latest block in the blockchain
func (blockChain *BlockChain) write(data []byte) (BN int, LN int, byteAdd int) {
	blockToWriteTo := blockChain.getLatestBlock()

	blockToWriteTo.Transactions = append(blockToWriteTo.Transactions, data)

	BN = blockToWriteTo.Index
	LN = blockToWriteTo.findLineNumber(data)
	byteAdd = 0

	return
}

// The 'read' funtion reads the data stored within the block number and the line
// number in the block chain
func (blockChain *BlockChain) read(BN int, LN int, byteAdd int) (data []byte) {
	data = blockChain.Blocks[BN].Transactions[LN]
	return
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

func calculateHash(b Block) []byte {
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
func isBlockValid(newBlock *Block, oldBlock *Block) bool {

	if bytes.Compare(newBlock.PrevPublicKey, oldBlock.PublicKey) != 0 {
		return false
	}

	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if bytes.Compare(calculateHash(*newBlock), newBlock.PublicKey) != 0 {
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
	firstTransaction := []byte{'a', 'b', 'c'}
	secondTransaction := []byte{'d', 'e', 'f'}

	initialBlockTransactions = append(initialBlockTransactions, firstTransaction)
	initialBlockTransactions = append(initialBlockTransactions, secondTransaction)

	// A block contains the Transactions, PublicKey, PrevPublicKey, and Index fields
	// respectively

	blockChain := &BlockChain{
		Blocks: []*Block{},
	}

	blockChain.addGenesisBlock(initialBlockTransactions)

	// Create a second block for testing purposes
	firstBlock := blockChain.getLatestBlock()

	// Build the second 2d slice
	secondBlockTransactions := [][]byte{}
	thirdTransaction := []byte{'g', 'h', 'i'}
	fourthTransaction := []byte{'j', 'k', 'l'}

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
	secondBlock := blockChain.getLatestBlock()

	thirdBlockTransactions := [][]byte{}
	fifthTransaction := []byte{'m', 'n', 'o'}
	sixthTransaction := []byte{'p', 'q', 'r'}

	thirdBlockTransactions = append(thirdBlockTransactions, fifthTransaction)
	thirdBlockTransactions = append(thirdBlockTransactions, sixthTransaction)

	newerBlock, err := createBlock(secondBlock, thirdBlockTransactions)

	blockChain.addBlock(newerBlock)

	// Testing reading/writing
	blockChain.write([]byte("hey it is me the guy"))
	blockChain.write([]byte("I am a beast and you are not"))

	spew.Dump(blockChain)
}
