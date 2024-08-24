package main

import (
	"fmt"
	"net"
	"os"
)

func  handleClient(conn net.Conn){
	// ensure we close the connection after we're done
	defer conn.Close()

	// Read data
	buf := make([]byte,1024)
	n,err := conn.Read(buf)

	if err != nil{
		fmt.Println("Error reading data from  the client",err.Error())
		return
	}
	fmt.Println("Received data",string(buf[:n]))


	message := []byte("Hello , server!")
	n,err := conn.Write(message)

	if err != nil{
		fmt.Println("Error sneding message to the  client",err.Error())
		return
	}
	fmt.Println("send %d bytes",n)



}

func main() {
	
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	// ensure to stop the tcp server when the program exits
	defer listener.close()

	fmt.Println("Server is listening on port 6379")
	for {
		// Block until we receive an incoming connection
		conn, err = listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)

		}
		// Hnadle client connection
		handleClient(conn)

		
	}
	

}
