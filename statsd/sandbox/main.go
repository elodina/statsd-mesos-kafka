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

package main

import (
	"github.com/stealthly/siesta"
	"github.com/stealthly/statsd-mesos-kafka/statsd"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go http.ListenAndServe(":31099", nil)

	producerConfig := siesta.NewProducerConfig()
	producerConfig.BrokerList = []string{"localhost:9092"}

	connectorConfig := siesta.NewConnectorConfig()
	connectorConfig.BrokerList = []string{"localhost:9092"}

	connector, err := siesta.NewDefaultConnector(connectorConfig)
	if err != nil {
		panic(err)
	}

	statsd.InitLogging("debug")
	statsd.Config.Topic = "metrics"

	producer := siesta.NewKafkaProducer(producerConfig, siesta.ByteSerializer, siesta.StringSerializer, connector)

	server := statsd.NewStatsDServer("0.0.0.0:8125", producer, func(message string, host string) interface{} {
		return message
	}, "127.0.0.1")

	go server.Start()

	select {}
}
