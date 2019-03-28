package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
)

type streamWrapper struct {
	net.Stream
	rw io.ReadWriter
}

type firstStream struct {
	connected bool
	origin    bool
}

type Client struct {
	testMap map[string]string

	// rw map[peer.ID]*bufio.ReadWriter
	rw map[net.Stream]*bufio.ReadWriter

	// streams map[net.Stream]net.Stream
	// Streams map[string]net.Stream
	Streams map[string]net.Stream
	// Streams map[string]interface{}

	// buildStreams map[string]bool
	buildStreams map[string]stream
	streamPorts  map[string]string
	// buildStreams map[string]firstStream

	host host.Host

	context context.Context
}

type stream struct {
	connected bool
	port      string
}

type jsonWrapper struct {
	ObjectType string
	// Object     interface{}
	Object []byte
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
	fmt.Println("AM I THE PEER ID THINGY")
	fmt.Printf("%+v\n", info.ID)
	fmt.Println("AM I THE PEER ID THINGY")
	return info.ID
}

func (c *Client) handleStream(s net.Stream) {
	log.Println("Got a new stream!")
	fmt.Printf("s.Conn().Stat(): %+v\n", s.Conn().Stat())
	fmt.Printf("s.Protocol(): %+v\n", s.Protocol())
	fmt.Printf("s.Conn().LocalPeer: %+v\n", s.Conn().LocalPeer)
	fmt.Printf("s.Conn().LocalPrivateKey: %+v\n", s.Conn().LocalPrivateKey)
	fmt.Printf("s.Conn().LocalMultiaddr: %+v\n", s.Conn().LocalMultiaddr)
	fmt.Printf("s.Conn().RemotePeer: %+v\n", s.Conn().RemotePeer)
	fmt.Printf("s.Conn().RemotePublicKey: %+v\n", s.Conn().RemotePublicKey)
	fmt.Printf("s.Conn().RemoteMultiaddr: %+v\n", s.Conn().RemoteMultiaddr)
	fmt.Printf("s.Stat(): %+v\n", s.Stat())
	fmt.Printf("s.Conn(): %+v\n", s.Conn())
	// fmt.Printf("s.Protocol(): %+v\n", s.Protocol

	// uuid := xid.New()
	// fmt.Println(uuid.String())
	// temp := map[xid.ID]net.Stream{}

	if val, ok := c.rw[s]; ok {
		fmt.Println("It's there!")
		fmt.Printf("%+v", val)
	}

	// Create a buffer stream for non blocking read and write.
	// c.rw[c.host.ID()] = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	// c.rw[s] = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	fmt.Printf("c.host.ID().Pretty(): %+v\n", c.host.ID().Pretty())
	firstIndexID := strings.LastIndex(fmt.Sprintf("%+v", s.Conn()), "(")
	lastIndexID := strings.LastIndex(fmt.Sprintf("%+v", s.Conn()), ")")
	parsedIncomingStream := fmt.Sprintf("%+v", s.Conn())[firstIndexID+1 : lastIndexID]
	firstIndexPort := strings.LastIndex(fmt.Sprintf("%+v", s.Conn()), "/")
	lastIndexPort := strings.LastIndex(fmt.Sprintf("%+v", s.Conn()), " ")
	parsedIncomingPort := fmt.Sprintf("%+v", s.Conn())[firstIndexPort+1 : lastIndexPort]
	fmt.Printf("parsedIncomingPort: %+v\n", parsedIncomingPort)

	c.rw[s] = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	// c.buildStreams[parsedIncomingStream] = true
	// reassign := c.buildStreams[parsedIncomingStream]
	// reassign.connected = true
	c.buildStreams[parsedIncomingStream] = stream{
		connected: true,
		port:      parsedIncomingPort,
	}

	fmt.Printf("BEFORE c.streams: %+v\n", c.Streams)

	// if val, ok := dict["foo"]; ok {
	// 	//do something here
	// }

	// fmt.Println(s)
	// out, err := json.Marshal(s)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(out))

	// c.streams[s] = s

	c.Streams[fmt.Sprintf("%+v", s)] = s

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

	fmt.Printf("AFTER c.streams: %+v\n", c.Streams)
	fmt.Printf("s: %+v\n", s)

	go c.readExampleData(s)
	go c.writeExampleData(s)

	fmt.Printf("%+v\n", c.Streams)
	fmt.Printf("%+v\n", s)

	// panic(errors.New("asldfjalfd"))

	c.writeStreams(s)

	// panic(errors.New("STAHP"))
	// go readData(rw)
	// go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func (c *Client) writeStreams(s net.Stream) {

	fmt.Println("writeStreams before json.Marshal")
	fmt.Printf("c.buildStreams: %+v\n", c.buildStreams)
	fmt.Printf("c.buildStreams: %+v\n", c.buildStreams)

	buildStreamBytes, err := json.Marshal(c.buildStreams)
	if err != nil {
		fmt.Println("JSON.MARSHALL PANIC")
		panic(err)
	}

	wrapper := &jsonWrapper{
		Object:     buildStreamBytes,
		ObjectType: "buildStreams",
	}

	wrapperBytes, err := json.Marshal(wrapper)
	if err != nil {
		fmt.Println("JSON.MARSHALL PANIC")
		panic(err)
	}

	// sendData = append(sendData, '\n')

	fmt.Println("WELL IS THIS WORKING")
	fmt.Printf("%+v\n", string(wrapperBytes))
	fmt.Printf("%+v\n", wrapperBytes)

	fmt.Println("STRRRRRREEEAMMMMSSSSSSS")
	fmt.Printf("%+v\n", c.Streams)

	c.rw[s].Write(wrapperBytes)
	c.rw[s].Flush()

	fmt.Println("POST FLUSH")

	// temp := fmt.Sprintf("%+v", c.streams)

	// fmt.Println("testing Sprintf")
	// fmt.Println(temp)

	// c.rw[s].Write([]byte(temp))
	// c.rw[s].Flush()
}

// func (c *Client) readStreams(s net.Stream) {
// 	// str, err := c.rw[s].ReadSlice('}')
// 	var str []byte
// 	err := json.NewDecoder(c.rw[s]).Decode(&str)
// 	// str, err := c.rw[s].Dee('}')
// 	fmt.Println("READING THAT STR")
// 	fmt.Println(string(str))
// 	fmt.Printf("READING: %+v\n", s)
// 	fmt.Println("STTTTRRRREEEEAAAAAMMMMS")
// 	fmt.Printf("%+v\n", c.streams)

// 	// fmt.Println("AFTER READSLICE")
// 	if err != nil {
// 		// fmt.Println("READSLICE PANIC")
// 		panic(err)
// 	}
// 	// fmt.Println("after readSlice")
// 	fmt.Printf("readslice: %+v\n", string(str))
// 	if len(str) > 0 {
// 		if err := json.Unmarshal(str, &c.testMap); err != nil {
// 			// fmt.Println("OR IS IT THIS PANIC")
// 			panic(err)
// 		}
// 		fmt.Printf("%+v\n", c.testMap)
// 		// fmt.Println("END OF ELSE")
// 		// fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
// 	} else {
// 		fmt.Println("Receieved value is 0")
// 		return
// 	}
// }

func (c *Client) readExampleData(s net.Stream) {
	// var testMap map[string]string
	for {
		// fmt.Println("FOR LOOP STARTED")
		// str, err := c.rw[c.host.ID()].ReadSlice('}')

		// IT'S STOPPING AT THE FIRST } THAT IT SEES, BUT THERE IS
		// MORE TO THE SLICE THAT IS INCOMING
		fmt.Println("BEFORE STR")
		// str, err := c.rw[s].ReadSlice('}')
		// str, err := c.rw[s].ReadSlice('\x00')
		wrapperBytes, err := c.rw[s].ReadSlice('}')
		// str, err := c.rw[s].ReadSlice('\n')

		fmt.Println("AFTER READSLICE")
		if err != nil {
			fmt.Println("READSLICE PANIC")
			panic(err)
		}

		// str = str[:len(str)-1]

		// if _, err := c.rw[s].ReadSlice('\x00'); err != nil {
		// 	fmt.Println("READSLICE PANIC")
		// 	panic(err)
		// }
		fmt.Println("AFTER STR")
		fmt.Printf("%+v\n", wrapperBytes)

		fmt.Printf("READING: %+v\n", s)
		fmt.Println("STTTTRRRREEEEAAAAAMMMMS")
		fmt.Printf("%+v\n", c.Streams)

		// fmt.Println("after readSlice")
		// str = append(str, '"')
		// str = append(str, '}')
		fmt.Printf("readslice: %+v\n", string(wrapperBytes))
		// str = append(str, '}')

		// panic(errors.New("ASDLKAJSDFLKASJDF;LAKDSFAFD"))

		// var temp map[string]interface{}

		// incomingBuildStreams := map[string]bool{}

		fmt.Println("MORE MORE MORE")
		fmt.Println("MORE MORE MORE")
		fmt.Println("MORE MORE MORE")

		if len(wrapperBytes) > 0 {

			incomingWrapper := &jsonWrapper{}

			fmt.Printf("string(wrapperBytes): %+v\n", string(wrapperBytes))

			if err := json.Unmarshal(wrapperBytes, &incomingWrapper); err != nil {
				fmt.Println("Cannot unmarshal incomingWrapper")
				panic(err)
			}

			fmt.Printf("incomingWrapper: %+v\n", incomingWrapper)
			fmt.Printf("string(incomingWrapper.Object): %+v\n", string(incomingWrapper.Object))
			fmt.Printf("incomingWrapper.ObjectType: %+v\n", incomingWrapper.ObjectType)

			switch incomingWrapper.ObjectType {
			case "buildStreams":
				c.buildNewStreams(*incomingWrapper)
			case "buildTestMaps":
				c.buildTestMaps(*incomingWrapper)
			}

			// incomingObject := make(map[string]bool)

			// if err := json.Unmarshal(incomingWrapper.Object, &incomingObject); err != nil {
			// 	fmt.Println("Cannot unmarshal incomingObject")
			// 	panic(err)
			// }

			// fmt.Printf("incomingObject: %+v\n", incomingObject)

			// if len(incomingObject) > 0 {
			// 	fmt.Println("INCOMING STREAMS")
			// 	fmt.Printf("\nc.bs: %+v\n", c.buildStreams)
			// 	fmt.Printf("ibs: %+v\n", incomingObject)
			// 	for key, _ := range incomingObject {
			// 		// c.buildStreams[key] = value
			// 		_, ok := c.buildStreams[key]
			// 		c.buildStreams[key] = ok
			// 	}
			// 	fmt.Printf("c.bs: %+v\n\n", c.buildStreams)
			// 	c.buildNewStreams()
			// }

			// if err := json.Unmarshal(str, &incomingBuildStreams); err == nil {
			// 	// if err := json.Unmarshal(str, &temp); err != nil {
			// 	// fmt.Println("json unmarshal c.streams error")
			// 	fmt.Println("incomingBuildStreams works")
			// 	// panic(err)
			// 	// continue
			// } else if err := json.Unmarshal(str, &c.testMap); err == nil {
			// 	// fmt.Println("json unmarshal c.testMap error")
			// 	// panic(err)
			// 	fmt.Println("c.testMap works")
			// 	// continue
			// } else {
			// 	panic(errors.New("Couldn't unmarshal in readData"))
			// }

			// fmt.Println("temp")
			// fmt.Printf("%+v\n", temp)
			// fmt.Println("c.testMap")
			// fmt.Printf("%+v\n", c.testMap)
			// fmt.Println("c.streams")
			// fmt.Printf("%+v\n", c.Streams)
			// fmt.Println("incomingBuildStreams")
			// fmt.Printf("%+v\n", incomingBuildStreams)

			// fmt.Println("END OF ELSE")
			// fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		} else {
			fmt.Println("Receieved value is 0")
			return
		}
		// fmt.Println("END OF A LOOP")
	}
}

func (c *Client) buildTestMaps(incomingWrapper jsonWrapper) {
	fmt.Println("buildTestMaps")
	fmt.Println("BEFORE")
	fmt.Printf("c.testMap: %+v\n", c.testMap)
	if err := json.Unmarshal(incomingWrapper.Object, &c.testMap); err != nil {
		fmt.Println("Cannot unmarshal incomingWrapper.Object for c.testMap")
		panic(err)
	}
	fmt.Println("AFTER")
	fmt.Printf("c.testMap: %+v\n", c.testMap)
}

func (c *Client) buildNewStreams(incomingWrapper jsonWrapper) {
	// incomingObject := make(map[string]bool)
	incomingObject := make(map[string]stream)

	if err := json.Unmarshal(incomingWrapper.Object, &incomingObject); err != nil {
		fmt.Println("Cannot unmarshal incomingObject")
		panic(err)
	}

	fmt.Printf("incomingObject: %+v\n", incomingObject)

	if len(incomingObject) > 0 {
		fmt.Println("INCOMING STREAMS")
		fmt.Printf("\nc.bs: %+v\n", c.buildStreams)
		fmt.Printf("ibs: %+v\n", incomingObject)
		for key, value := range incomingObject {
			// _, ok := c.buildStreams[key]
			// c.buildStreams[key] = ok
			val, ok := c.buildStreams[key]
			if !ok {
				c.buildStreams[key] = stream{
					connected: false,
					port:      value.port,
				}
			} else {
				if val.port == "" {
					reassign := val
					reassign.port = value.port
					c.buildStreams[key] = reassign
				}
			}
		}
		fmt.Printf("c.bs: %+v\n\n", c.buildStreams)
	}

	fmt.Println("HOW MANY TIMES AM I RUNNING")
	fmt.Println("HOW MANY TIMES AM I RUNNING")
	fmt.Println("HOW MANY TIMES AM I RUNNING")
	fmt.Println("HOW MANY TIMES AM I RUNNING")
	fmt.Println("HOW MANY TIMES AM I RUNNING")
	fmt.Println("HOW MANY TIMES AM I RUNNING")
	fmt.Println("HOW MANY TIMES AM I RUNNING")
	fmt.Println("HOW MANY TIMES AM I RUNNING")
	fmt.Println("HOW MANY TIMES AM I RUNNING")
	fmt.Println("HOW MANY TIMES AM I RUNNING")
	fmt.Println("c.buildNewStreams")
	fmt.Printf("%+v\n", c.buildStreams)
	for key, value := range c.buildStreams {
		hostID := c.host.ID().Pretty()
		fmt.Printf("hostID: %+v\n", hostID)
		fmt.Printf("key: %+v\n", key)
		if hostID == key || value.connected == true {
			fmt.Println("I AM HEREREEE")
			fmt.Println("I AM HEREREEE")
			fmt.Println("I AM HEREREEE")
			fmt.Println("I AM HEREREEE")
			fmt.Printf("key: %+v\n", key)
			fmt.Printf("value: %+v\n", value)
		} else {
			fmt.Println("some value got through")
			fmt.Printf("key: %+v\n", key)
			fmt.Printf("value: %+v\n", value)
			// c.buildStreams[key].connected = true
			reassign := c.buildStreams[key]
			reassign.connected = true
			c.buildStreams[key] = reassign
			// typeCast, err := peer.IDFromString(key)
			// if err != nil {
			// 	fmt.Println("Type cast didn't work")
			// 	panic(err)
			// }
			// fmt.Println("key: " + key)
			// fmt.Printf("c.host.ID(): %+v\n", c.host.ID())
			// fmt.Printf("value: %+v\n", value)
			// fmt.Printf("typeCast: %+v\n", typeCast)
			// panic(errors.New("ASDLKAJSDFLKASJDF;LAKDSFAFD"))
			typeCast, err := peer.IDB58Decode(key)
			fmt.Printf("key: %+v\n", key)
			fmt.Printf("typeCast: %+v\n", typeCast)
			if err != nil {
				fmt.Println("Typecast decode panic")
				// fmt.Println("key: " + key)
				panic(err)
			}
			// ================================================================================================
			// ==============================LET'S COPY AND PASTE==============================================
			// ================================================================================================
			dest := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%s", c.streamPorts[key], key)
			fmt.Println("This node's multiaddresses:")
			for _, la := range c.host.Addrs() {
				fmt.Printf(" - %v\n", la)
			}
			fmt.Println()

			fmt.Printf("dest: %+v\n", dest)

			// Turn the destination into a multiaddr.
			maddr, err := multiaddr.NewMultiaddr(dest)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("dest: %+v\n", dest)
			fmt.Printf("maddr: %+v\n", maddr)

			// Extract the peer ID from the multiaddr.
			info, err := peerstore.InfoFromP2pAddr(maddr)
			if err != nil {
				log.Fatalln(err)
			}

			// fmt.Printf("dest: %+v\n", *dest)
			// fmt.Printf("maddr: %+v\n", maddr)
			fmt.Printf("info: %+v\n", info)

			c.host.SetStreamHandler("/chat/1.0.0", c.handleStream)

			// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
			var port string
			for _, la := range c.host.Network().ListenAddresses() {
				if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
					port = p
					break
				}
			}

			if port == "" {
				panic("was not able to find actual local port")
			}

			fmt.Printf("Run './chat-exec -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, c.host.ID().Pretty())

			// fmt.Printf("Run './chat-exec -d /ip4/127.0.0.1/tcp//p2p/%s' on another console.\n", host.ID().Pretty())

			// Add the destination's peer multiaddress in the peerstore.
			// This will be used during connection and stream creation by libp2p.
			c.host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

			// nodeAddress := strings.Split(*dest, "/")
			// nodeID := nodeAddress[len(nodeAddress)-1]

			// id := host.ID()
			// idString := c.host.ID().Pretty()
			// sampleClient.buildStreams[idString] = id
			// c.buildStreams[idString] = true
			// c.buildStreams[nodeID] = true

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
			fmt.Printf("c.buildStreams: %+v\n", c.buildStreams)

			// Start a stream with the destination.
			// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
			fmt.Printf("info.ID: %+v\n", info.ID)

			s, err := c.host.NewStream(c.context, typeCast, "/chat/1.0.0")
			if err != nil {
				fmt.Println("Cannot dail this peer/node")
				panic(err)
			}
			// ================================================================================================
			// ==============================LET'S COPY AND PASTE==============================================
			// ================================================================================================

			// Create a buffered stream so that read and writes are non blocking.
			c.rw[s] = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
			go c.readExampleData(s)
			go c.writeExampleData(s)
		}
	}
	// panic(errors.New("stahp"))
}

func (c *Client) writeExampleData(s net.Stream) {
	// testMap := make(map[string]string)
	stdReader := bufio.NewReader(os.Stdin)
	count := 0
	for {

		fmt.Println("count: %i", count)
		fmt.Print("before> ")
		data, err := stdReader.ReadString('\n')
		fmt.Printf("WRITING: %+v\n", s)
		// fmt.Println("AFTER READSTRING")
		// fmt.Println("data: " + data)

		fmt.Println("WHAT")
		fmt.Printf("c.testMap before: %+v\n", c.testMap)
		fmt.Printf("data: %+v\n", data)
		c.testMap[data] = data
		fmt.Printf("c.testMap after: %+v\n", c.testMap)
		// fmt.Println("AFTER TESTMAP")
		buildTestMaps, err := json.Marshal(c.testMap)
		if err != nil {
			fmt.Println("JSON.MARSHALL PANIC")
			panic(err)
		}

		wrapper := &jsonWrapper{
			Object:     buildTestMaps,
			ObjectType: "buildTestMaps",
		}

		wrapperBytes, err := json.Marshal(wrapper)
		if err != nil {
			fmt.Println("JSON.MARSHALL PANIC")
			panic(err)
		}

		fmt.Println("WELL IS THIS WORKING")
		fmt.Printf("%+v\n", string(wrapperBytes))
		fmt.Printf("%+v\n", wrapperBytes)

		for _, writer := range c.rw {
			writer.Write(wrapperBytes)
			writer.Flush()
		}

		// old way
		// c.rw[s].Write(buildTestMaps)
		// c.rw[s].Flush()
		// old way

		count++
	}
}

// func readData(rw *bufio.ReadWriter) {
// 	for {
// 		str, _ := rw.ReadString('\n')

// 		if str == "" {
// 			return
// 		}
// 		if str != "\n" {
// 			// Green console colour: 	\x1b[32m
// 			// Reset console colour: 	\x1b[0m
// 			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
// 		}

// 	}
// }

// func writeData(rw *bufio.ReadWriter) {
// 	stdReader := bufio.NewReader(os.Stdin)

// 	for {
// 		fmt.Print("> ")
// 		sendData, err := stdReader.ReadString('\n')

// 		if err != nil {
// 			panic(err)
// 		}

// 		rw.WriteString(fmt.Sprintf("%s\n", sendData))
// 		rw.Flush()
// 	}

// }

func main() {
	fmt.Println("START OF MAIN FUNCTION")
	sourcePort := flag.Int("sp", 0, "Source port number")
	dest := flag.String("d", "", "Destination multiaddr string")
	help := flag.Bool("help", false, "Display help")
	debug := flag.Bool("debug", false, "Debug generates the same node ID on every execution")
	// testMap := make(map[string]string)
	sampleClient := &Client{context: context.Background()}

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
		sampleClient.rw[s] = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// sampleClient := &Client{
		// 	testMap: make(map[string]string),
		// 	rw:      rw,
		// }

		// testMap := make(map[string]string)

		go sampleClient.readExampleData(s)
		go sampleClient.writeExampleData(s)

		fmt.Println("LET'S CHECK OUT THOSE STREAMS")
		fmt.Println("%+v\n", sampleClient.Streams)

		// Create a thread to read and write data.
		// go writeData(rw)
		// go readData(rw)

		// Hang forever.
		select {}
	}
}
