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
	"errors"
	"flag"
	"fmt"
	"github.com/stealthly/statsd-mesos-kafka/statsd"
	"os"
)

func main() {
	if err := exec(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func exec() error {
	args := os.Args
	if len(args) == 1 {
		handleHelp(nil)
		return errors.New("No command supplied")
	}

	command := args[1]
	commandArgs := args[1:]
	os.Args = commandArgs

	if command == "help" {
		handleHelp(commandArgs)
		return nil
	} else if command == "scheduler" {
		return handleScheduler(commandArgs)
	} else if command == "start" || command == "stop" {
		return handleStartStop(commandArgs, command == "start")
	} else if command == "update" {
		return handleUpdate(commandArgs)
	} else {
		return fmt.Errorf("Unknown command: %s\n", command)
	}
}

func handleHelp(commandArgs []string) {
	fmt.Println("help message") //TODO
}

func handleScheduler(commandArgs []string) error {
	var api string
	var user string
	var logLevel string

	flag.StringVar(&statsd.Config.Master, "master", "", "Mesos Master addresses.")
	flag.StringVar(&api, "api", "", "Binding host:port for http/artifact server. Optional if SM_API env is set.")
	flag.StringVar(&user, "user", "", "Mesos user. Defaults to current system user")
	flag.StringVar(&logLevel, "log.level", statsd.Config.LogLevel, "Log level. trace|debug|info|warn|error|critical. Defaults to info.")
	flag.StringVar(&statsd.Config.FrameworkName, "framework.name", statsd.Config.FrameworkName, "Framework name.")
	flag.StringVar(&statsd.Config.FrameworkRole, "framework.role", statsd.Config.FrameworkRole, "Framework name.")

	flag.Parse()

	if err := resolveApi(api); err != nil {
		return err
	}

	if err := statsd.InitLogging(logLevel); err != nil {
		return err
	}

	if statsd.Config.Master == "" {
		return errors.New("--master flag is required.")
	}

	return new(statsd.Scheduler).Start()
}

func handleStartStop(commandArgs []string, start bool) error {
	var api string
	flag.StringVar(&api, "api", "", "Binding host:port for http/artifact server. Optional if SM_API env is set.")

	flag.Parse()

	if err := resolveApi(api); err != nil {
		return err
	}

	apiMethod := "start"
	if !start {
		apiMethod = "stop"
	}

	request := statsd.NewApiRequest(statsd.Config.Api + "/api/" + apiMethod)
	response := request.Get()

	fmt.Println(response.Message)

	return nil
}

func handleUpdate(commandArgs []string) error {
	var api string
	flag.StringVar(&api, "api", "", "Binding host:port for http/artifact server. Optional if SM_API env is set.")
	flag.StringVar(&statsd.Config.ProducerProperties, "producer.properties", "", "Producer.properties file name.")
	flag.StringVar(&statsd.Config.Topic, "topic", "", "Topic to produce data to.")
	flag.StringVar(&statsd.Config.Transform, "transform", "", "Transofmation to apply to each metric. none|avro|proto")
	flag.StringVar(&statsd.Config.SchemaRegistryUrl, "schema.registry.url", "", "Avro Schema Registry url for transform=avro")

	flag.Parse()

	if err := resolveApi(api); err != nil {
		return err
	}

	request := statsd.NewApiRequest(statsd.Config.Api + "/api/update")
	request.AddParam("producer.properties", statsd.Config.ProducerProperties)
	request.AddParam("topic", statsd.Config.Topic)
	request.AddParam("transform", statsd.Config.Transform)
	request.AddParam("schema.registry.url", statsd.Config.SchemaRegistryUrl)
	response := request.Get()

	fmt.Println(response.Message)

	return nil
}

func resolveApi(api string) error {
	if api != "" {
		statsd.Config.Api = api
		return nil
	}

	if os.Getenv("SM_API") != "" {
		statsd.Config.Api = os.Getenv("SM_API")
		return nil
	}

	return errors.New("Undefined API url. Please provide either a CLI --api option or SM_API env.")
}
