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
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
	"os"
	"os/signal"
)

type Scheduler struct {
}

func (s *Scheduler) Start() error {
	Logger.Infof("Starting scheduler with configuration: \n%s", Config)

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	frameworkInfo := &mesosproto.FrameworkInfo{
		User:       proto.String(Config.User),
		Name:       proto.String(Config.FrameworkName),
		Role:       proto.String(Config.FrameworkRole),
		Checkpoint: proto.Bool(true),
	}

	driverConfig := scheduler.DriverConfig{
		Scheduler: s,
		Framework: frameworkInfo,
		Master:    Config.Master,
	}

	driver, err := scheduler.NewMesosSchedulerDriver(driverConfig)
	go func() {
		<-ctrlc
		s.Shutdown(driver)
	}()

	if err != nil {
		return fmt.Errorf("Unable to create SchedulerDriver: %s", err)
	}

	if stat, err := driver.Run(); err != nil {
		Logger.Infof("Framework stopped with status %s and error: %s\n", stat.String(), err)
		return err
	}

	return nil
}

func (s *Scheduler) Registered(scheduler.SchedulerDriver, *mesos.FrameworkID, *mesos.MasterInfo) {

}

func (s *Scheduler) Reregistered(scheduler.SchedulerDriver, *mesos.MasterInfo) {

}

func (s *Scheduler) Disconnected(scheduler.SchedulerDriver) {

}

func (s *Scheduler) ResourceOffers(scheduler.SchedulerDriver, []*mesos.Offer) {

}

func (s *Scheduler) OfferRescinded(scheduler.SchedulerDriver, *mesos.OfferID) {

}

func (s *Scheduler) StatusUpdate(scheduler.SchedulerDriver, *mesos.TaskStatus) {

}

func (s *Scheduler) FrameworkMessage(scheduler.SchedulerDriver, *mesos.ExecutorID, *mesos.SlaveID, string) {

}

func (s *Scheduler) SlaveLost(scheduler.SchedulerDriver, *mesos.SlaveID) {

}

func (s *Scheduler) ExecutorLost(scheduler.SchedulerDriver, *mesos.ExecutorID, *mesos.SlaveID, int) {

}

func (s *Scheduler) Error(scheduler.SchedulerDriver, string) {

}

func (s *Scheduler) Shutdown(driver *scheduler.MesosSchedulerDriver) {
	Logger.Info("Shutdown triggered, stopping driver")
	driver.Stop(false)
}
