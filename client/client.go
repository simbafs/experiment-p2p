package client

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "start a client",
	Run: func(cmd *cobra.Command, args []string) {
		// log.Println("client")
		client(cmd, args)
	},
}

var (
	tag  string
	addr string
	port int

	serverIP net.IP
)

func init() {
	Cmd.Flags().StringVarP(&tag, "tag", "t", "", "client tag")
	Cmd.Flags().StringVarP(&addr, "addr", "a", "139.162.86.217", "server addr")
	Cmd.Flags().IntVarP(&port, "port", "p", 9527, "client port")

	serverIP = net.ParseIP(addr)
	if serverIP == nil {
		serverIP = net.IPv4zero
	}
}

func client(cmd *cobra.Command, args []string) {
	srcAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 9901,
	}
	dstAddr := &net.UDPAddr{
		IP:   serverIP,
		Port: port,
	}
	log.Printf("%s:%d <-> %s:%d\n", srcAddr.IP, srcAddr.Port, dstAddr.IP, dstAddr.Port)
	// error
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Fatalf("DialUDP: %s\n", err)
	}
	if _, err = conn.Write([]byte("hello,I'm new peer:" + tag)); err != nil {
		log.Fatalf("conn.Write: %s\n", err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		log.Fatalf("error during read: %s", err)
	}
	log.Printf("n: %d, remoteAddr:%v\n", n, remoteAddr)
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	log.Printf("local:%s server:%s:%d another:%s:%d\n", srcAddr, remoteAddr.IP, remoteAddr.Port, anotherPeer.IP, anotherPeer.Port)
	bidirectionHole(srcAddr, anotherPeer)
}

func parseAddr(addr string) *net.UDPAddr {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatalf("net.SplitHostPort: %s\n", err)
	}
	portInt, err := strconv.Atoi(port)
	return &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: portInt,
	}
}

func bidirectionHole(srcAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		log.Println("send handshake:", err)
	}

	// send
	go func() {
		r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
		for {
			data := fmt.Sprintf("[%s]: %d", tag, r1.Intn(100))
			time.Sleep(10 * time.Second)
			log.Printf("send: %s\n", data)
			if _, err = conn.Write([]byte(data)); err != nil {
				log.Fatalf("send msg fail: %s", err)
			}
		}
	}()

	// receive
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("receive: %s\n", data[:n])
		}
	}
}
