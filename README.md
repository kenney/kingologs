kingologs
=========

King o' Logs, or KinGoLogs is a demultiplexer for syslog to reliably send logs to 3rd party syslog servers as well as AWS Kinesis.

In its simplest mode of operation both kingologs and the syslog service run on the local box, with syslog connecting to kingologs and kingologs forwarding on to remote destinations.

King o' Logs is currently very experimental and not full featured.  Using a system like fluentd instead is encouraged.

King o' Logs will eventually offer key features for reliable and efficient syslog delivery, such as:
* queuing up remote log delivery to protect against loss of log data during service disruptions
* deduping local records to save storage
* efficiently sending traffic to AWS Kinesis

King o' Logs is not meant to replace syslog servers such as rsyslog or syslog-ng. It complements these solutions but provides for added reliability and capabilities.

## Usage
---
```
Usage: kingologs [args]
 -config configfile     configuration yaml file to process (defaults to /etc/kingologs/config.yml)
 -debug true            force debug mode
```

### Example
Start the daemon:
```
$ go run kingologs.go
```

Send a test message
```
$ echo "Syslog Test" | nc localhost 65140
```

Example output
```
INFO: 2015/11/13 22:45:01 server.go:62: Getting ready to create a new server
INFO: 2015/11/13 22:45:01 kinesis.go:22: Beginning to start kinesis
INFO: 2015/11/13 22:45:01 kinesis.go:70: Checking to see if we have a 'kingologs-test' stream already
INFO: 2015/11/13 22:45:02 kinesis.go:90: Checking to see if kingologs-test stream exists
INFO: 2015/11/13 22:45:02 kinesis.go:44: Done starting server
INFO: 2015/11/13 22:45:02 server.go:76: Beginning to start server
INFO: 2015/11/13 22:45:02 kinesis.go:56: Starting to watch channel
INFO: 2015/11/13 22:45:02 server.go:20: Starting server at: 127.0.0.1:65140
INFO: 2015/11/13 22:45:40 kinesis.go:135: Attempting to put record into stream
INFO: 2015/11/13 22:45:40 kinesis.go:157: Done putting record!
```

# Kinesis integration
kingologs is capable of forwarding on the syslog traffic to the AWS Kinesis service.  

## Configuration

### AWS Configuration
In order to use the Kinesis capabilities the AWS credentials must be present in your ENV and the AWS region must be defined in the config.yml file.

### Config options
The config file has options for:
* service.Name - the name of the daemon
* service.Hostname - the hostname to run as
* Connection.TCP.Enabled - whether to run the TCP syslog server
* Connection.TCP.Host - what host to run as
* Connection.TCP.Port - what port to run on
* Kinesis.StreamName - what stream name to use (will be created if it doesn't exist)
* Kinesis.Region - which AWS region to send the Kinesis data to
* Debug.Verbose - turn on verbose logging

## Testing
If you'd like to test relaying traffic from syslog you can run the syslogtail.sh script in this repo. It tails the syslog file on disk and pumps it to kingologs.

## TODO
Work not yet done

### TODO - Advanced Kinesis integration
By default the system should utilize batch PUT mode to reduce the total number of Kinesis API calls.

Planned options for flushing to Kinesis
* batchPutIdealQuantity:  how many records to queue before performing a batch PUT
* batchPutIdealSize:  ideal packet size
* flushIntervalTime:  time between flushes in the event the ideal
* recordDelimiter: what character(s) by which to separate each syslog record
