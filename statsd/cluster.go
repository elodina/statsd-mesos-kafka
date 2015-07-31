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
	mesos "github.com/mesos/mesos-go/mesosproto"
	"sync"
)

type Cluster struct {
	tasks    map[string]*mesos.TaskInfo
	taskLock sync.Mutex
}

func NewCluster() *Cluster {
	return &Cluster{
		tasks: make(map[string]*mesos.TaskInfo),
	}
}

func (c *Cluster) Exists(hostname string) bool {
	c.taskLock.Lock()
	defer c.taskLock.Unlock()

	_, exists := c.tasks[hostname]
	return exists
}

func (c *Cluster) Add(hostname string, task *mesos.TaskInfo) {
	c.taskLock.Lock()
	defer c.taskLock.Unlock()

	if _, exists := c.tasks[hostname]; exists {
		// this should never happen. would mean a bug if so
		panic(fmt.Sprintf("statsd-kafka at %s already exists", hostname))
	}

	c.tasks[hostname] = task
}

func (c *Cluster) Remove(hostname string) {
	c.taskLock.Lock()
	defer c.taskLock.Unlock()

	delete(c.tasks, hostname)
}

func (c *Cluster) GetAllTasks() []*mesos.TaskInfo {
	c.taskLock.Lock()
	defer c.taskLock.Unlock()

	tasks := make([]*mesos.TaskInfo, 0)
	for _, task := range c.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}
