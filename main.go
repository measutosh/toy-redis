package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Command int

const (
	Get Command = iota + 1
	Set
	Incr
	Del
)

type commandMessage struct {
	commandName     Command
	key             string
	value           string
	responseChannel chan string
}

func handleDB(commandChannel chan CommandMessage) {
    db := make(map[string]string)

    for {
        select {
            case command := <- commandChannel:
                switch command.commandName {
                    case Get:
                        command.responseChannel <- db[command.key]
                    case Set:
                        db[command.key] = command.value
                        command.responseChannel <- "OK"
                    case Incr:
                        value, ok := db[command, key]
                        var response string

                        if ok {
                            intValue, err := strconv.Atoi(value)
                            if err != nil {
                                response = "ERR value is noy an integer or out of range"
                            } else {
                                response = strconv.Itoa(intValue + 1)
                                db[command.key] = response
                            }
                        } else {
                            response = "1"
                            db[command.key] = response
                        }
    
                        command.responseChannel <- response
                    case Del:
                        _, ok := db[command.key]
                    var response string

                    if ok {
                        delete(db, command.key)
                        response = "1"
                    } else {
                        response = "0"
                    }
                    command.responseChannel <- response
                }
        }
    }
}

func handleConnection(commandChannel chan commandMessage, client net.Conn) {

	defer client.Close()

	fmt.Printf("Serving %s\n", client.RemoteAddr().String())

	for {
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
				commandMessage := commandMessage{
					commandName:     Get,
					key:             key,
					responseChannel: make(chan string),
				}
				commandChannel <- commandMessage
				response = <-commandMessage.responseChannel
			} else {
				response = "ERR wrong number of arguments for 'get' command"
			}
		case "SET":
			if len(parts) > 2 {
				key := parts[1]
				value := parts[2]
				commandMessage := commandMessage{
					commandName:     Set,
					key:             key,
					value:           value,
					responseChannel: make(chan string),
				}
				commandChannel <- commandMessage
				response = <-commandMessage.responseChannel

			} else {
				response = "ERR wrong number of arguments for 'set' command"
			}
		case "DEL":
			key := parts[1]
			commandMessage := commandMessage{
				commandName:     Del,
				key:             key,
				responseChannel: make(chan string),
			}

			commandChannel <- commandMessage
			response = <-commandMessage.responseChannel
		case "INCR":
			// if values exists then convert the text into int
			// increase the string
			// convert it back to string and store in the map
			if len(parts) > 1 {
				key := parts[1]
				commandMessage := commandMessage{
					commandName:     Incr,
					key:             key,
					responseChannel: make(chan string),
				}
				commandChannel <- commandMessage
				response = <-commandMessage.responseChannel
			} else {
				response = "ERR wrong number of arguments for 'incr' command"
			}
		default:
			response = "ERR unknown command"
		}
		client.Write([]byte(response + "\n"))
	}

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


	commandChannel := make(chan commandMessage)
	go handleDB(commandChannel)

	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(commandChannel, client)
	}
}
