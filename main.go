package main

import (
	"fmt"
	"time"
)

type txState int

const (
	Idle txState = iota
	SendData
)

type txEvent int

const (
	Start txEvent = iota
	Timeout
)

type transport struct {
	cli client
}

func (tp *transport) sendMSG() {
	tp.cli = client{
		state:     make(chan txState),
		event:     make(chan txEvent),
		currState: Idle,
	}

	go tp.cli.run()
	tp.cli.event <- Start
}

type client struct {
	state     chan txState
	event     chan txEvent
	currState txState
}

func (c *client) run() {
	for {
		select {
		case state := <-c.state:
			c.currState = state
		case event := <-c.event:
			switch c.currState {
			case Idle:
				c.Idle(event)
			case SendData:
				c.SendData(event)
			}
		}
	}
}

func (c *client) Idle(ev txEvent) {
	fmt.Printf("Idle state: %d\n", ev)
	c.state <- SendData
	time.Sleep(2 * time.Second)
	c.event <- Timeout
}

func (c *client) SendData(ev txEvent) {
	fmt.Printf("Sending data: %d\n", ev)
}

func main() {
	tp := transport{}
	go tp.sendMSG()

	var input string
	fmt.Scan(&input)
}
