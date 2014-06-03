# statsd traffic loader

## Overview
CLI script to generate traffic for a statsd instance with configurable rates

## Usage
```
% go get
% go run main.go -s 20000
usage: statsd-traffic-loader [-sp] hostname
  -p=8125: port to send to
  -s=30000: packets per second to send
exit status 1
```
