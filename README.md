# statsd traffic loader

## Overview
CLI script to generate traffic for a statsd instance with configurable rates.
It also has the ability to load metrics key names from a file.

## Usage
```
% go get
% go run main.go -s 20000
usage: statsd-traffic-loader [-cpst] hostname
  -c="stats_counter_keys.txt": file with example counter keys
  -p=8125: port to send to
  -s=30000: packets per second to send
  -t="stats_timer_keys.txt": file with example timer keys
```
