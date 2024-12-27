package rcon

import (
	"errors"
	"fmt"
	"log"
)

// Client acts as the entrypoint to send RCON messages. Client hold
// all the configuration necessary to connect to Minecraft server.
type Client struct {
	addr     string
	password string
	channel  *Conn
}

// NewClient creates and returns a configured RCON client.
func NewClient(addr, password string) *Client {
	return &Client{
		addr:     addr,
		password: password,
	}
}

// Send establishes a new authenticated connection to the Minecraft
// server and transmits requested command. Not all commands generate
// a response from the server. Any response from the server is returned
// to requester. If connection failure occurs an error is returned.
//
// This function is concurrency-safe since each command sent to the
// Minecraft server creates a new connection. Upon completion of the
// request the established connection is closed.
func (c *Client) Send(command string) (string, error) {
	err := errors.New("")
	c.channel, err = Dial(c.addr, c.password)
	if err != nil {
		return "", fmt.Errorf("failed to establish connection: %w", err)
	}
	defer func() {
		if err = c.channel.Close(); err != nil {
			log.Fatalf("close connection fail:%s", err)
		}
	}()

	return c.channel.SendCommand(command)
}

func (c *Client) AutoReconnect() {
	if !c.channel.isClosed {
		return
	}
	c.channel.start()
}
