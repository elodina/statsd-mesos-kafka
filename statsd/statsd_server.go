/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package statsd

import (
	"bufio"
	"net"
	"sync"
)

type StatsDServer struct {
	UDPAddr  string
	incoming chan string

	closeChan chan struct{}
	closed    bool
	closeLock sync.Mutex
}

func NewStatsDServer(addr string) *StatsDServer {
	return &StatsDServer{
		UDPAddr:   addr,
		incoming:  make(chan string, 100), //TODO buffer size should be configurable
		closeChan: make(chan struct{}, 1),
	}
}

func (s *StatsDServer) Start() {
	s.startUDPServer()
	s.startProducer()
}

func (s *StatsDServer) Stop() {
	s.closeLock.Lock()
	defer s.closeLock.Unlock()

	if s.closed {
		return
	}

	Logger.Info("Stopping StatsD server")
	s.closeChan <- struct{}{}
	close(s.incoming)
	s.closed = true
}

func (s *StatsDServer) startUDPServer() {
	Logger.Debugf("Starting StatsD server at %s", s.UDPAddr)
	udpAddr, err := net.ResolveUDPAddr("udp", s.UDPAddr)
	if err != nil {
		panic(err)
	}

	connection, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-s.closeChan:
				return
			default:
			}

			s.scan(connection)
		}
	}()
	Logger.Infof("Listening for messages at UDP %s", s.UDPAddr)
}

func (s *StatsDServer) scan(connection net.Conn) {
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		s.incoming <- scanner.Text()
	}
}

func (s *StatsDServer) startProducer() {
	for message := range s.incoming {
		Logger.Infof("Got message: %s", message) //TODO
	}
}
