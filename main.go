package main

// https://www.glendimplex.se/media/15650/nobo-hub-api-v-1-1-integration-for-advanced-users.pdf


//unc main() {
	// broadcast udp 10000
	// multicast 239.0.1.187 port 10001
//}

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	commandSetVersion = "1.1"
	timeFormat        = "20060102150405"
)

func main() {
	hubIP, hubPartialSerial := autoDiscoverMulticast()

	if len(os.Args) <= 1 {
		log.Fatal(`usage: gonobo <hubSerial>
hint: the output so far should show the first 9 digits of the serial,
      and you must provide all 12 digits`)
	}
	hubSerial := os.Args[1]
	log.Print("Looking for hub with serial ", hubSerial)

	if !strings.HasPrefix(hubSerial, hubPartialSerial) {
		log.Fatal("discovered hub with wrong serial (only one hub supported)")
	}

	conn := dial(hubIP)
	doHandshake(conn, hubSerial)
	doCommandG00(conn)
}

// autoDiscoverMulticast returns the IP address of the first Nobø hub
// that sends a multicast packet on the network.
// Multiple hubs are not supported yet (but can be added in the future).
func autoDiscoverMulticast() (net.IP, string)  {
	log.Print("Auto discovering Nobø hub by listening to multicast...")
	maxDatagramSize := 8192
	addr, err := net.ResolveUDPAddr("udp", "239.0.1.187:10001")
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	_ = l.SetReadBuffer(maxDatagramSize)
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		if string(b[0:11]) == "__NOBOHUB__" {
			partialSerial := string(b[11:n])
			log.Print("Found hub at ", src.IP, " with serial starting with ", partialSerial)
			return src.IP, partialSerial
		}
	}
}

func dial(ip net.IP) net.Conn {
	// Connect
	addr := net.TCPAddr{
		IP:   ip,
		Port: 27779,
	}
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func doHandshake(conn net.Conn, hubSerial string) {
	// Hello message to hub
	timeStr := time.Now().UTC().Format(timeFormat)
	helloMsg := fmt.Sprintf("HELLO 1.1 %s %s\r", hubSerial, timeStr)
	_, err := conn.Write([]byte(helloMsg))
	if err != nil {
		log.Print("cannot write to tcp: ", err)
	}
	log.Print("To hub: ", helloMsg)

	// Hello message from hub
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		log.Print("err reading TCP: ", err)
	}
	helloReply := string(buf[:n])
	expectedHelloReply := "HELLO 1.1\r"
	if helloReply != expectedHelloReply {
		log.Print("unexpected message during handshake: ", helloReply)
	} else {
		log.Print("From hub: ", helloReply)
	}

	// Handshake message to hub
	handshakeMsg := "HANDSHAKE\r"
	_, err = conn.Write([]byte(handshakeMsg))
	if err != nil {
		log.Print("cannot write to tcp: ", err)
	}
	log.Print("To hub: ", handshakeMsg)

	// Hello message from hub
	n, err = conn.Read(buf)
	if err != nil {
		log.Print("err reading TCP: ", err)
	}
	handshakeReply := string(buf[:n])
	expectedHandshakeReply := "HANDSHAKE\r"
	if handshakeReply != expectedHandshakeReply {
		log.Print("unexpected message during handshake: ", helloReply)
	} else {
		log.Print("From hub: ", handshakeReply)
	}
}

func doCommandG00(conn net.Conn) {
	// G00 command to download all relevant information from the hub
	g00Msg := "G00\r"
	_, err := conn.Write([]byte(g00Msg))
	if err != nil {
		log.Print("cannot write to tcp: ", err)
	}
	log.Print("To hub: ", g00Msg)

	buf := make([]byte, 8196)
	for {
		// Get info from hub
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal("err reading TCP: ", err)
		}
		infoReply := string(buf[:n])
		parts := strings.Split(strings.Trim(infoReply, "\r"), " ")
		command := parts[0]
		fmt.Println(infoReply)

		if command == "H05" {
			//break
		}
	}
}