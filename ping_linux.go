package main

import (
	"flag"
	"fmt"
	"github.com/nperez-messagebird/sctp"
	"log"
	"strconv"
	"time"
)

func servePing(conn *sctp.SCTPConn) {
	for {
		msg := make([]byte, 256)
		n, err := conn.Read(msg)
		if err != nil {
			log.Printf("error reading message: %s", err)
			_ = conn.Close()
			return
		}
		n, err = conn.Write(msg[:n])
		if err != nil {
			log.Printf("error writing message: %s", err)
			_ = conn.Close()
			return
		}
	}
}

var sockInit = *sctp.NewDefaultInitMsg()

func pingServer(addr *sctp.SCTPAddr) {
	l, err := sctp.NewSCTPListener(addr, sockInit, sctp.OneToOne)
	if err != nil {
		log.Fatalf("error listening sctp socket: %s", err)
	}
	for {
		conn, err := l.AcceptSCTP()
		if err != nil {
			log.Fatalf("failed to accept: %s", err)
		}
		go servePing(conn)
	}
}

func pingClient(bind *sctp.SCTPAddr, peer *sctp.SCTPAddr) {
	var buf = make([]byte, 0, 256)
	c, err := sctp.NewSCTPConnection(bind, peer, sockInit, sctp.OneToOne)
	if err != nil {
		log.Fatalf("error connecting to server: %s", err)
	}
	for i := int64(0);; i++ {
		buf = buf[0:0]
		buf = strconv.AppendInt(buf, i, 10)
		ts := time.Now()
		n, err := c.Write(buf)
		if err != nil {
			log.Fatalf("error writing to server: %s", err)
		}
		fmt.Print("c -> ")
		buf = buf[:]
		n, err = c.Read(buf)
		if err != nil {
			log.Fatalf("error writing to server: %s", err)
		}
		fmt.Printf("s (%0.3fms): %s\n", float64(time.Now().Sub(ts)) / float64(time.Millisecond), string(buf[:n]))
		time.Sleep(time.Second)
	}
}

func main() {
	var (
		bindAddrStr = flag.String("l", "0.0.0.0:1234", "local address")
		peerAddrStr = flag.String("r", "", "remote address")
		bindAddr, peerAddr *sctp.SCTPAddr
		err error
	)

	flag.Parse()

	if bindAddr, err = ResolveSCTPAddr("sctp", *bindAddrStr); err != nil {
		log.Fatalf("can't parse bind address: %s", err.Error())
	}

	if *peerAddrStr != "" {
		if peerAddr, err = ResolveSCTPAddr("sctp", *peerAddrStr); err != nil {
			log.Fatalf("can't parse remote peer address: %s", err.Error())
		}
		pingClient(bindAddr, peerAddr)
	} else {
		pingServer(bindAddr)
	}

}