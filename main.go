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
		timerChan: make(chan time.Time),
	}

	go tp.cli.run()
	tp.cli.event <- Start
}

type client struct {
	state     chan txState
	event     chan txEvent
	currState txState
	timerChan <-chan time.Time
	timer     *time.Timer
}

func (c *client) run() {
	for {
		select {
		case state := <-c.state:
			c.currState = state
		case <-c.timerChan:
			go func() {
				c.event <- Timeout
			}()
		case event := <-c.event:
			switch c.currState {
			case Idle:
				go c.Idle(event)
			case SendData:
				go c.SendData(event)
			}
		}
	}
}

func (c *client) startTimer() {
	c.timer = time.NewTimer(2 * time.Second)
	c.timerChan = c.timer.C
}

func (c *client) Idle(ev txEvent) {
	fmt.Printf("Idle state: %d\n", ev)
	c.startTimer()
	c.state <- SendData
	fmt.Println("state changed to sending data")
}

func (c *client) SendData(ev txEvent) {
	fmt.Printf("Sending data: %d\n", ev)
}

func main() {
	tp := transport{}
	go tp.sendMSG()

	var input string
	fmt.Scanln(&input)
}
