
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        fmt.Println("Can not get IP")
				os.Exit(100)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Press J to join a game or press H to host a game")
		var answer string
		fmt.Scanln(&answer)
		if ((strings.ToLower(strings.TrimSpace(answer))) == "h"){
			fmt.Println("You chose to host a game your IP is " + GetOutboundIP().String())
			fmt.Println("Choose a port for your friend to connect to.")
			var PORT string
			fmt.Scanln(&PORT)
			fmt.Println("Your IP and port are " + GetOutboundIP().String() + ":" + PORT + " Please give this to your friend so they can connect.")
			startServer(PORT)
			os.Exit(100)
		}

		if ((strings.ToLower(strings.TrimSpace(answer))) == "j"){
			fmt.Println("Please make a username")
			var username string
			fmt.Scanln(&username)
			fmt.Println("You have chosen to join a game please provide IP and port to conect to a server.")
			var IP string
			fmt.Scanln(&IP)
			fmt.Println("Connecting")
			startClient(IP, username)
		}
	}
}

func startClient(IPPORT string, username string){

	c, err := net.Dial("tcp", IPPORT)
	if err != nil {
		fmt.Println(err)
		os.Exit(100)
	}

	for{
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(username + ">> ")
		text, _ := reader.ReadString('\n')
		text = username + ":" + IPPORT + " ->" + text
		fmt.Fprintf(c, text + "\n")
		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("-> " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("Quiting TCP client")
			return
		}
	}


}


func startServer(PORT string) {
		fmt.Println("Init Server")
		l, err := net.Listen("tcp", ":" + PORT)
		if err != nil {
			fmt.Println(err)
			os.Exit(100)
		}
		defer l.Close()

		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(100)
		}

		for {
			netData, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				os.Exit(100)
			}

			fmt.Print("-> ", string(netData))
			c.Write([]byte(netData))
			if strings.TrimSpace(string(netData)) == "STOP" {
				fmt.Println("Exiting TCP server!")
				return
			}
		}
}