// mystack-controller api
// https://github.com/topfreegames/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"encoding/json"
	"io"
	"net"
	"sync"

	"github.com/Sirupsen/logrus"
)

type Proxy struct {
	from     string
	to       string
	done     chan struct{}
	messages map[string]interface{}
	log      *logrus.Entry
}

func NewProxy(from, to string, messages map[string]interface{}) *Proxy {
	return &Proxy{
		from: from,
		to:   to,
		done: make(chan struct{}),
		log: logrus.WithFields(logrus.Fields{
			"from": from,
		}),
		messages: messages,
	}
}

func (p *Proxy) Start() error {
	listener, err := net.Listen("tcp", p.from)
	if err != nil {
		return err
	}
	p.run(listener)
	return nil
}

func (p *Proxy) Stop() {
	p.log.Infoln("Stopping proxy")
	if p.done == nil {
		return
	}
	close(p.done)
	p.done = nil
}

func (p *Proxy) run(listener net.Listener) {
	for {
		select {
		case <-p.done:
			return
		default:
			connection, err := listener.Accept()
			if err == nil {
				go p.handle(connection)
			} else {
				p.log.WithField("err", err).Errorln("Error accepting conn")
			}
		}
	}
}

func (p *Proxy) handle(connection net.Conn) {
	p.log.Debugln("Handling", connection)
	defer p.log.Debugln("Done handling", connection)
	defer connection.Close()
	remote, err := net.Dial("tcp", p.to)
	if err != nil {
		p.log.WithField("err", err).Errorln("Error dialing remote host")
		return
	}
	defer remote.Close()

	p.handshake(remote)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go p.copy(remote, connection, wg)
	go p.copy(connection, remote, wg)

	wg.Wait()
}

func (p *Proxy) handshake(remote net.Conn) {
	bts, err := json.Marshal(p.messages)
	if err != nil {
		p.log.Fatalf("invalid handshake message: %s\n", err.Error())
	}

	bts = append(bts, '\n')

	_, err = remote.Write(bts)
	if err != nil {
		p.log.Fatalf("handshake error: %s\n", err.Error())
	}

	authRes := make([]byte, 1024)
	_, err = remote.Read(authRes)
	if err != nil {
		p.log.Fatalf("read ack error: %s\n", err.Error())
	}
}

func (p *Proxy) copy(from, to net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-p.done:
		return
	default:
		if _, err := io.Copy(to, from); err != nil {
			p.log.WithField("err", err).Errorln("Error from copy")
			p.Stop()
			return
		}
	}
}
