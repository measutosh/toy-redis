package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleConnection(db map[string]string, client net.Conn) {
  // recieves the message sent by the client using bufio
  // runs the goroutine infinitely

  defer client.Close()

  fmt.Printf("Serving %s\n", client.RemoteAddr().String())

  for{
    netData, err := bufio.NewReader(client).ReadString('\n')
    if err != nil {
      fmt.Println("error reading: ", err)
      return
    }

    var response string
    commandString := strings.TrimSpace(netData)
    parts := strings.Split(commandString, " ")
    command := parts[0]

    switch command {
      case "STOP", "QUIT":
        return
      case "GET":
        if len(parts) > 1 {
        	key := parts[1]
        
        	response = db[key]
        } else {
        	response = "ERR wrong number of arguments for 'get'"
        }
      case "SET":
        if len(parts) > 2 {
        	key := parts[1]
        	value := parts[2]
        
        	db[key] = value
        	response = "OK"
        } else {
        	response = "ERR wrong number of arguments for 'set' command"
        }
      case "DEL":
        if len(parts) > 1 {
          key := parts[1]
          _, ok := db[key]

          if ok {
            delete(db, key)
            response = "1"
          } else {
            response = "0"
          }
        } else {
          response = "ERR wrong number of arguments for 'del'"
        }
      case "INCR":
        // if values exists then convert the text into int
        // increase the string
        // convert it back to string and store in the map
        if len(parts) > 1 {
        	key := parts[1]
        	value, ok := db[key]
        
        	if ok {
        		intValue, err := strconv.Atoi(value)
        		if err != nil {
        			response = "ERR value is not an integer or out of range"
        		} else {
        			response = strconv.Itoa(intValue + 1)
        			db[key] = response
        		}
        	} else {
        		response = "1"
        		db[key] = response
        	}
        } else {
        	response = "ERR wrong number of arguments for 'incr' command"
        }
      default:
        response = "ERR unknown command"
    }
    client.Write([]byte(response + "\n"))
  }

  fmt.Println("Closing the client")
  client.Close()
}

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
    return 
  }
  defer server.Close()

  // main goroutine runs infinitely
  // when it gets a request from a client, 
  // it runs that request in a separate goroutine
  // map acts as database and is passed to all goroutines
  db := make(map[string]string)

  for{
    client, err := server.Accept()
    if err != nil {
      fmt.Println(err)
      return
    }

    go handleConnection(db, client)
  }
}










