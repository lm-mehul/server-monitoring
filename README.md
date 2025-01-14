# server-monitoring
server-monitoring


# Setup 

go get github.com/shirou/gopsutil

go run main.go


The script above provides the functionality :

1. /start-monitoring Endpoint (POST):

Starts monitoring system metrics every 0.5 seconds for 30 seconds.
Saves the collected data in a JSON file, overwriting the previous results.

2. /cpu-metrics Endpoint (GET):

Returns the content of the JSON file containing the collected metrics.
