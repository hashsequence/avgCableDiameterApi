# AvgCableDiameterApi

## Challenge

At Oden we deal with a lot of time-series data. We process millions of IoT metrics every day and have to make sure our end users can visualize the data in a timely manner. One common use case for our users is to monitor metrics that are indicators of product quality in real-time. For our cable extrusion customers, this happens to be the cable diameter.

http://takehome-backend.oden.network/?metric=cable-diameter

The API returns the current value for the metric. Your task is to assemble an application that polls the Oden API, calculates a one minute moving average of the Cable Diameter, and* exposes that moving average via an HTTP service when issued a GET request to localhost:8080/cable-diameter.

Example:

```
curl localhost:8080/cable-diameter
11.24
```

Your moving average should be updated, if possible, once per second and, after each update, the new moving average is logged to STDOUT.

## Design

* From my interpretation of the prompt, it seems that the Api must calculate the one minute moving average of the Cable Diameter, meaning we must calculate 
the average of all the cableApi values retrieve from the poll api. By utilizing goroutines, I can have a process polling the poll api every second, and after 60 seconds has passed since polling, we can start popping off the oldest values every second, thus in our dataStore we will always maintain
the subset of the polled values in our one minute interval prior to the current time in real-time.

breakdown of responsibilities: 

* poller
    * Responsible for polling the oden's Api

* dataStore 
    * data structure to store moving average value and related data

* Routes
    * defines Handlers to handle endpoints

* logger
    * maintains the logging of web service

### Scope 

* The specification did not specify authentication nor encryption requirements and is not part of the core challenge, though if implemented the certificates can be create and signed using openssl, I have taken the liberty to create the certificates for server side authentication and client side authentication,
and a bash script is available in the ssl folder (genCerts.sh) to create said certificates, and we can use the following to enable tls if needed: 


    https://golang.org/pkg/net/http/#Server.ListenAndServeTLS 
    
    "crypto/tls"

    "crypto/x509"


Though in my implementation the Api is an insecure public API using http

* The DataStore used to store the data for the running average is a self-implmenented concurrent datastore to store the running average,
perhaps the use of a third party in memory dataStore like redis would be more effective if attempting to scale the api, but in the context
of the challenge I felt it was overkill and the use of a simple concurrent data struct is sufficient.

* Another consideration we must account for is the consistent behavior of poller, since there will be variable latency times in getting the response 
from the oden api, we would not be able to reliably call the oden api the same number of times per minute. In one one minute window, oden api may be called
59 times and in the future in another one minute window, there might be 57 calls, which is the nature of real time feed.

### local development environment

* Ubuntu 16.04

* golang
    * current version:
```
$ go version
go version go1.15 linux/amd64
```


### Questions

1. How should we run your solution?



A makefile is created in the working directory (\/avgCableDiameterApi)

my working directory:

```
avwong13@avwong13:~/avgCableDiameterApi$ pwd
/home/avwong13/avgCableDiameterApi
```

run the make run command
```
$ make run
```

In Addition there is a config.json that can be used to switch the server's configuration

* Address specifies the TCP address for the server to listen on, in the form "host:port". If empty, ":http" (port 80) is used.
* ResponseType can be "plain" or "json"
* TimeWindow is in seconds
* File is path to the log, if empty then outputs to stdout
* pollApi is the url for the Get Api to retrieve the cableDiameter values

```json
{
    "Address"        : "0.0.0.0:8080",
    "ReadTimeout"    : 10,
    "WriteTimeout"   : 600,
    "Static"         : "public",
    "pollApi"        : "http://takehome-backend.oden.network/?metric=cable-diameter",
    "File"           : "",
    "TimeWindow"     : 60,
    "ResponseType"   : "plain"
}
```

to call the API use curl in the terminal or use web browser
```
$ curl http://0.0.0.0:8080/cable-diameter
{ 
  "Value": 10.768100226094273
}
```
if the ResponseType is set to plain then we get:
```
$ curl http://0.0.0.0:8080/cable-diameter
10.592566
```

sample run of server locally
```
$ make run
make build 
make[1]: Entering directory '/home/avwong13/avgCableDiameterApi'
make clean 
make[2]: Entering directory '/home/avwong13/avgCableDiameterApi'
rm -rf server 
make[2]: Leaving directory '/home/avwong13/avgCableDiameterApi'
go build -o ./server ./cmd/server/ 
make[1]: Leaving directory '/home/avwong13/avgCableDiameterApi'
./server 
AvgCableDiameter Web Service started at : 0.0.0.0:8080
polledApi Value: 8.337609383825345
sum:  8.337609383825345  numCount: 1  movingAverage:  8.337609383825345
polledApi Value: 8.735794584816379
sum:  17.073403968641724  numCount: 2  movingAverage:  8.536701984320862
polledApi Value: 10.419738121229939
sum:  27.49314208987166  numCount: 3  movingAverage:  9.164380696623887
polledApi Value: 11.220133863656503
sum:  38.71327595352817  numCount: 4  movingAverage:  9.678318988382042

```
2. How long did you spend on the take home? What would you add to your solution if you had more time and expected it to be used in a production setting?

The design doc and solution took roughly 3 hours. Though, documenting and error cases was done throughout the following day which add a couple of hours to the take home.

3. If you used any libraries not in the language’s standard library, why did you use them?

I did not use any libraries outside the standard library for golang, just to keep things simple

4. If you have any feedback, feel free to share your thoughts!