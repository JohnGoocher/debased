/*
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 Juan Batiz-Benet
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 * This program demonstrate a simple chat application using p2p communication.
 *
 */

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

	"github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
)

type Client struct {
	testMap map[string]string

	// rw map[peer.ID]*bufio.ReadWriter
	rw map[net.Stream]*bufio.ReadWriter

	// streams map[peer.ID]net.Stream
	streams map[net.Stream]net.Stream
	host    host.Host
}

// type Node struct {
// }

// type Chat interface {

// }

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
	go c.writeExampleData(s)

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
	if len(str) > 0 {
		if err := json.Unmarshal(str, &c.testMap); err != nil {
			// fmt.Println("OR IS IT THIS PANIC")
			panic(err)
		}
		fmt.Printf("%+v\n", c.testMap)
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
		if len(str) > 0 {
			if err := json.Unmarshal(str, &c.testMap); err != nil {
				// fmt.Println("OR IS IT THIS PANIC")
				panic(err)
			}
			fmt.Printf("%+v\n", c.testMap)
			// fmt.Println("END OF ELSE")
			// fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		} else {
			fmt.Println("Receieved value is 0")
			return
		}
		// fmt.Println("END OF A LOOP")
	}
}

func (c *Client) writeExampleData(s net.Stream) {
	// testMap := make(map[string]string)
	stdReader := bufio.NewReader(os.Stdin)
	count := 0
	for {
		// fmt.Println("count: %i", count)
		fmt.Print("before> ")
		data, err := stdReader.ReadString('\n')
		// fmt.Printf("WRITING: %+v\n", s)
		// fmt.Println("AFTER READSTRING")
		// fmt.Println("data: " + data)

		// fmt.Println("WHAT")
		// fmt.Printf("c.testMap before: %+v\n", c.testMap)
		// fmt.Printf("data: %+v\n", data)
		c.testMap[data] = data
		// fmt.Printf("c.testMap after: %+v\n", c.testMap)
		// fmt.Println("AFTER TESTMAP")
		sendData, err := json.Marshal(c.testMap)
		if err != nil {
			fmt.Println("JSON.MARSHALL PANIC")
			panic(err)
		}
		// fmt.Println("AFTER SEND DATA")
		// fmt.Println("before write")
		// fmt.Printf("sendData: %+v\n", string(sendData))
		// fmt.Println("STRRRRRREEEAMMMMSSSSSSS")
		// fmt.Printf("%+v\n", c.streams)

		for _, writer := range c.rw {
			writer.Write(sendData)
			writer.Flush()
		}

		// old way
		// c.rw[s].Write(sendData)
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

	if *dest == "" {
		// Set a function as stream handler.
		// This function is called when a peer connects, and starts a stream with this protocol.
		// Only applies on the receiving side.
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

		host.SetStreamHandler("/chat/1.0.0", sampleClient.handleStream)
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

		sampleClient.testMap = make(map[string]string)
		sampleClient.rw[s] = rw

		// sampleClient := &Client{
		// 	testMap: make(map[string]string),
		// 	rw:      rw,
		// }

		// testMap := make(map[string]string)

		go sampleClient.readExampleData(s)
		go sampleClient.writeExampleData(s)

		// fmt.Println("LET'S CHECK OUT THOSE STREAMS")
		// fmt.Println("%+v\n", sampleClient.streams)

		// Create a thread to read and write data.
		// go writeData(rw)
		// go readData(rw)

		// Hang forever.
		select {}
	}
}
