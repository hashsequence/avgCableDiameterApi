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
the subset of the polled values in our one minute interval prior to the current time in real-time. Also,since the API only handles one route "/cable-diameter" there is no need for a router/multiplexer.

breakdown of responsibilities: 

* poller
    * Responsible for polling the oden's Api

* dataStore 
    * data structure to store moving average value and related data

* Routes
    * defines Handlers to handle endpoints

* logger
    * maintains the logging of web service

### Scope and Contentious Issues

* The specification did not specify authentication nor encryption requirements and is not part of the core challenge, though if implemented the certificates can be created and signed using openssl, I have taken the liberty to create the certificates for server side authentication and client side authentication,
and a bash script is available in the ssl folder (genCerts.sh) to create said certificates (self-signed), the following resources would be useful in enabling tls for web service: 


    https://golang.org/pkg/net/http/#Server.ListenAndServeTLS 
    
    "crypto/tls"

    "crypto/x509"


Though in my implementation of the web server the we will be using ListenAndServe and will be using http

* The DataStore used to store the data for the running average is a self-implmenented concurrent datastore to store the running average,
perhaps the use of a third party in memory dataStore like redis would be more effective if attempting to scale the api, but in the context
of the challenge I felt it was overkill and the use of a simple concurrent data struct is sufficient.

* Another consideration we must account for is the consistent behavior of poller, since there will be variable latency times in getting the response 
from the oden api, every time we restart the server we may end up with 59,58,or 60 values in the buffer for the first minute. After the minute, we 
will start popping the oldest value every second. For Example, we have 59 values in the buffer after the first minute since starting. If in one interval of time it takes 5 seconds to poll a new value, then 5 old values will be popped, and so we end up with 55 in the buffer. I chose to have a separate goroutine 
for polling new values and popping old values just so the two processes are independent and don't intefere with one another and to ensure that the buffer only contains values that are within the last one minute.

* As for the format of response in the curl example, the value looks like plainText up to 2 significant figures. However, it was not specified what type of response the http web api will serve nor the response value has to be 2 significate figures, so I implemented a conigurable option to respond with a json or a plainText, and the format of value will be float64 to be consitent with the response of the oden Api

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

#### Linux (Ubuntu 16.04)
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
    "Address"        : "localhost:8080",
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
$ curl http://localhost:8080/cable-diameter
{"Value":10.915435784411157}
```
if the ResponseType is set to plain then we get:
```
$ curl http://localhost:8080/cable-diameter
10.772065570961349
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
AvgCableDiameter Web Service started at : localhost:8080
polledApi Value: 12.976316132644891
sum: 12.976316132644891 numCount: 1 movingAverage: 12.976316132644891
polledApi Value: 13.09765038624688
sum: 26.073966518891773 numCount: 2 movingAverage: 13.036983259445886
polledApi Value: 12.451448141041464
sum: 38.52541465993323 numCount: 3 movingAverage: 12.841804886644411
polledApi Value: 11.507305131582514
sum: 50.032719791515746 numCount: 4 movingAverage: 12.508179947878936
polledApi Value: 9.046249575430384
sum: 59.078969366946126 numCount: 5 movingAverage: 11.815793873389225
polledApi Value: 8.51331886708611
sum: 67.59228823403224 numCount: 6 movingAverage: 11.265381372338707
polledApi Value: 8.790965032812514
sum: 76.38325326684475 numCount: 7 movingAverage: 10.911893323834963
polledApi Value: 10.286953565353977
sum: 86.67020683219873 numCount: 8 movingAverage: 10.833775854024841
polledApi Value: 11.473288019372903
sum: 98.14349485157163 numCount: 9 movingAverage: 10.904832761285737
polledApi Value: 12.534110013650851
sum: 110.67760486522248 numCount: 10 movingAverage: 11.067760486522248
...
popped:  12.976316132644891
polledApi Value: 13.07753015906396
sum: 636.6589871983313 numCount: 59 movingAverage: 10.79083029149714
polledApi Value: 13.036937594707014
sum: 649.6959247930383 numCount: 60 movingAverage: 10.828265413217306
popped:  13.09765038624688
popped:  12.451448141041464
polledApi Value: 11.548358014454482
sum: 635.6951842802044 numCount: 59 movingAverage: 10.774494648817024
popped:  11.507305131582514
polledApi Value: 10.40057525450755
sum: 634.5884544031295 numCount: 59 movingAverage: 10.75573651530728
popped:  9.046249575430384
polledApi Value: 8.833609016249282
sum: 634.3758138439484 numCount: 59 movingAverage: 10.752132438033025
popped:  8.51331886708611
polledApi Value: 8.338822411480074
sum: 634.2013173883424 numCount: 59 movingAverage: 10.749174870988854
...
```

#### Windows

* Install the latest version of Go

* Go to Command Prompt

* Run following command in working directory:

    ```
    go build -o server.exe cmd/server/main.go
    ```

* Run server.exe

    ```
    server.exe
    ```

* go to web browser and goto localhost:8080/cable-diameter

2. How long did you spend on the take home? What would you add to your solution if you had more time and expected it to be used in a production setting?

The design doc and solution took roughly 3 hours. Though, documenting and error cases was done throughout the following day which add a couple of hours to the take home.

3. If you used any libraries not in the language’s standard library, why did you use them?

I did not use any libraries outside the standard library for golang, since the challenge was not too complicated. However, I did use the assert library to
handle assertions for testing.

4. If you have any feedback, feel free to share your thoughts!

Had a lot of fun working on this project! 

## Project Layout

```
.
├── cmd
│   └── server
│       └── main.go
├── config.json
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── Oden Software Engineer Takehome.pdf
├── pkg
│   ├── dataStore
│   │   ├── dataStore.go
│   │   └── dataStore_test.go
│   ├── poll
│   │   ├── poll.go
│   │   └── poll_test.go
│   ├── routes
│   │   ├── routes.go
│   │   └── routes_test.go
│   └── utils
│       └── utils.go
├── README.md
├── server
└── ssl
    ├── ca-cert.pem
    ├── ca-cert.srl
    ├── ca-key.pem
    ├── client-cert.pem
    ├── client-ext.cnf
    ├── client-key.pem
    ├── client-req.pem
    ├── genCerts.sh
    ├── server-cert.pem
    ├── server-ext.cnf
    ├── server-key.pem
    └── server-req.pem

```
## Creating and Running Docker Image

in working directory run:
```
$ sudo docker build -t avgcablediameterapi .
[sudo] password for avwong13: 
Sending build context to Docker daemon  33.09MB
Step 1/7 : FROM golang
 ---> 9f495162f677
Step 2/7 : RUN mkdir ./avgCableDiameterApi
 ---> Running in ed8960b534dd
Removing intermediate container ed8960b534dd
 ---> dc798d43ac25
Step 3/7 : WORKDIR ./avgCableDiameterApi
 ---> Running in 88ee49672ddc
Removing intermediate container 88ee49672ddc
 ---> 04110030aede
Step 4/7 : COPY . .
 ---> d437028b2329
Step 5/7 : RUN make build
 ---> Running in f117e53df5b5
make clean 
make[1]: Entering directory '/go/avgCableDiameterApi'
rm -rf server 
make[1]: Leaving directory '/go/avgCableDiameterApi'
go build -o ./server ./cmd/server/ 
Removing intermediate container f117e53df5b5
 ---> b1a9b78bc2e5
Step 6/7 : ENTRYPOINT ["./server"]
 ---> Running in 9c0fe45e2dd6
Removing intermediate container 9c0fe45e2dd6
 ---> def5979402a3
Step 7/7 : EXPOSE 8080
 ---> Running in a354e3eebc0e
Removing intermediate container a354e3eebc0e
 ---> d7ba6755bc9e
Successfully built d7ba6755bc9e
Successfully tagged avgcablediameterapi:latest

```

run docker image exposing on host port 8000 from docker containers port 8080

```
$ sudo docker container run --rm -p 8000:8080 avgcablediameterapi

```

## Testing
Though testing was not part of the challenge, I took the liberty of making test cases with 
Go standard testing libraries

here is a sample run using main test:

```
$ make test
go test -v ./...
?   	github.com/hashsequence/avgCableDiameterApi/cmd/server	[no test files]
=== RUN   TestNewDataStore
sum:  324.626  numCount: 1  movingAverage:  324.626
sum:  324949.252  numCount: 2  movingAverage:  162474.626
sum:  329616.85199999996  numCount: 3  movingAverage:  109872.28399999999
sum:  329617.31399999995  numCount: 4  movingAverage:  82404.32849999999
sum:  350941.93999999994  numCount: 5  movingAverage:  70188.38799999999
popping 324.626
popping 324624.626
popping 4667.6
popping 0.462
popping 21324.626
--- PASS: TestNewDataStore (0.00s)
PASS
ok  	github.com/hashsequence/avgCableDiameterApi/pkg/dataStore	(cached)
=== RUN   TestPoller
polledApi Value: 8.606891976231747
sum:  8.606891976231747  numCount: 1  movingAverage:  8.606891976231747
polledApi Value: 8.30000178079741
sum:  16.90689375702916  numCount: 2  movingAverage:  8.45344687851458
polledApi Value: 9.202394541700805
sum:  26.109288298729965  numCount: 3  movingAverage:  8.703096099576655
polledApi Value: 10.653421271966947
sum:  36.76270957069691  numCount: 4  movingAverage:  9.190677392674228
polledApi Value: 11.622196723468107
sum:  48.38490629416502  numCount: 5  movingAverage:  9.676981258833004
polledApi Value: 13.089514203367397
sum:  61.47442049753242  numCount: 6  movingAverage:  10.245736749588737
polledApi Value: 12.577677435877836
sum:  74.05209793341025  numCount: 7  movingAverage:  10.578871133344322
polledApi Value: 10.251780080680044
sum:  84.3038780140903  numCount: 8  movingAverage:  10.537984751761288
polledApi Value: 11.717770338364097
sum:  96.0216483524544  numCount: 9  movingAverage:  10.6690720391616
--- PASS: TestPoller (16.00s)
PASS
ok  	github.com/hashsequence/avgCableDiameterApi/pkg/poll	(cached)
=== RUN   TestCableDiameterRouteJsonResponse
polledApi Value: 11.014617821958128
sum:  11.014617821958128  numCount: 1  movingAverage:  11.014617821958128
polledApi Value: 9.56807397672658
sum:  20.582691798684706  numCount: 2  movingAverage:  10.291345899342353
polledApi Value: 8.492795949688682
sum:  29.07548774837339  numCount: 3  movingAverage:  9.691829249457797
currentAverage: 9.691829249457797
--- PASS: TestCableDiameterRouteJsonResponse (5.00s)
=== RUN   TestCableDiameterRoutePLainResponse
polledApi Value: 8.387346659929147
sum:  37.46283440830254  numCount: 4  movingAverage:  9.365708602075635
polledApi Value: 9.360237891446603
sum:  46.82307229974914  numCount: 5  movingAverage:  9.364614459949829
polledApi Value: 10.565137947341825
sum:  57.38821024709097  numCount: 6  movingAverage:  9.564701707848494
polledApi Value: 11.045150351770614
sum:  11.045150351770614  numCount: 1  movingAverage:  11.045150351770614
polledApi Value: 11.165611128426463
sum:  68.55382137551743  numCount: 7  movingAverage:  9.793403053645347
polledApi Value: 11.358646736541402
sum:  22.403797088312018  numCount: 2  movingAverage:  11.201898544156009
polledApi Value: 12.141072084142653
sum:  80.69489345966008  numCount: 8  movingAverage:  10.08686168245751
polledApi Value: 12.87168753334768
sum:  35.2754846216597  numCount: 3  movingAverage:  11.758494873886567
currentAverage: 11.758494873886567
plaintext response:  11.758495 type:  float64
--- PASS: TestCableDiameterRoutePLainResponse (5.00s)
PASS
ok  	github.com/hashsequence/avgCableDiameterApi/pkg/routes	10.006s
?   	github.com/hashsequence/avgCableDiameterApi/pkg/utils	[no test files]
```