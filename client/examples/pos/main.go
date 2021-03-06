package main

import (
	// TODO: import blockchain and related structs as bc
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	//"bytes"
	"os"
	"strings"
	"encoding/json"
)

// TODO: Table data and fields need to be encrypted during tableCreation and data entry

// Transfer : transaction used to move funds
type Transfer struct {
	ToAcctID      []byte
	Ammount       float64
	FromAcctID    []byte
}

// TableCreation : transaction used to create a new table
type TableCreation struct {
	ID     []byte
	Fields []string
	Types  []string
	//special check, owner in permission must == AcctID of publicKey used by signer
	PermissionByTable *TablePermission
	//                      AcctID : UserPermission
  PermissionByAcct    map[string]UserPermission
	FromAcctID          []byte
}

// Write : transaction used to write data to a table
type Write struct {
	TableID           []byte
	Cells             []*Cell
	//                [cell1Data, cell2Data, ...]
	Data              [][]byte
	FromAcctID        []byte
}

// Cell : identifies the column and row af the conceptual table
type Cell struct {
	X *big.Int
	Y *big.Int
}

// Edit : transaction used to change the value of cells within a table
type Edit struct {
	TableID []byte
	Cells   []*Cell
	//                [cell1Data, cell2Data, ...]
	NewDataByCell     [][]byte
  FromAcctID        []byte
}

// DeletionStub : holds onto data edited over
type DeletionStub struct {
	TableID           []byte
	Cells             []RetiredCellInfo
}

// Delete : transaction used to delete rows from a table
type Delete struct {
	TableID           []byte
	Cells             []*Cell
  FromAcctID        []byte
}

// ChangePermissions : transaction that grants an account permissions for a table
type ChangePermissions struct {
	TableID             []byte
	PermissionByTable   *TablePermission
	//                  acctID : UserPermission
  PermissionByAcct    map[string]UserPermission
  FromAcctID          []byte
}

// BlockGenerationBid : used to offer next block to network at a price
//Future: Create system to track market cost for each transaction type
type BlockGenerationBid struct {
	BidPrice      float64
	Stake         float64
	BlockNumber   *big.Int
	EstGenTime    float64
  FromAcctID    []byte
}

// GeneratedBlock : used to send generated block from bid winner to network
type GeneratedBlock struct {
	BlockHeight   big.Int
	// TODO: merge with blockchain and block type
  //CreatedBlock  bc.Block
  // Future: Updated Metadata could be a delta of current Metadata
	UpdatedMD     *DebasedMetadata
}

// Bet : used for staking during PoS
//Future: bets can contain evidence if believe block is wrong
type Bet struct {
	Stake         float64
	Position      bool
	Confidence    int
	Round         int
	BlockHash     []byte
	FromAcctID    []byte
}

// UserPermission : stores a user's access to a given table
type UserPermission struct {
	//only one owver allowed (gets billed and can assign permissions to others)
	//Multiple Admins allowed (can assign permissions to others)
	//ReadHistory is currently unused
	//read is currently unused
	//Owner,	Admin, WriteData,	EditData , DeleteData,	DeleteTable, ReadHistory,	ReadData
	Roles [8]bool
}

// AccountInfo : stores account balance and permissions
type AccountInfo struct {
	LiquidBalance float64
	//total amount currently staked
	IlliquidBalance float64
	//Permissions key:TableID value:Permission
	Permissions map[string]UserPermission
}

// RecordLocation : stores the block number and position of a transaction
type RecordLocation struct {
	BlockNumber big.Int
	//position 0 is the first transaction in a block
	Position big.Int
}

// RetiredCellInfo : used to store location of delete data
type RetiredCellInfo struct {
	Cell2DCord        Cell
	Location          CellLocation
}

// CellLocation : used to map cells from a table to the blockchain
type CellLocation struct {
	BlockNumber *big.Int
	//position 0 is the first transaction in a block
	Position *big.Int
	//byte position the cell begins at
	PostionInRecord *big.Int
}

// TablePermission : stores all accounts with a given access level
type TablePermission struct {
	Owner []byte
	//Multiple Admins allowed (can assign permissions to others)
	Admins [][]byte
	//readers is currently unused
	Readers       [][]byte
	Writers       [][]byte
	Editers       [][]byte
	DataDeleters  [][]byte
	TableDeleters [][]byte
	//HistoryReaders is currently unused
	HistoryReaders [][]byte
}

// TableInfo : stores tableSchema, location of each row, which accounts have what access
type TableInfo struct {
	CreationStub RecordLocation
	ID           []byte
	Fields       []string
	Types        []string
	//position 0 is the oldest
	//               row column
	Cells             [][]CellLocation
	Writes            []RecordLocation
	Edits             []DeletionStub
	Deletions         []DeletionStub
	DeadTable         bool
	Permission        TablePermission
}

// DebasedMetadata : stores account balance, permissions, and table info
type DebasedMetadata struct {
	//Accounts key:AcctNumber value: AcctInfo
	Accounts map[string]AccountInfo
	//Tables key:TableID value: TableInfo
	Tables map[string]TableInfo
}

// Transactions : stores slices of every transaction type
type Transactions struct {
	Transfers         []*Transfer
	TableCreations    []*TableCreation
	Writes            []*Write
	Edits             []*Edit
	Deletes           []*Delete
	PermissionChanges []*ChangePermissions
}

//DebasedSystem : model for the nodes' entire view of the debased pos/blockchain system
type DebasedSystem struct {
	PrivateKeysFromSession  []*ecdsa.PrivateKey
	// TODO: Integrate with blockchain
	CurrentBlockHeight      big.Int
	DuringConsensus         bool
	//Blockchain            *bc.Blockchain
	Metadata                *DebasedMetadata
	CurrentBids             []*BlockGenerationBid
	UnconfirmedBlock        *GeneratedBlock
  CurrentBets             []*Bet
	PendingBetPayouts       []*Transfer
	PendingTransactions     *Transactions
	HoldingPenTransactions  *Transactions
}

// GenerateBlock : creates a new block using the given DebasedSystem state
func (debasedS *DebasedSystem) GenerateBlock() GeneratedBlock {
	newDebasedMD := debasedS.Metadata
	var newBlockHeight big.Int
	newBlockHeight.Add(&debasedS.CurrentBlockHeight, big.NewInt(1))
	// Not needed pending new interface to Blockchain
	// TODO: UPDATE currentRecordLocation whenever a transaction is added to the block
	currentRecordLocation := RecordLocation{BlockNumber: newBlockHeight, Position: *big.NewInt(0)}
	// Not needed pending new interface to Blockchain/
	// TODO: UPDATE currentRecordLocation whenever a transaction or is added to the block and intra-transaction when adding data
	currentCellLocation := CellLocation{BlockNumber: big.NewInt(0), Position: big.NewInt(0), PostionInRecord: big.NewInt(0)}
	for _, transfer := range debasedS.PendingBetPayouts {
    var x = newDebasedMD.Accounts[string(transfer.FromAcctID[:])]
		x.IlliquidBalance -= transfer.Ammount
		newDebasedMD.Accounts[string(transfer.FromAcctID[:])] = x
		x = newDebasedMD.Accounts[string(transfer.ToAcctID[:])]
		x.IlliquidBalance -= transfer.Ammount
		x.LiquidBalance += transfer.Ammount * 2
  	newDebasedMD.Accounts[string(transfer.ToAcctID[:])] = x
	}
	for _, transfer := range debasedS.PendingTransactions.Transfers {
		var x = newDebasedMD.Accounts[string(transfer.FromAcctID[:])]
		x.IlliquidBalance -= transfer.Ammount
		newDebasedMD.Accounts[string(transfer.FromAcctID[:])] = x
		if acct, exist := newDebasedMD.Accounts[string(transfer.ToAcctID[:])]; exist {
      acct.LiquidBalance += transfer.Ammount
    } else {
			newDebasedMD.Accounts[string(transfer.ToAcctID[:])] = AccountInfo{transfer.Ammount, 0, make(map[string]UserPermission)}
		}
	}
	for _, create := range debasedS.PendingTransactions.TableCreations {
		create.PermissionByTable.Owner = create.FromAcctID
    newDebasedMD.Tables[string(create.ID[:])] = TableInfo{CreationStub: currentRecordLocation,
			                                         ID: create.ID,
																							 Fields: create.Fields,
                                               Types: create.Types,
																							 Cells: make([][]CellLocation, 0),
																							 Writes: make([]RecordLocation, 0),
																							 Edits: make([]DeletionStub, 0),
																							 Deletions: make([]DeletionStub, 0),
																							 DeadTable: false,
																							 Permission : *(create.PermissionByTable),
		                                          }
		for acctID, userPermish := range create.PermissionByAcct {
			newDebasedMD.Accounts[acctID].Permissions[string(create.ID[:])] = userPermish
		}
	}
	for _, add := range debasedS.PendingTransactions.Writes {
    //check is FromAcctID has write access to table
		if newDebasedMD.Accounts[string(add.FromAcctID[:])].Permissions[string(add.TableID[:])].Roles[2] {
			// TODO: add data to the blockchain
			var x = newDebasedMD.Tables[string(add.TableID[:])]
      x.Writes = append(newDebasedMD.Tables[string(add.TableID[:])].Writes, currentRecordLocation)
			newDebasedMD.Tables[string(add.TableID[:])] = x
			for _, eachCell := range add.Cells {
				newDebasedMD.Tables[string(add.TableID[:])].Cells[eachCell.Y.Uint64()][eachCell.X.Uint64()] = currentCellLocation
			}
		}
	}
	for _, editRequest := range debasedS.PendingTransactions.Edits {
    //check is FromAcctID has edit access to table
		if newDebasedMD.Accounts[string(editRequest.FromAcctID[:])].Permissions[string(editRequest.TableID[:])].Roles[3] {
			var deletionRecord DeletionStub
			deletionRecord.TableID = editRequest.TableID
      for _, cell := range editRequest.Cells {
				deletionRecord.Cells = append(deletionRecord.Cells, RetiredCellInfo{Cell{cell.X, cell.Y}, newDebasedMD.Tables[string(editRequest.TableID[:])].Cells[cell.X.Uint64()][cell.Y.Uint64()]})
			}
			var x = newDebasedMD.Tables[string(editRequest.TableID[:])]
			x.Edits = append(newDebasedMD.Tables[string(editRequest.TableID[:])].Edits, deletionRecord)
			newDebasedMD.Tables[string(editRequest.TableID[:])] = x
			// TODO: add data to the blockchain
			for _, eachCell := range editRequest.Cells {
				newDebasedMD.Tables[string(editRequest.TableID[:])].Cells[eachCell.Y.Uint64()][eachCell.X.Uint64()] = currentCellLocation
			}
		}
	}
	for _, deltionRequest := range debasedS.PendingTransactions.Deletes {
		if newDebasedMD.Accounts[string(deltionRequest.FromAcctID[:])].Permissions[string(deltionRequest.TableID[:])].Roles[4] {
			var deletionRecord DeletionStub
			deletionRecord.TableID = deltionRequest.TableID
      for _, cell := range deltionRequest.Cells {
				deletionRecord.Cells = append(deletionRecord.Cells, RetiredCellInfo{Cell{cell.X, cell.Y}, newDebasedMD.Tables[string(deltionRequest.TableID[:])].Cells[cell.X.Uint64()][cell.Y.Uint64()]})
			}
			var x = newDebasedMD.Tables[string(deltionRequest.TableID[:])]
			x.Edits = append(newDebasedMD.Tables[string(deltionRequest.TableID[:])].Edits, deletionRecord)
			newDebasedMD.Tables[string(deltionRequest.TableID[:])] = x
			for _, eachCell := range deltionRequest.Cells {
				newDebasedMD.Tables[string(deltionRequest.TableID[:])].Cells[eachCell.Y.Uint64()][eachCell.X.Uint64()] = CellLocation{}
			}
		}
	}
	return GeneratedBlock{newBlockHeight, newDebasedMD}
}

//CheckUnconfirmedBlock : used to verify a GeneratedBlock received from the network
func (debasedS *DebasedSystem) CheckUnconfirmedBlock() bool {
	var nextBlockHeight big.Int
	nextBlockHeight.Add(&debasedS.CurrentBlockHeight, big.NewInt(1))
	if debasedS.UnconfirmedBlock.BlockHeight.Uint64() != nextBlockHeight.Uint64() {
		return false
	}
  return reflect.DeepEqual(debasedS.UnconfirmedBlock, debasedS.GenerateBlock)
}

// JSONWrapper : used to sign and verify obj sent overnetwork
type JSONWrapper struct {
	PK        *ecdsa.PublicKey
	R         *big.Int
	S         *big.Int
	Type      string
	//JSON of the encolesed struct
	Contents  []byte
}

// Sign : assgins PK, R, S to JSONWrapper
func (wrapper *JSONWrapper) Sign(privateKey *ecdsa.PrivateKey) error{
	var err error
	wrapper.R, wrapper.S, err = ecdsa.Sign(rand.Reader, privateKey, append([]byte(wrapper.Type), wrapper.Contents...))
	wrapper.PK = &privateKey.PublicKey
	return err
}

// VerifySignature : verifies the signature in JSONWrapper
func (wrapper *JSONWrapper) VerifySignature() bool {
	return ecdsa.Verify(wrapper.PK, append([]byte(wrapper.Type), wrapper.Contents...), wrapper.R, wrapper.S)
}

//NOW TODOs
// TODO: MAKE consensus occur on a timer
//       A node creates a block a sends it to the network
//       Each node then checks the block
//       The nodes each place their votes
//       After X time from the strat of the process, the voting is closed
//       IF approved:
//                    determine payouts and clear CurrentBets
//										each node updates their MD and BC
//                    each node flushes PendingTransactions and PendingBetPayouts
//                    and move holding bay transactions into pending
//       IF not repeat
//                    determine payouts and clear CurrentBets
//										repeat process

//BIG TESTs
// Be able to generate a block and update metadata

//BIG TODOs
// TODO: Be able to convert between struct <----> JSON
// TODO: Be able to check if a transaction conforms to debased rules
// TODO: Be able to choose a node to gen next block
// TODO: Be able to check the correctness of a newly generate block and metadata
// TODO: Be able to bet
// TODO: Be able to recognize consensus
// TODO: Be able to payout accounts with correct votes (part of block generation)

//BIG FUTUREs (WSB STYLE)
// Future: have decearnment on which transactions to include vs exclude in newly generated blocks
// Future: avoid dead transactions by allowing killing
// Future: voting order and complicated payouts
// Future: confidence system for inter-node relations
// Future: track accounts/nodes with skin in the game, how much, and where
// Future: table deletion
// Future: PermissionChanges

//AccountNumber : determines account number from PublicKey
func AccountNumber(publicKey ecdsa.PublicKey) []byte {
	hash := sha256.Sum256(append(publicKey.X.Bytes(), publicKey.Y.Bytes()...))
	return hash[12:]
}

//creates account that can be used on the debased network
func createAcct() (*ecdsa.PrivateKey, []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return privateKey, AccountNumber(privateKey.PublicKey)
}

func main() {
	nodeDebasedSystem := DebasedSystem{PrivateKeysFromSession: make([]*ecdsa.PrivateKey, 0),
		                                 CurrentBlockHeight: *big.NewInt(0),
																		 DuringConsensus: false,
																		 Metadata: &DebasedMetadata{Accounts: make(map[string]AccountInfo), Tables: make(map[string]TableInfo)},
																		 CurrentBids: make([]*BlockGenerationBid, 0),
																		 UnconfirmedBlock: nil,
																		 CurrentBets: make([]*Bet, 0),
																		 PendingBetPayouts: make([]*Transfer, 0),
																		 PendingTransactions: &Transactions{},
																		 HoldingPenTransactions: &Transactions{},
	                                   }

	dummyMetadata := dummyDebasedMetaData()
	fmt.Printf("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" {
			os.Exit(0)
		}
		args := strings.Fields(line)

		if args[0] == "checkBalance" && len(args) == 2{
			if val, ok := nodeDebasedSystem.Metadata.Accounts[args[1]]; ok {
				balance := val.IlliquidBalance + val.LiquidBalance
				fmt.Println("Account balance:", balance)
			} else {
				fmt.Println("Account not found")
			}
		}
		if args[0] == "checkRoles" {
			if _, ok := dummyMetadata.Accounts[args[1]]; !ok {
				fmt.Println("Account not found")
			} else {
				allRoles := []string{"Owner", "Admin", "WriteData", "EditData", "DeleteData", "DeleteTable", "ReadHistory", "ReadData"}
				currentMsg := ""
				for k, v := range dummyMetadata.Accounts[args[1]].Permissions {
					currentMsg = "TableID: " + k + "\n" + "Roles: "
					for index, role := range allRoles {
						if v.Roles[index] {
							currentMsg = currentMsg + role + ", "
						}
					}
					fmt.Println(currentMsg)
				}
			}
		}
		if args[0] == "tableAttributes" {
			if _, ok := dummyMetadata.Tables[args[1]]; !ok {
				fmt.Println("Table not found")
			} else {
				currentMsg := "["
				for _, v := range dummyMetadata.Tables[args[1]].Fields {
					currentMsg = currentMsg + v + ", "
				}
				currentMsg = currentMsg + "]"
				fmt.Println(currentMsg)
			}
		}
		if args[0] == "tableData" {
			if _, ok := dummyMetadata.Tables[args[1]]; !ok {
				fmt.Println("Table not found")
			} else {
				for _, row := range dummyMetadata.Tables[args[1]].Cells {
					fmt.Print("[")
					for _, column := range row {
						fmt.Print("(BN: ", column.BlockNumber.String())
						fmt.Print(", POS: ", column.Position.String())
						fmt.Print(", PIR: ", column.PostionInRecord.String(), "), ")
					}
					fmt.Println("]")
				}
			}
		}
		if args[0] == "AccessibleTables" {
			if _, ok := dummyMetadata.Accounts[args[1]]; !ok {
				fmt.Println("Account not found")
			} else {
				fmt.Println("TableIDs: ")
				for k := range dummyMetadata.Accounts[args[1]].Permissions {
					fmt.Print(k, ", ")
				}
				fmt.Println()
			}
		}
		if args[0] == "tableHistory" {
			if _, ok := dummyMetadata.Tables[args[1]]; !ok {
				fmt.Println("Table not found")
			} else {
				fmt.Print("[")
				for _, write := range dummyMetadata.Tables[args[1]].Writes {
					fmt.Print("(BN: ", write.BlockNumber.String())
					fmt.Print(", POS: ", write.Position.String(), "), ")
				}
				fmt.Println("]")
			}
		}

		if args[0] == "accounts" {
			for _, account := range dummyMetadata.Accounts {
				fmt.Printf("%+v\n", account)
			}
		}

		if args[0] == "tables" {
			for _, table := range dummyMetadata.Tables {
				fmt.Printf("%+v\n", table)
			}
		}

		if args[0] == "addAccount" {
			if _, ok := dummyMetadata.Accounts[args[1]]; !ok {
				initialPermission := make(map[string]UserPermission)
				for key, account := range dummyMetadata.Accounts {
					initialPermission[key] = UserPermission{
						Roles: [8]bool{
							false, false, false, false, false, false, false, false,
						},
					}
					account.Permissions[args[1]] = UserPermission{
						Roles: [8]bool{
							false, false, false, false, false, false, false, false,
						},
					}
				}
				dummyMetadata.Accounts[args[1]] = AccountInfo{
					LiquidBalance:   1,
					IlliquidBalance: 1,
					Permissions:     initialPermission,
				}
				fmt.Printf("%+v\n", dummyMetadata.Accounts[args[1]])
			} else {
				fmt.Println("No accountID provided")
			}
		}

		if args[0] == "deleteAccount" {
			if _, ok := dummyMetadata.Accounts[args[1]]; ok {
				delete(dummyMetadata.Accounts, args[1])
				for _, account := range dummyMetadata.Accounts {
					if _, ok := account.Permissions[args[1]]; ok {
						delete(account.Permissions, args[1])
					}
				}
			} else {
				fmt.Println("Account does not exist")
			}
		}
		if args[0] == "createAcct" && len(args) == 1 {
      privateKey, acctID := createAcct()
			nodeDebasedSystem.PrivateKeysFromSession = append(nodeDebasedSystem.PrivateKeysFromSession, privateKey)
			fmt.Println(privateKey)
			fmt.Println(acctID)
		}
		// transfer ToAcct Amount FromAcctPrivateKeyIndex
		if args[0] == "transfer" && len(args) == 4 {
			var ammount float64
			var err error
			if ammount, err = strconv.ParseFloat(args[2], 64); err != nil {
        panic(err)
      }
			var payment Transfer
			toAcctIndex, err := strconv.ParseInt(args[3], 10, 64)
			payment = Transfer{[]byte(args[1]), ammount, AccountNumber(nodeDebasedSystem.PrivateKeysFromSession[toAcctIndex].PublicKey)}
			if !nodeDebasedSystem.DuringConsensus {
				nodeDebasedSystem.PendingTransactions.Transfers = append(nodeDebasedSystem.PendingTransactions.Transfers, &payment)
				fmt.Println(nodeDebasedSystem.PendingTransactions.Transfers[len(nodeDebasedSystem.PendingTransactions.Transfers)-1])
			} else {
				nodeDebasedSystem.HoldingPenTransactions.Transfers = append(nodeDebasedSystem.HoldingPenTransactions.Transfers, &payment)
				fmt.Println(nodeDebasedSystem.HoldingPenTransactions.Transfers[len(nodeDebasedSystem.HoldingPenTransactions.Transfers)-1])
			}
		}
		// genBlock currentAcctIndexInCurrentAccts
		if args[0] == "genBlock" && len(args) == 2 {
			// fmt.Println(nodeDebasedSystem.GenerateBlock())
			var newBlock = nodeDebasedSystem.GenerateBlock()
			b, err := json.Marshal(newBlock)
	    if err != nil {
		    fmt.Println("error:", err)
	    }
			wrapper := &JSONWrapper{Type:"GeneratedBlock", Contents:b}
			toAcctIndex, err := strconv.ParseInt(args[1], 10, 64)
			wrapper.Sign(nodeDebasedSystem.PrivateKeysFromSession[toAcctIndex])
			//fmt.Println("ORIG WRAPPER")
			fmt.Println(wrapper)
			fmt.Println(wrapper.VerifySignature())

			//fmt.Println("WRAPPER AFTER MARSHAL UNMARSHAL")
			//if err != nil {
		  //  fmt.Println("error:", err)
	    //}
			//jsonOfWrapper, err := json.Marshal(*wrapper)
			//var unMarshaledWrapper JSONWrapper
			//err = json.Unmarshal(jsonOfWrapper, &unMarshaledWrapper)
			//if err != nil {
		  //  fmt.Println("error:", err)
	    //}
			//fmt.Println(unMarshaledWrapper)
      //fmt.Println("Signature is:", (*x).VerifySignature())
			//fmt.Println("Signer matches currect account:", )
			//nodeDebasedSystem.CurrentBlockHeight = newBlock.BlockHeight
			//nodeDebasedSystem.Metadata = newBlock.UpdatedMD
			//fmt.Printf("%+v\n", newBlock)
		}

		if args[0] == "never" {
			fmt.Println(dummyMetadata)
		}
		fmt.Printf("> ")
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
