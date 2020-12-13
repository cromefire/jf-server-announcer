package main

import (
	"encoding/json"
	"github.com/jessevdk/go-flags"
	"log"
	"net"
	"os"
)

var opts struct {
	Id      string `short:"i" long:"id" description:"The ID of the Jellyfin Server to advertise" required:"true"`
	Address string `short:"a" long:"address" description:"The address at which the clients should reach the server" required:"true"`
	Name    string `short:"n" long:"name" description:"The name of the server the clients should show" required:"true"`
}

type AnnounceMsg struct {
	Id              *string `json:"Id"`
	Name            *string `json:"Name"`
	Address         *string `json:"Address"`
	EndpointAddress *string `json:"EndpointAddress"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	// listen to incoming udp packets
	s, err := net.ResolveUDPAddr("udp", ":7359")
	pc, err := net.ListenUDP("udp", s)
	if err != nil {
		log.Fatal(err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer pc.Close()

	msg := AnnounceMsg{Id: &opts.Id, Name: &opts.Name, Address: &opts.Address, EndpointAddress: nil}
	send, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Server running...")
	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		req := string(buf[0:n])
		if req != "who is JellyfinServer?" {
			continue
		}

		go func() {
			_, err = pc.WriteTo(send, addr)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}
