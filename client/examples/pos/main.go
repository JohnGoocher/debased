package main

import (
	// TODO: import blockchain and related structs as bc
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"bufio"
  "os"
	"fmt"
	"math/big"
)

// TODO: Table data and fields need to be encrypted during tableCreation and data entry

// Transfer : transaction used to move funds
type Transfer struct {
	ToAcctID      []byte
	Ammount       float
	fromAcctID    []byte
}

// TableCreation : transaction used to create a new table
type TableCreation struct {
	ID                  []byte
	Fields              []string
	Types               []string
  //special check, owner in permission must == AcctID of publicKey used by signer
	PermissionByTable   *TablePermission
	//                      AcctID : UserPermission
  PermissionByAcct    map[string]UserPermission
	fromAcctID          []byte
}

// Write : transaction used to write data to a table
type Write struct {
	TableID           []byte
	//               row column data
	Data              [][][]byte
	fromAcctID        []byte
}

// Cell : identifies the column and row af the conceptual table
type Cell struct {
	X                 *big.int
	Y                 *big.int
}

// Edit : transaction used to change the value of cells within a table
type Edit struct {
	TableID           []byte
	Cells             []*Cell
	//                [cell1Data, cell2Data, ...]
	NewDataByCell     [][]byte
  fromAcctID        []byte
}

// Delete : transaction used to delete rows from a table
type Delete struct {
	TableID           []byte
	Rows              []*big.int
  fromAcctID        []byte
}

// ChangePermissions : transaction that grants an account permissions for a table
type ChangePermissions struct {
	TableID             []byte
	PermissionByTable   *TablePermission
  PermissionByAcct    map[string]UserPermission
  fromAcctID          []byte
}

// BlockGenerationBid : used to offer next block to network at a price
//Future: Create system to track market cost for each transaction type
type BlockGenerationBid struct {
	BidPrice      float
	Stake         float
	BlockNumber   *big.int
	EstGenTime    float
  fromAcctID    []byte
}

// GeneratedBlock : used to send generated block from bid winner to network
type GeneratedBlock struct {
	// TODO: merge with blockchain and block type
  CreatedBlock  bc.Block
  // Future: Updated Metadata could be a delta of current Metadata
	UpdatedMD     DebasedMetadata
}

// Bet : used for staking during PoS
//Future: bets can contain evidence if believe block is wrong
type Bet struct {
	Stake         float
	Position      bool
	Confidence    int
	Round         int
	BlockHash     []byte
	fromAcctID    []byte
}

// UserPermission : stores a user's access to a given table
type UserPermission struct {
	//only one owver allowed (gets billed and can assign permissions to others)
	//Multiple Admins allowed (can assign permissions to others)
	//ReadHistory is currently unused
  //read is currently unused
	//Owner,	Admin, WriteData,	EditData , DeleteData,	DeleteTable, ReadHistory,	ReadData
	Roles             [8]bool
}

// AccountInfo : stores account balance and permissions
type AccountInfo struct {
	LiquidBalance     float
  //total amount currently staked
	IlliquidBalance   float
	//Permissions key:TableID value:Permission
	Permissions       map[string]UserPermission
}

// RecordLocation : stores the block number and position of a transaction
type RecordLocation struct {
	BlockNumber       big.int
	//position 0 is the first transaction in a block
	Position          big.int
}

// CellLocation : used to map cells from a table to the blockchain
type CellLocation struct {
	BlockNumber       *big.int
	//position 0 is the first transaction in a block
	Position          *big.int
	//byte position the cell begins at
	PostionInRecord   *big.int
}

// TablePermission : stores all accounts with a given access level
type TablePermission struct {
	Owner             []byte
	//Multiple Admins allowed (can assign permissions to others)
	Admins            [][]byte
	//readers is currently unused
	Readers           [][]byte
	Writers           [][]byte
	Editers           [][]byte
	DataDeleters      [][]byte
	TableDeleters     [][]byte
	//HistoryReaders is currently unused
	HistoryReaders    [][]byte
}

// TableInfo : stores tableSchema, location of each row, which accounts have what access
type TableInfo struct {
	CreationStub      RecordLocation
	ID                []byte
	Fields            []string
	Types             []string
	//position 0 is the oldest
	//               row column
	Cells             [][]*CellLocation
	Writes            []*RecordLocation
	Edits             []*RecordLocation
	Deletions         []*RecordLocation
	DeadTable         bool
	Permission        TablePermission
}

// DebasedMetadata : stores account balance, permissions, and table info
type DebasedMetadata struct {
	//Accounts key:AcctNumber value: AcctInfo
  Accounts          map[string]AccountInfo
	//Tables key:TableID value: TableInfo
	Tables            map[string]TableInfo
}

// Transactions : stores slices of every transaction type
type Transactions struct {
	Transfers           []*Transfer
	TableCreations      []*TableCreation
	Writes              []*Write
	Edits               []*Edit
	Deletes             []*Delete
	PermissionChanges   []*ChangePermissions
}

//DebasedSystem : model for the nodes' entire view of the debased pos/blockchain system
type DebasedSystem struct {
	// TODO: Integrate with blockchain
	CurrentBlockHeight    big.int
	Blockchain            *bc.Blockchain
	Metadata              *DebasedMetadata
	CurrentBids           []*BlockGenerationBid
	UnconfirmedBlock      *GeneratedBlock
  CurrentBets           []*Bet
	PendingBetPayouts     []*Transfer
	PendingTransactions   *Transactions
	// Future: voting order and complicated payouts
	// Future: confidence system for inter-node relations
	// TODO: track accounts/nodes with skin in the game, how much, and where
	// avoid stuck transactions by allowing killing
}

// GenerateBlock : creates a new block using the given DebasedSystem state
func (debasedS DebasedSystem) GenerateBlock() *GeneratedBlock {
	newDebasedMD := debasedS.Metadata
	newDebasedMD.CurrentBlockHeight++
	// TODO: UPDATE currentRecordLocation whenever a transaction is added to the block
	currentRecordLocation = RecordLocation{BlockNumber: newDebasedMD.CurrentBlockHeight, Position: 0}
	for _, transfer := range debasedS.PendingBetPayouts {
    newDebasedMD.Accounts[string(transfer.fromAcctID[:])].IlliquidBalance -= transfer.Ammount
    newDebasedMD.Accounts[string(transfer.ToAcctID[:])].IlliquidBalance -= transfer.Ammount
  	newDebasedMD.Accounts[string(transfer.ToAcctID[:])].LiquidBalance += transfer.Ammount * 2
	}
	for _, transfer := range debasedS.PendingTransactions.Transfers {
		newDebasedMD.Accounts[string(transfer.fromAcctID[:])].IlliquidBalance -= transfer.Ammount
		//HANDLE IF KEY NOT IN map
		if val, ok := newDebasedMD.Accounts[string(transfer.ToAcctID[:])]; ok {
      val.LiquidBalance += transfer.Ammount
    } else {
			newDebasedMD.Accounts[string(transfer.ToAcctID[:])] = AccountInfo{0, transfer.Ammount, make(map[string]UserPermission)}
		}
	}
	for _, create := range debasedS.PendingTransactions.TableCreations {
		create.Permission.Owner = fromAcctID
    newDebasedMD.Tables[create.ID] = TableInfo{CreationStub: currentRecordLocation,
			                                         ID: create.ID,
																							 Fields: create.Fields,
                                               Types: create.Types,
																							 Cells: make([][]*CellLocation, 0),
																							 Writes: make([]*RecordLocation, 0),
																							 Edits: make([]*RecordLocation, 0),
																							 Deletions: make([]*RecordLocation, 0),
																							 DeadTable: false,
																							 Permission : &(create.PermissionByTable),
		                                          }
		for acctID, userPermish := range create.PermissionByAcct {
        newDebasedMD.Accounts[acctID].Permissions[create.ID] = userPermish
    }
	}
	for _, add := range debasedS.PendingTransactions.Writes {
    //check is fromAcctID has write access to table

	}
}

//BIG TODOs
// TODO: Be able to convert between struct <----> JSON
// TODO: Be able to check if a transaction conforms to debased rules
// TODO: Be able to choose a node to gen next block
// TODO: Be able to generate a block and update metadata
// TODO: Be able to check the correctness of a newly generate block and metadata
// TODO: Be able to bet
// TODO: Be able to recognize consensus
// TODO: Be able to payout accounts with correct votes (part of block generation)


//AccountNumber : determines account number from PublicKey
func (publicKey ecdsa.PublicKey) AccountNumber() []byte {
	return sha256.Sum256(append(publicKey.X.Bytes(),publicKey.Y.Bytes()...))[12:]
}

//creates account that can be used on the debased network
func createAcct() (*ecdsa.PrivateKey, []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return privateKey, &privateKey.PublicKey.AccountNumber()
}

func main() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: (0x%x, 0x%x)\n", r, s)

	valid := ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s)
	fmt.Println("signature verified:", valid)


	scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    line := scanner.Text()
		//money has to be transfered to the new acct to notify the network of its creation
		if line == "create account" {
			privateKey, accountNumber, err := createAcct()
			//does NOT actually return error
			if err != nil {
				panic(err)
			}
			// TODO: Need to format both of these
			fmt.Println(privateKey)
			fmt.Println(accountNumber)
		}
    if line == "exit" {
      os.Exit(0)
    }
    args := strings.Fields(line)
		//transfer fromPrivKey toPubKey ammount
		if args[0] == "transfer" {
			// TODO: Create Transfer from CLI
		}
  }
  if err := scanner.Err(); err != nil {
      fmt.Fprintln(os.Stderr, "reading standard input:", err)
  }
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
