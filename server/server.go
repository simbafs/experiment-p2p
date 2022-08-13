package server

import (
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "start a server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("server")
		server(cmd, args)
	},
}

var (
	port int
)

func init() {
	Cmd.Flags().IntVarP(&port, "port", "p", 9527, "port to listen")
}

func server(cmd *cobra.Command, args []string) {
	log.Println("port", port)

	// create a udp listener
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	})
	if err != nil {

		log.Println(err)
	}
	log.Printf("local address: <%s> \n", listener.LocalAddr().String())
	peers := make([]net.UDPAddr, 0, 2)
	data := make([]byte, 1024)
	for {

		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {

			log.Printf("err during read: %s\n", err)
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
		peers = append(peers, *remoteAddr)
		if len(peers) == 2 {

			log.Printf("Creating UDP tunnel...... pear: %s <--> %s \n", peers[0].String(), peers[1].String())
			// writeToUDP send data through channel c.
			// if timeout, it return an error, which is very rare
			listener.WriteToUDP([]byte(peers[1].String()), &peers[0])
			listener.WriteToUDP([]byte(peers[0].String()), &peers[1])
			time.Sleep(time.Second * 8)
			log.Println("server exit")
			return
		}
	}
}
