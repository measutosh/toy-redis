package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number")
		return
	}

  port := ":" + arguments[1]
  server, err := net.Listen("tcp4", port)
  if err != nil {
    fmt.Println(err)
  }
  defer server.Close()

  // main goroutine runs infinitely
  // when it gets a request from a client, 
  // it runs that request in a separate goroutine

  for{
    client, err := server.Accept()
    if err != nil {
      fmt.Println(err)
      return
    }

    go handleConnection(client)
  }
}










