Statsd Mesos Kafka Framework
============================

Installation
------------

Install go 1.4 (or higher) http://golang.org/doc/install

Install godep https://github.com/tools/godep

Clone and build the project

    # git clone git@github.com:stealthly/statsd-mesos-kafka.git
    # cd statsd-mesos-kafka
    # godep restore
    # go build .
    # go build cli.go

Usage
-----

Statsd framework ships with command-line utility to manage schedulers and executors:

    # ./cli -help
    Usage:
        help: show this message
        scheduler: configure and start scheduler
        start: start statsd server
        stop: stop statsd server
        update: update configuration
        status: get current status of cluster
    More help you can get from ./cli <command> -h


Scheduler Configuration
-----------------------

The scheduler is configured through the command line.

    # ./cli scheduler <options>

Following options are available:

-master="": Mesos Master addresses.
-api="": Binding host:port for http/artifact server. Optional if SM_API env is set.
-user="": Mesos user. Defaults to current system user.
-log.level="info": Log level. trace|debug|info|warn|error|critical. Defaults to info.
-framework.name="statsd-kafka": Framework name.
-framework.role="*": Framework role.

Starting and Stopping a Server
------------------------------

    # ./cli start|stop <options>

Options available:
    -api="": Binding host:port for http/artifact server. Optional if SM_API env is set.

Updating Server Preferences
---------------------------

    # ./cli update <options>

Following options are available:

    -api="": Binding host:port for http/artifact server. Optional if SM_API env is set.
    -producer.properties="": Producer.properties file name.
    -topic="": Topic to produce data to.
    -transform="": Transofmation to apply to each metric. none|avro|proto
    -schema.registry.url="": Avro Schema Registry url for transform=avro

