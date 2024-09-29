package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
)

const (
	ClrfDelimeter = "\r\n"
)

func handleClient(conn net.Conn) {
	// ensure we close the connection after we're done
	defer conn.Close()

	// Read data
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading data from  the client", err.Error())
			return
		}

		input_buf := string(buf)
		output, err := commands.Execute(input_buf)
		fmt.Println("Output string is", output)

		if err != nil {
			fmt.Println("Error executing command: ", err.Error())
		}

		// message := []byte("+PONG\r\n")
		_, err = conn.Write([]byte(output))

		if err != nil {
			fmt.Println("Error sending message to the  client", err.Error())
			return
		}
		handle_psync(input_buf, conn)
		// fmt.Printf("send %d bytes", n)

	}

}

func handle_psync(input_buf string, conn net.Conn) {
	command_slice := commands.Decode(input_buf)
	fmt.Println("Command slice is  ", command_slice)

	name := strings.ToLower(command_slice[0])
	if name == "psync" {
		sync_data := "$87\r\nREDIS0011\xfa\tredis-ver\x057.2.0\xfa\nredis-bits\xc0@\xfa\x05ctime\xc2m\x08\xbce\xfa\x08used-mem\xc2\xb0\xc4\x10\x00\xfa\x08aof-base\xc0\x00\xff\xf0n;\xfe\xc0\xffZ"
		_, err := conn.Write([]byte(sync_data))
		fmt.Println("Sent empty RDB file")
		if err != nil {
			fmt.Println("Error sending empty RDB to the client  ", err.Error())
			return
		}

	}
}

func HandShake(master_host string, master_port string, listening_port string) {

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", master_host, master_port))
	if err != nil {
		fmt.Println("Error in handshaking:", err)
		os.Exit(1)
	}
	defer conn.Close()

	message_slice := [][]string{{"PING"},
		{"REPLCONF", "listening-port", listening_port},
		{"REPLCONF", "capa", "psync2"},
		{"PSYNC", "?", "-1"},
	}

	for _, message := range message_slice {

		bulk_message, _ := commands.CreateMessage(message)
		// fmt.Println("Bulk message is :", bulk_message)

		_, err = conn.Write([]byte(bulk_message))
		if err != nil {
			fmt.Println("Error sending message over Handshake:", err)
			os.Exit(1)

		}
		fmt.Println("Sent message:", message)

		// Wait for the server's response
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving response:", err)
			return
		}

		// Print the server's response
		fmt.Println("Server response:", response)

	}

	fmt.Println("Successfully HandShake done and Message sent")

}

func main() {
	port := flag.String("port", "6379", "server port")
	replicaof := flag.String("replicaof", "master", "decalre server as slave or master")

	// Parse the flags
	flag.Parse()

	server_port := fmt.Sprintf("0.0.0.0:%s", *port)

	if *replicaof != "master" {
		commands.DEFAULTROLE = "slave"
		master_info := strings.Split(*replicaof, " ")
		master_host := master_info[0]
		master_port := master_info[1]
		fmt.Println("Master info", master_host, master_port)
		go HandShake(master_host, master_port, *port)

	}

	// listener, err := net.Listen("tcp", "0.0.0.0:6379")
	listener, err := net.Listen("tcp", server_port)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	// ensure to stop the tcp server when the program exits
	defer listener.Close()

	fmt.Println("Server is listening on port:", *port)
	// Block until we receive an incoming connection

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)

		}
		// Hnadle client connection
		go handleClient(conn)

	}

}
