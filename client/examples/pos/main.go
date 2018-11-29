package main

import (
	// TODO: import blockchain and related structs as bc
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	mrand "math/rand"
	"reflect"
	"strconv"

	//"bytes"
	"encoding/json"
	"os"
	"strings"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// TODO: Table data and fields need to be encrypted during tableCreation and data entry

// Transfer : transaction used to move funds
type Transfer struct {
	ToAcctID   []byte
	Ammount    float64
	FromAcctID []byte
}

// TableCreation : transaction used to create a new table
type TableCreation struct {
	ID     []byte
	Fields []string
	Types  []string
	//special check, owner in permission must == AcctID of publicKey used by signer
	PermissionByTable *TablePermission
	//                      AcctID : UserPermission
	PermissionByAcct map[string]UserPermission
	FromAcctID       []byte
}

// Write : transaction used to write data to a table
type Write struct {
	TableID []byte
	Cells   []*Cell
	//                [cell1Data, cell2Data, ...]
	Data       [][]byte
	FromAcctID []byte
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
	NewDataByCell [][]byte
	FromAcctID    []byte
}

// DeletionStub : holds onto data edited over
type DeletionStub struct {
	TableID []byte
	Cells   []RetiredCellInfo
}

// Delete : transaction used to delete rows from a table
type Delete struct {
	TableID    []byte
	Cells      []*Cell
	FromAcctID []byte
}

// ChangePermissions : transaction that grants an account permissions for a table
type ChangePermissions struct {
	TableID           []byte
	PermissionByTable *TablePermission
	//                  acctID : UserPermission
	PermissionByAcct map[string]UserPermission
	FromAcctID       []byte
}

// BlockGenerationBid : used to offer next block to network at a price
//Future: Create system to track market cost for each transaction type
type BlockGenerationBid struct {
	BidPrice    float64
	Stake       float64
	BlockNumber *big.Int
	EstGenTime  float64
	FromAcctID  []byte
}

// GeneratedBlock : used to send generated block from bid winner to network
type GeneratedBlock struct {
	BlockHeight big.Int
	// TODO: merge with blockchain and block type
	//CreatedBlock  bc.Block
	// Future: Updated Metadata could be a delta of current Metadata
	UpdatedMD *DebasedMetadata
}

// Bet : used for staking during PoS
//Future: bets can contain evidence if believe block is wrong
type Bet struct {
	Stake      float64
	Position   bool
	Confidence int
	Round      int
	BlockHash  []byte
	FromAcctID []byte
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
	Cell2DCord Cell
	Location   CellLocation
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
	Cells      [][]CellLocation
	Writes     []RecordLocation
	Edits      []DeletionStub
	Deletions  []DeletionStub
	DeadTable  bool
	Permission TablePermission
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
	PrivateKeysFromSession []*ecdsa.PrivateKey
	// TODO: Integrate with blockchain
	CurrentBlockHeight big.Int
	DuringConsensus    bool
	//Blockchain            *bc.Blockchain
	Metadata               *DebasedMetadata
	CurrentBids            []*BlockGenerationBid
	UnconfirmedBlock       *GeneratedBlock
	CurrentBets            []*Bet
	PendingBetPayouts      []*Transfer
	PendingTransactions    *Transactions
	HoldingPenTransactions *Transactions
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
			ID:         create.ID,
			Fields:     create.Fields,
			Types:      create.Types,
			Cells:      make([][]CellLocation, 0),
			Writes:     make([]RecordLocation, 0),
			Edits:      make([]DeletionStub, 0),
			Deletions:  make([]DeletionStub, 0),
			DeadTable:  false,
			Permission: *(create.PermissionByTable),
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
	PK   *ecdsa.PublicKey
	R    *big.Int
	S    *big.Int
	Type string
	//JSON of the encolesed struct
	Contents []byte
}

// Sign : assgins PK, R, S to JSONWrapper
func (wrapper *JSONWrapper) Sign(privateKey *ecdsa.PrivateKey) error {
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

type Client struct {
	testMap map[string]string

	// rw map[peer.ID]*bufio.ReadWriter
	rw map[net.Stream]*bufio.ReadWriter

	// streams map[peer.ID]net.Stream
	streams map[net.Stream]net.Stream
	stream  net.Stream
	host    host.Host
	json    *JSONWrapper
}

/*
 * addAddrToPeerstore parses a peer multiaddress and adds
 * it to the given host's peerstore, so it knows how to
 * contact it. It returns the peer ID of the remote peer.
 * @credit examples/http-proxy/proxy.go
 */
func addAddrToPeerstore(h host.Host, addr string) peer.ID {
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		log.Fatalln(err)
	}

	info, err := peerstore.InfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatalln(err)
	}

	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	// fmt.Println("AM I THE PEER ID THINGY")
	// fmt.Printf("%+v\n", info.ID)
	// fmt.Println("AM I THE PEER ID THINGY")
	return info.ID
}

func (c *Client) handleStream(s net.Stream) {
	log.Println("Got a new stream!")
	// fmt.Printf("s.Conn().LocalPeer: %+v\n", s.Conn().LocalPeer)
	// fmt.Printf("s.Conn().LocalPrivateKey: %+v\n", s.Conn().LocalPrivateKey)
	// fmt.Printf("s.Conn().LocalMultiaddr: %+v\n", s.Conn().LocalMultiaddr)
	// fmt.Printf("s.Conn().RemotePeer: %+v\n", s.Conn().RemotePeer)
	// fmt.Printf("s.Conn().RemotePublicKey: %+v\n", s.Conn().RemotePublicKey)
	// fmt.Printf("s.Conn().RemoteMultiaddr: %+v\n", s.Conn().RemoteMultiaddr)
	// fmt.Printf("s.Stat(): %+v\n", s.Stat())
	// fmt.Printf("s.Protocol(): %+v\n", s.Protocol())

	// uuid := xid.New()
	// fmt.Println(uuid.String())
	// temp := map[xid.ID]net.Stream{}

	if val, ok := c.rw[s]; ok {
		fmt.Println("It's there!")
		fmt.Printf("%+v", val)
	}

	// Create a buffer stream for non blocking read and write.
	// c.rw[c.host.ID()] = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	c.rw[s] = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	// fmt.Printf("c.host.ID().Pretty(): %+v\n", c.host.ID().Pretty())
	// fmt.Printf("BEFORE c.streams: %+v\n", c.streams)

	// if val, ok := dict["foo"]; ok {
	// 	//do something here
	// }

	c.streams[s] = s

	// =================================================
	// =================================================
	// host.SetStreamHandler("/chat/1.0.0", sampleClient.handleStream)

	// // Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
	// var port string
	// for _, la := range host.Network().ListenAddresses() {
	// 	if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
	// 		port = p
	// 		break
	// 	}
	// }
	// =================================================
	// =================================================

	// fmt.Printf("AFTER c.streams: %+v\n", c.streams)

	go c.readExampleData(s)
	// go c.writeExampleData(s)

	// go readData(rw)
	// go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func (c *Client) writeStreams(s net.Stream) {
	sendData, err := json.Marshal(c.streams)
	if err != nil {
		fmt.Println("JSON.MARSHALL PANIC")
		panic(err)
	}
	// fmt.Println("AFTER SEND DATA")
	// fmt.Println("before write")
	// fmt.Printf("sendData: %+v\n", string(sendData))

	// fmt.Println("STRRRRRREEEAMMMMSSSSSSS")
	// fmt.Printf("%+v\n", c.streams)

	c.rw[s].Write(sendData)
	c.rw[s].Flush()
}

func (c *Client) readStreams(s net.Stream) {
	str, err := c.rw[s].ReadSlice('}')

	// fmt.Printf("READING: %+v\n", s)
	// fmt.Println("STTTTRRRREEEEAAAAAMMMMS")
	// fmt.Printf("%+v\n", c.streams)

	// fmt.Println("AFTER READSLICE")
	if err != nil {
		// fmt.Println("READSLICE PANIC")
		panic(err)
	}
	// fmt.Println("after readSlice")
	fmt.Printf("readslice: %+v\n", string(str))

	// var someJson interface{}
	// json := &interface{}

	if len(str) > 0 {
		if err := json.Unmarshal(str, &c.json); err != nil {
			// fmt.Println("OR IS IT THIS PANIC")
			// panic(err)
			// continue
			fmt.Printf("%+v\n", c.json)
		}
		// fmt.Printf("c.json: %+v\n", c.json)
		// fmt.Println("END OF ELSE")
		// fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
	} else {
		fmt.Println("Receieved value is 0")
		return
	}
}

func (c *Client) readExampleData(s net.Stream) {
	// var testMap map[string]string
	for {
		// fmt.Println("FOR LOOP STARTED")
		// str, err := c.rw[c.host.ID()].ReadSlice('}')
		str, err := c.rw[s].ReadSlice('\n')

		str = str[:len(str)-1]
		// s = s[:sz-1]

		// fmt.Printf("READING: %+v\n", s)
		// fmt.Println("STTTTRRRREEEEAAAAAMMMMS")
		// fmt.Printf("%+v\n", c.streams)

		// fmt.Println("AFTER READSLICE")
		if err != nil {
			// fmt.Println("READSLICE PANIC")
			panic(err)
		}
		// fmt.Println("after readSlice")
		// fmt.Printf("readslice: %+v\n", string(str))

		// var someJson interface{}

		if len(str) > 0 {
			// if err := json.Unmarshal(str, &someJson); err != nil {
			if err := json.Unmarshal(str, &c.json); err != nil {
				// fmt.Println("OR IS IT THIS PANIC")
				// panic(err)
			}
			fmt.Printf("c.json: %+v\n", c.json)
			// fmt.Printf("someJson: %+v\n", someJson.PK)
			// fmt.Println("END OF ELSE")
			// fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		} else {
			fmt.Println("Receieved value is 0")
			return
		}
		// fmt.Println("END OF A LOOP")
	}
}

func (c *Client) writeExampleData(data []byte, s net.Stream) {
	// testMap := make(map[string]string)
	// stdReader := bufio.NewReader(os.Stdin)
	// count := 0
	// for {
	// fmt.Println("count: %i", count)
	// fmt.Print("before> ")
	// data, err := stdReader.ReadString('\n')
	// fmt.Printf("WRITING: %+v\n", s)
	// fmt.Println("AFTER READSTRING")
	// fmt.Println("data: " + data)

	// fmt.Println("WHAT")
	// fmt.Printf("c.testMap before: %+v\n", c.testMap)
	// fmt.Printf("data: %+v\n", data)
	// c.testMap[data] = data
	// fmt.Printf("c.testMap after: %+v\n", c.testMap)
	// fmt.Println("AFTER TESTMAP")
	// sendData, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println("JSON.MARSHALL PANIC")
	// 	panic(err)
	// }
	// fmt.Println("AFTER SEND DATA")
	// fmt.Println("before write")
	// fmt.Printf("sendData: %+v\n", string(sendData))
	// fmt.Println("STRRRRRREEEAMMMMSSSSSSS")
	// fmt.Printf("%+v\n", c.streams)

	data = append(data, '\n')

	for _, writer := range c.rw {
		writer.Write(data)
		writer.Flush()
	}

	// old way
	// c.rw[s].Write(sendData)
	// c.rw[s].Flush()
	// old way

	// count++
	// }
}

func BeginStream(client *Client, dest *string, host host.Host) {

	if *dest == "" {
		// Set a function as stream handler.
		// This function is called when a peer connects, and starts a stream with this protocol.
		// Only applies on the receiving side.
		host.SetStreamHandler("/chat/1.0.0", client.handleStream)

		// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
		var port string
		for _, la := range host.Network().ListenAddresses() {
			if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
				port = p
				break
			}
		}

		if port == "" {
			panic("was not able to find actual local port")
		}

		fmt.Printf("Run './chat-exec -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, host.ID().Pretty())
		fmt.Println("You can replace 127.0.0.1 with public IP as well.")
		fmt.Printf("\nWaiting for incoming connection\n\n")

		// Hang forever
		<-make(chan struct{})
	} else {
		fmt.Println("This node's multiaddresses:")
		for _, la := range host.Addrs() {
			fmt.Printf(" - %v\n", la)
		}
		fmt.Println()

		// Turn the destination into a multiaddr.
		maddr, err := multiaddr.NewMultiaddr(*dest)
		if err != nil {
			log.Fatalln(err)
		}

		// Extract the peer ID from the multiaddr.
		info, err := peerstore.InfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}

		host.SetStreamHandler("/chat/1.0.0", client.handleStream)
		fmt.Printf("Run './chat-exec -d /ip4/127.0.0.1/tcp//p2p/%s' on another console.\n", host.ID().Pretty())

		// Add the destination's peer multiaddress in the peerstore.
		// This will be used during connection and stream creation by libp2p.
		host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		// Start a stream with the destination.
		// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
		s, err := host.NewStream(context.Background(), info.ID, "/chat/1.0.0")
		if err != nil {
			panic(err)
		}

		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		client.stream = s

		client.testMap = make(map[string]string)
		client.rw[s] = rw

		// client := &Client{
		// 	testMap: make(map[string]string),
		// 	rw:      rw,
		// }

		// testMap := make(map[string]string)

		go client.readExampleData(s)
		// go client.writeExampleData(s)

		// fmt.Println("LET'S CHECK OUT THOSE STREAMS")
		// fmt.Println("%+v\n", sampleClient.streams)

		// Create a thread to read and write data.
		// go writeData(rw)
		// go readData(rw)

		// Hang forever.
		select {}
	}
}

func dummyDebasedMetaData() DebasedMetadata {
	dummyData := &DebasedMetadata{
		Accounts: map[string]AccountInfo{
			"cantSeeMe": AccountInfo{
				LiquidBalance:   1,
				IlliquidBalance: 2,
				Permissions: map[string]UserPermission{
					"theOvertaker": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"thePebble": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"bruceBanner": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"robForwardlund": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
				},
			},
			"theOvertaker": AccountInfo{
				LiquidBalance:   3,
				IlliquidBalance: 4,
				Permissions: map[string]UserPermission{
					"cantSeeMe": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"thePebble": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"bruceBanner": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"robForwardlund": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
				},
			},
			"thePebble": AccountInfo{
				LiquidBalance:   5,
				IlliquidBalance: 6,
				Permissions: map[string]UserPermission{
					"cantSeeMe": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"theOverTaker": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"bruceBanner": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"robForwardlund": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
				},
			},
			"bruceBanner": AccountInfo{
				LiquidBalance:   7,
				IlliquidBalance: 8,
				Permissions: map[string]UserPermission{
					"theOvertaker": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"cantSeeMe": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"thePebble": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"robForwardlund": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
				},
			},
			"robForwardlund": AccountInfo{
				LiquidBalance:   9,
				IlliquidBalance: 10,
				Permissions: map[string]UserPermission{
					"theOvertaker": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"cantSeeMe": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"thePebble": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"bruceBanner": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
				},
			},
		},
		Tables: map[string]TableInfo{
			"first": TableInfo{
				Fields: []string{
					"first",
					"second",
					"third",
				},
				Cells: [][]CellLocation{
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(1),
							Position:        big.NewInt(1),
							PostionInRecord: big.NewInt(1),
						},
						CellLocation{
							BlockNumber:     big.NewInt(1),
							Position:        big.NewInt(2),
							PostionInRecord: big.NewInt(2),
						},
						CellLocation{
							BlockNumber:     big.NewInt(1),
							Position:        big.NewInt(3),
							PostionInRecord: big.NewInt(3),
						},
					},
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(2),
							Position:        big.NewInt(1),
							PostionInRecord: big.NewInt(1),
						},
						CellLocation{
							BlockNumber:     big.NewInt(2),
							Position:        big.NewInt(2),
							PostionInRecord: big.NewInt(2),
						},
						CellLocation{
							BlockNumber:     big.NewInt(2),
							Position:        big.NewInt(3),
							PostionInRecord: big.NewInt(3),
						},
					},
				},
				Writes: []RecordLocation{
					RecordLocation{
						BlockNumber: *big.NewInt(1),
						Position:    *big.NewInt(1),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(2),
						Position:    *big.NewInt(2),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(3),
						Position:    *big.NewInt(3),
					},
				},
			},
			"photos": TableInfo{
				Fields: []string{
					"jpeg",
					"gif",
					"jif",
				},
				Cells: [][]CellLocation{
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(5),
							Position:        big.NewInt(9),
							PostionInRecord: big.NewInt(23),
						},
						CellLocation{
							BlockNumber:     big.NewInt(5),
							Position:        big.NewInt(8),
							PostionInRecord: big.NewInt(2),
						},
						CellLocation{
							BlockNumber:     big.NewInt(4),
							Position:        big.NewInt(3),
							PostionInRecord: big.NewInt(37),
						},
					},
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(7),
							Position:        big.NewInt(67),
							PostionInRecord: big.NewInt(14),
						},
						CellLocation{
							BlockNumber:     big.NewInt(24),
							Position:        big.NewInt(32),
							PostionInRecord: big.NewInt(3),
						},
						CellLocation{
							BlockNumber:     big.NewInt(4),
							Position:        big.NewInt(36),
							PostionInRecord: big.NewInt(7),
						},
					},
				},
				Writes: []RecordLocation{
					RecordLocation{
						BlockNumber: *big.NewInt(2),
						Position:    *big.NewInt(4),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(4),
						Position:    *big.NewInt(7),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(4),
						Position:    *big.NewInt(9),
					},
				},
			},
			"champions": TableInfo{
				Fields: []string{
					"cantSeeMe",
					"theOvertaker",
					"TeddyChu",
				},
				Cells: [][]CellLocation{
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(2),
							Position:        big.NewInt(6),
							PostionInRecord: big.NewInt(1000),
						},
						CellLocation{
							BlockNumber:     big.NewInt(3),
							Position:        big.NewInt(17),
							PostionInRecord: big.NewInt(121),
						},
						CellLocation{
							BlockNumber:     big.NewInt(91),
							Position:        big.NewInt(13),
							PostionInRecord: big.NewInt(29),
						},
					},
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(4),
							Position:        big.NewInt(7),
							PostionInRecord: big.NewInt(23),
						},
						CellLocation{
							BlockNumber:     big.NewInt(31),
							Position:        big.NewInt(82),
							PostionInRecord: big.NewInt(43),
						},
						CellLocation{
							BlockNumber:     big.NewInt(12),
							Position:        big.NewInt(6),
							PostionInRecord: big.NewInt(9),
						},
					},
				},
				Writes: []RecordLocation{
					RecordLocation{
						BlockNumber: *big.NewInt(9),
						Position:    *big.NewInt(21),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(12),
						Position:    *big.NewInt(91),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(32),
						Position:    *big.NewInt(56),
					},
				},
			},
		},
	}

	return *dummyData
}

func main() {
	fmt.Println("START OF MAIN FUNCTION")
	sourcePort := flag.Int("sp", 0, "Source port number")
	dest := flag.String("d", "", "Destination multiaddr string")
	help := flag.Bool("help", false, "Display help")
	debug := flag.Bool("debug", false, "Debug generates the same node ID on every execution")
	// testMap := make(map[string]string)
	sampleClient := &Client{}

	flag.Parse()

	if *help {
		// fmt.Printf("This program demonstrates a simple p2p chat application using libp2p\n\n")
		fmt.Println("Usage: Run './chat-exec -sp <SOURCE_PORT>' where <SOURCE_PORT> can be any port number.")
		fmt.Println("Now run './chat-exec -d <MULTIADDR>' where <MULTIADDR> is multiaddress of previous listener host.")

		os.Exit(0)
	}

	// If debug is enabled, use a constant random source to generate the peer ID. Only useful for debugging,
	// off by default. Otherwise, it uses rand.Reader.
	var r io.Reader
	if *debug {
		// Use the port number as the randomness source.
		// This will always generate the same host ID on multiple executions, if the same port number is used.
		// Never do this in production code.
		r = mrand.New(mrand.NewSource(int64(*sourcePort)))
	} else {
		r = rand.Reader
	}

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)

	sampleClient.host = host
	// sampleClient.streams = make(map[peer.ID]net.Stream)
	// sampleClient.rw = make(map[peer.ID]*bufio.ReadWriter)
	sampleClient.streams = make(map[net.Stream]net.Stream)
	sampleClient.rw = make(map[net.Stream]*bufio.ReadWriter)

	// fmt.Println("HOST.ID().PRETTY()")
	// fmt.Println("HOST.ID().PRETTY()")
	// fmt.Printf("%+v\n", host.ID().Pretty())
	// fmt.Println("HOST.ID().PRETTY()")
	// fmt.Println("HOST.ID().PRETTY()")

	if err != nil {
		panic(err)
	}

	go BeginStream(sampleClient, dest, host)
	nodeDebasedSystem := DebasedSystem{PrivateKeysFromSession: make([]*ecdsa.PrivateKey, 0),
		CurrentBlockHeight:     *big.NewInt(0),
		DuringConsensus:        false,
		Metadata:               &DebasedMetadata{Accounts: make(map[string]AccountInfo), Tables: make(map[string]TableInfo)},
		CurrentBids:            make([]*BlockGenerationBid, 0),
		UnconfirmedBlock:       nil,
		CurrentBets:            make([]*Bet, 0),
		PendingBetPayouts:      make([]*Transfer, 0),
		PendingTransactions:    &Transactions{},
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

			data, err := json.Marshal(payment)

			if err != nil {
				panic(err)
			}

			x := JSONWrapper{
				Type:     "Transfer",
				Contents: data,
			}
			err = x.Sign(nodeDebasedSystem.PrivateKeysFromSession[toAcctIndex])
			if err != nil {
				panic(err)
			}
			fmt.Println("here")
			fmt.Printf("x: %+v\n", x)

			test := x.VerifySignature()

			if !test {
				panic(errors.New("Signature not valid"))
			}

			// fmt.Println("after")

			data, err = json.Marshal(x)

			if err != nil {
				panic(err)
			}

			sampleClient.writeExampleData(data, sampleClient.stream)

		}
		// genBlock currentAcctIndexInCurrentAccts
		if args[0] == "genBlock" && len(args) == 2 {
			// fmt.Println(nodeDebasedSystem.GenerateBlock())
			var newBlock = nodeDebasedSystem.GenerateBlock()
			nodeDebasedSystem.CurrentBlockHeight = newBlock.BlockHeight
			nodeDebasedSystem.Metadata = newBlock.UpdatedMD
			b, err := json.Marshal(nodeDebasedSystem.Metadata)
			if err != nil {
				fmt.Println("error:", err)
			}
			wrapper := &JSONWrapper{Type:"GeneratedBlock", Contents:b}
			toAcctIndex, err := strconv.ParseInt(args[1], 10, 64)
			wrapper.Sign(nodeDebasedSystem.PrivateKeysFromSession[toAcctIndex])
			fmt.Println(wrapper)
      fmt.Println("Signature is valid:", wrapper.VerifySignature())
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
