package broker

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/network"
	"fmt"
	"net"
)

type Processor struct {
	Conn net.Conn
	Tr   *network.Transport
}

func (p *Processor) brokerProcessMes() (err error) {
	p.Tr = &network.Transport{
		Conn: p.Conn,
	}

	for {
		mes, err := p.Tr.ReceiveMessage(p.Tr.Conn)
		if err != nil {
			if err.Error() == "EOF" {
				log.LogMessage("INFO", "Client closed the connection")
				p.removeConsumer(p.Conn)
				return nil
			}
			p.removeConsumer(p.Conn)
			return err
		}

		log.LogMessage("INFO", "Broker received a message: "+string(mes.Payload))

		// Broadcast the message to all consumers
		consumersMutex.Lock()
		var activeConsumers []net.Conn
		for _, consumer := range consumers {
			tr := &network.Transport{
				Conn: consumer,
			}
			err := tr.SendMessage(consumer, mes)
			if err != nil {
				// fmt.Printf("Failed to send message to consumer: %v\n", err)
				log.LogMessage("ERROR", fmt.Sprintf("Failed to send message to consumer: %v", err))
				consumer.Close() // Close the connection if sending fails
			} else {
				activeConsumers = append(activeConsumers, consumer)
			}
		}
		consumers = activeConsumers
		consumersMutex.Unlock()
	}
}

func (p *Processor) removeConsumer(conn net.Conn) {
	consumersMutex.Lock()
	defer consumersMutex.Unlock()

	for i, consumer := range consumers {
		if consumer == conn {
			consumers = append(consumers[:i], consumers[i+1:]...)
			break
		}
	}
}

package main

import (
	"COMP47250-Team-Software-Project/internal/network"
	"fmt"
	"net"
)

type Processor struct {
	Conn net.Conn
}

func (p *Processor) brokerProcessMes() (err error) {
	tr := &network.Transport{
		Conn: p.Conn,
	}

	mes, err := tr.ReceiveMessage(tr.Conn)
	mes, err := p.Tr.ReceiveMessage(p.Tr.Conn)
	if err != nil {
		fmt.Println("Broker can not receive message successfully!!")
		return
	}

	fmt.Println("Broker received a message: ", mes.Payload)
	// fmt.Println("The message detail:")
	// fmt.Println("ID: ", mes.ID)
	// fmt.Println("Timestamp: ", mes.Timestamp)
	// fmt.Println("Type: ", mes.Type)
}
