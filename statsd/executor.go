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
	"github.com/mesos/mesos-go/executor"
	"github.com/mesos/mesos-go/mesosproto"
	"time"
)

type Executor struct {
}

func (e *Executor) Registered(driver executor.ExecutorDriver, executor *mesosproto.ExecutorInfo, framework *mesosproto.FrameworkInfo, slave *mesosproto.SlaveInfo) {
	Logger.Infof("[Registered] framework: %s slave: %s", framework.GetId().GetValue(), slave.GetId().GetValue())
}

func (e *Executor) Reregistered(driver executor.ExecutorDriver, slave *mesosproto.SlaveInfo) {
	Logger.Infof("[Reregistered] slave: %s", slave.GetId().GetValue())
}

func (e *Executor) Disconnected(executor.ExecutorDriver) {
	Logger.Info("[Disconnected]")
}

func (e *Executor) LaunchTask(driver executor.ExecutorDriver, task *mesosproto.TaskInfo) {
	Logger.Infof("[LaunchTask] %s", task)

	time.Sleep(1 * time.Minute) //stub
}

func (e *Executor) KillTask(driver executor.ExecutorDriver, id *mesosproto.TaskID) {
	Logger.Infof("[KillTask] %s", id.GetValue())
}

func (e *Executor) FrameworkMessage(driver executor.ExecutorDriver, message string) {
	Logger.Infof("[FrameworkMessage] %s", message)
}

func (e *Executor) Shutdown(driver executor.ExecutorDriver) {
	Logger.Infof("[Shutdown]")
}

func (e *Executor) Error(driver executor.ExecutorDriver, message string) {
	Logger.Errorf("[Error] %s", message)
}
