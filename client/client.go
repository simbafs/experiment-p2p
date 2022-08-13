package client

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "start a client",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("client")
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
	Cmd.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1", "server addr")
	Cmd.Flags().IntVarP(&port, "port", "p", 9901, "client port")

	serverIP = net.ParseIP(addr)
	if serverIP == nil {
		serverIP = net.IPv4zero
	}
}

const HAND_SHAKE_MSG = "this is a tunnel msg"

func client(cmd *cobra.Command, args []string) {
	srcAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 9527,
	}
	dstAddr := &net.UDPAddr{
		IP:   serverIP,
		Port: port,
	}
	// error
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Println(err)
	}
	if _, err = conn.Write([]byte("hello,I'm new peer:" + tag)); err != nil {
		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		log.Printf("error during read: %s", err)
	}
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	log.Printf("local:%s server:%s another:%s\n", srcAddr, remoteAddr, anotherPeer)
	bidirectionHole(srcAddr, &anotherPeer)
}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {

		fmt.Println("send handshake:", err)
	}
	go func() {

		for {

			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + tag + "]")); err != nil {

				log.Println("send msg fail", err)
			}
		}
	}()

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
