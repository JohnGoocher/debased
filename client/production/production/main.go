package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/ecdsa"
//	"crypto/elliptic"
	"encoding/json"
	"math/big"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	fmt.Println("START OF MAIN FUNCTION")
	sourcePort := flag.Int("sp", 0, "Source port number")
	dest := flag.String("d", "", "Destination multiaddr string")
	help := flag.Bool("help", false, "Display help")
	debug := flag.Bool("debug", false, "Debug generates the same node ID on every execution")
	sampleClient := &Client{context: context.Background()}
	sampleClient.testMap = make(map[string]string)

	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p chat application using libp2p\n\n")
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

	fmt.Printf("sourcePort: %+v\n", *sourcePort)
	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		sampleClient.context,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)

	sampleClient.host = host
	// sampleClient.streams = make(map[peer.ID]net.Stream)
	// sampleClient.rw = make(map[peer.ID]*bufio.ReadWriter)

	// sampleClient.Streams = make(map[string]interface{})
	sampleClient.Streams = make(map[string]net.Stream)
	// sampleClient.buildStreams = make(map[string]bool)
	sampleClient.buildStreams = make(map[string]stream)

	sampleClient.rw = make(map[net.Stream]*bufio.ReadWriter)
	sampleClient.streamPorts = make(map[string]string)

	fmt.Println("HOST.ID().PRETTY()")
	fmt.Println("HOST.ID().PRETTY()")
	fmt.Printf("%+v\n", host.ID().Pretty())
	fmt.Println("HOST.ID().PRETTY()")
	fmt.Println("HOST.ID().PRETTY()")

	if err != nil {
		panic(err)
	}

	//pos SHIT here
	//localDebasedSystem should be a pointer
	msgHolder := &[][]byte{}
	localDebasedSystem := &DebasedSystem{PrivateKeysFromSession: make([]*ecdsa.PrivateKey, 0),
		                                 CurrentBlockHeight: *big.NewInt(0),
																		 DuringConsensus: false,
																		 Metadata: &DebasedMetadata{Accounts: make(map[string]AccountInfo), Tables: make(map[string]TableInfo)},
																		 CurrentBids: make([]*BlockGenerationBid, 0),
																		 UnconfirmedBlock: nil,
																		 CurrentBets: make([]*Bet, 0),
																		 PendingBetPayouts: make([]*Transfer, 0),
																		 PendingTransactions: &Transactions{},
																		 HoldingPenTransactions: &Transactions{},
																		 P2PMsgsToSend: msgHolder,
	                                   }
	sampleClient.localDebasedSystem = localDebasedSystem
	sampleClient.msgsToBeSent = msgHolder

	testAcctPrivKeyFrom, acctNumFrom := createAcct()
	_, acctNumTo := createAcct()
	testTransfer := Transfer{ToAcctID: acctNumTo, Ammount: 365.50, FromAcctID: acctNumFrom}
	fmt.Println("created Test Transfer")
	fmt.Println(testTransfer)
	marshalledTestTransfer, err := json.Marshal(testTransfer)
	if err != nil {
		fmt.Println("error Marshalling testPosWrapper")
		panic(err)
	}
	testPosWrapper := POSWrapper{Type: "Transfer", Contents: marshalledTestTransfer}
	testPosWrapper.Sign(testAcctPrivKeyFrom)
	fmt.Println("It signed!")
	if testPosWrapper.VerifySignature() == true {
		fmt.Println("HURRAY IT SIGNED CORRECTLY")
	}
	testByteSlice, err := json.Marshal(testPosWrapper)
	if err != nil {
		fmt.Println("error Marshalling testPosWrapper")
		panic(err)
	}
	*localDebasedSystem.P2PMsgsToSend = append(*localDebasedSystem.P2PMsgsToSend, testByteSlice)

	if *dest == "" {
		// Set a function as stream handler.
		// This function is called when a peer connects, and starts a stream with this protocol.
		// Only applies on the receiving side.
		host.SetStreamHandler("/chat/1.0.0", sampleClient.handleStream)
		// id := host.ID()
		idString := host.ID().Pretty()
		// sampleClient.buildStreams[idString] = id
		// sampleClient.buildStreams[idString] = stream{connected: true}
		// sampleClient.buildStreams[idString] = firstStream{
		// 	connected: true,
		// 	origin:    true,
		// }

		fmt.Printf("sampleClient.buildStreams: %+v\n", sampleClient.buildStreams)

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

		sampleClient.buildStreams[idString] = stream{
			connected: true,
			port:      port,
		}

		fmt.Printf("sampleClient.buildStreams: %+v\n", sampleClient.buildStreams)

		sampleClient.streamPorts[idString] = port

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

		fmt.Printf("dest: %+v\n", *dest)

		// Turn the destination into a multiaddr.
		maddr, err := multiaddr.NewMultiaddr(*dest)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("dest: %+v\n", *dest)
		fmt.Printf("maddr: %+v\n", maddr)

		// Extract the peer ID from the multiaddr.
		info, err := peerstore.InfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("dest: %+v\n", *dest)
		fmt.Printf("maddr: %+v\n", maddr)
		fmt.Printf("info: %+v\n", info)

		host.SetStreamHandler("/chat/1.0.0", sampleClient.handleStream)

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

		// fmt.Printf("Run './chat-exec -d /ip4/127.0.0.1/tcp//p2p/%s' on another console.\n", host.ID().Pretty())

		// Add the destination's peer multiaddress in the peerstore.
		// This will be used during connection and stream creation by libp2p.
		host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		nodeAddress := strings.Split(*dest, "/")
		nodeID := nodeAddress[len(nodeAddress)-1]

		// id := host.ID()
		idString := host.ID().Pretty()
		// sampleClient.buildStreams[idString] = id
		sampleClient.streamPorts[idString] = port
		// sampleClient.buildStreams[idString] = true
		sampleClient.buildStreams[idString] = stream{connected: true, port: port}
		// sampleClient.buildStreams[nodeID] = true
		reassign := sampleClient.buildStreams[nodeID]
		reassign.connected = true
		sampleClient.buildStreams[nodeID] = reassign

		// x := info.ID.Pretty()

		// y, err := peer.IDB58Decode(x)

		// if err != nil {
		// 	fmt.Println("ARE YOU THE NEW PANIC???")
		// 	panic(err)
		// }

		// fmt.Printf("x: %+v\n", x)
		// fmt.Printf("y: %+v\n", y)

		fmt.Printf("info.ID.String(): %+v\n", info.ID.String())
		fmt.Printf("info.ID: %+v\n", info.ID)

		// MAYBE TRY USING A SETTTTTTT ORRR  A A A A A A STIRNGGGGG
		fmt.Printf("sampleClient.buildStreams: %+v\n", sampleClient.buildStreams)

		// Start a stream with the destination.
		// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
		fmt.Printf("info.ID: %+v\n", info.ID)

		s, err := host.NewStream(sampleClient.context, info.ID, "/chat/1.0.0")
		if err != nil {
			panic(err)
		}

		// panic(errors.New("asdfaf"))

		sampleClient.testMap = make(map[string]string)
		sampleClient.POSWrapperSlice = make([]*POSWrapper, 0)
		sampleClient.rw[s] = bufio.NewReadWriter(bufio.NewReaderSize(s, 5000), bufio.NewWriterSize(s, 5000))

		// sampleClient := &Client{
		// 	testMap: make(map[string]string),
		// 	rw:      rw,
		// }

		// testMap := make(map[string]string)

		go sampleClient.readExampleData(s)
		go sampleClient.writeExampleData(s)
		go sampleClient.writePOSWrappedData(s)

		fmt.Println("LET'S CHECK OUT THOSE STREAMS")
		fmt.Println("%+v\n", sampleClient.Streams)

		// Create a thread to read and write data.
		// go writeData(rw)
		// go readData(rw)


		// Hang forever.
		select {}
	}
}
