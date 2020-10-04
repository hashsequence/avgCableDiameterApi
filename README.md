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
    * will be pushing the new values into a dataStore and popping values older than a minute from the dataStore, frequeuncy of updates is every second

* dataStore 
    * data structure to store moving average value and related data
    * will be using a queue implemented with a ring buffer to store the values in the one minute interval
        * reasons for using ring buffer over slices + append or linked-lists
        * size we know the interval is going to be a minute, we only need to store values up to a minute old, thus the buffer is fixed
        * with a ring buffer we don't need to reallocate memory for new elements since we are reusing the same space, whereas slices and linkeded-list implementation will need to allocate new memory for each pop (linked list allocates new memory for add)
        * the number of garbage collection cycles will decrease due to the fact we are not allocating as much for new puses and pops every second

* Routes
    * defines Handlers to handle endpoints
    * only need to handle one route which is "/cable-diameter"
    * we don't need a router/multiplexer such as http.ServeMux, go-chi, or gorilla/mux since there is only one route

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

### Local Development Environment

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
Started Polling
polledApi Value: 11.74917035760581
sum: 11.74917035760581 numCount: 1 movingAverage: 11.74917035760581
polledApi Value: 11.030011904135948
sum: 22.77918226174176 numCount: 2 movingAverage: 11.38959113087088
polledApi Value: 9.106767232155976
sum: 31.885949493897733 numCount: 3 movingAverage: 10.628649831299244
polledApi Value: 8.350898441419519
sum: 40.23684793531725 numCount: 4 movingAverage: 10.059211983829313
polledApi Value: 8.332320103055508
sum: 48.56916803837276 numCount: 5 movingAverage: 9.713833607674552
polledApi Value: 9.1249480204555
sum: 57.69411605882826 numCount: 6 movingAverage: 9.61568600980471
polledApi Value: 11.106765960712716
sum: 68.80088201954098 numCount: 7 movingAverage: 9.828697431362997
polledApi Value: 12.278551470943949
sum: 81.07943349048493 numCount: 8 movingAverage: 10.134929186310616
polledApi Value: 12.307122550101495
sum: 93.38655604058643 numCount: 9 movingAverage: 10.376284004509603
polledApi Value: 13.078199252325842
sum: 106.46475529291227 numCount: 10 movingAverage: 10.646475529291227
polledApi Value: 12.458113293487468
sum: 118.92286858639974 numCount: 11 movingAverage: 10.811169871490886
polledApi Value: 10.344222940836296
sum: 129.26709152723603 numCount: 12 movingAverage: 10.772257627269669
polledApi Value: 10.103244052644046
sum: 139.37033557988008 numCount: 13 movingAverage: 10.72079504460616
polledApi Value: 8.517363777973337
sum: 147.88769935785342 numCount: 14 movingAverage: 10.56340709698953
polledApi Value: 8.546598747845277
sum: 156.4342981056987 numCount: 15 movingAverage: 10.42895320704658
polledApi Value: 9.08544031110982
sum: 165.51973841680854 numCount: 16 movingAverage: 10.344983651050534
polledApi Value: 10.519748861174364
sum: 176.03948727798291 numCount: 17 movingAverage: 10.355263957528408
polledApi Value: 11.222766020141968
sum: 187.26225329812488 numCount: 18 movingAverage: 10.403458516562493
polledApi Value: 12.855978053028958
sum: 200.11823135115384 numCount: 19 movingAverage: 10.532538492165992
polledApi Value: 12.834033270185007
sum: 212.95226462133886 numCount: 20 movingAverage: 10.647613231066943
polledApi Value: 11.745333159960666
sum: 224.69759778129952 numCount: 21 movingAverage: 10.699885608633311
polledApi Value: 10.170976089069724
sum: 234.86857387036923 numCount: 22 movingAverage: 10.675844266834964
polledApi Value: 9.97489721081071
sum: 244.84347108117993 numCount: 23 movingAverage: 10.645368307877389
polledApi Value: 8.939985419870883
sum: 253.7834565010508 numCount: 24 movingAverage: 10.574310687543784
polledApi Value: 8.300041961674747
sum: 262.0834984627256 numCount: 25 movingAverage: 10.483339938509022
polledApi Value: 9.4147569780949
sum: 271.49825544082046 numCount: 26 movingAverage: 10.44224059387771
polledApi Value: 10.274919843909611
sum: 281.7731752847301 numCount: 27 movingAverage: 10.436043529064078
polledApi Value: 12.359498024362281
sum: 294.13267330909235 numCount: 28 movingAverage: 10.504738332467584
polledApi Value: 13.063820597121442
sum: 307.1964939062138 numCount: 29 movingAverage: 10.59298254849013
polledApi Value: 12.753235446540508
sum: 319.9497293527543 numCount: 30 movingAverage: 10.664990978425143
polledApi Value: 11.807605634346931
sum: 331.7573349871012 numCount: 31 movingAverage: 10.701849515712942
polledApi Value: 11.079594800734023
sum: 342.83692978783523 numCount: 32 movingAverage: 10.713654055869851
polledApi Value: 9.217941949157579
sum: 352.0548717369928 numCount: 33 movingAverage: 10.668329446575541
polledApi Value: 8.71181806148321
sum: 360.766689798476 numCount: 34 movingAverage: 10.610784994072825
polledApi Value: 8.66858760038153
sum: 369.4352773988576 numCount: 35 movingAverage: 10.55529363996736
polledApi Value: 8.792259155737973
sum: 378.22753655459553 numCount: 36 movingAverage: 10.506320459849876
polledApi Value: 10.298090882045571
sum: 388.5256274366411 numCount: 37 movingAverage: 10.500692633422734
polledApi Value: 11.450681376275355
sum: 399.9763088129165 numCount: 38 movingAverage: 10.525692337182013
polledApi Value: 13.098112329968703
sum: 413.0744211428852 numCount: 39 movingAverage: 10.591651824176544
polledApi Value: 12.984499022035477
sum: 426.0589201649207 numCount: 40 movingAverage: 10.651473004123018
polledApi Value: 11.278059889957678
sum: 437.3369800548784 numCount: 41 movingAverage: 10.666755611094594
polledApi Value: 11.01853304772657
sum: 448.35551310260496 numCount: 42 movingAverage: 10.675131264347737
polledApi Value: 8.804741910186264
sum: 457.1602550127912 numCount: 43 movingAverage: 10.631633837506772
polledApi Value: 8.480733762877543
sum: 465.64098877566875 numCount: 44 movingAverage: 10.582749744901562
polledApi Value: 8.313148912335835
sum: 473.9541376880046 numCount: 45 movingAverage: 10.532314170844547
polledApi Value: 10.066055721958
sum: 484.0201934099626 numCount: 46 movingAverage: 10.522178117607883
polledApi Value: 10.5322244313521
sum: 494.5524178413147 numCount: 47 movingAverage: 10.522391868964142
polledApi Value: 12.776293576668992
sum: 507.3287114179837 numCount: 48 movingAverage: 10.569348154541327
polledApi Value: 13.094465146663474
sum: 520.4231765646472 numCount: 49 movingAverage: 10.620881154380555
polledApi Value: 12.716142501245509
sum: 533.1393190658927 numCount: 50 movingAverage: 10.662786381317853
polledApi Value: 11.585664201001949
sum: 544.7249832668946 numCount: 51 movingAverage: 10.68088202484107
polledApi Value: 9.623950678455087
sum: 554.3489339453497 numCount: 52 movingAverage: 10.660556422025955
polledApi Value: 9.420427453316755
sum: 563.7693613986664 numCount: 53 movingAverage: 10.637157762238989
polledApi Value: 8.321599028648585
sum: 572.090960427315 numCount: 54 movingAverage: 10.594277044950278
polledApi Value: 8.747583847306371
sum: 580.8385442746214 numCount: 55 movingAverage: 10.560700804993116
polledApi Value: 10.186799368149547
sum: 591.025343642771 numCount: 56 movingAverage: 10.55402399362091
polledApi Value: 10.447597710567551
sum: 601.4729413533386 numCount: 57 movingAverage: 10.552156865848046
polledApi Value: 12.500368083296886
sum: 613.9733094366354 numCount: 58 movingAverage: 10.585746714424749
polledApi Value: 12.999975284663442
sum: 626.9732847212989 numCount: 59 movingAverage: 10.62666584273388
polledApi Value: 12.548803801632475
sum: 639.5220885229314 numCount: 60 movingAverage: 10.658701475382191
popped:  11.74917035760581
polledApi Value: 12.059964366615977
sum: 639.8328825319416 numCount: 60 movingAverage: 10.663881375532359
popped:  11.030011904135948
polledApi Value: 10.721856367793523
sum: 639.5247269955992 numCount: 60 movingAverage: 10.658745449926652
popped:  9.106767232155976
polledApi Value: 8.834368133877568
sum: 639.2523278973207 numCount: 60 movingAverage: 10.654205464955345
popped:  8.350898441419519
popped:  8.332320103055508
polledApi Value: 8.34697431833056
sum: 630.9160836711762 numCount: 59 movingAverage: 10.693492943579258
polledApi Value: 8.800808267422983
sum: 639.7168919385992 numCount: 60 movingAverage: 10.661948198976654
popped:  9.1249480204555
polledApi Value: 9.468856380509935
sum: 640.0608002986537 numCount: 60 movingAverage: 10.667680004977562
popped:  11.106765960712716
polledApi Value: 11.578832339846894
sum: 640.532866677788 numCount: 60 movingAverage: 10.675547777963134
popped:  12.278551470943949
polledApi Value: 12.300727657630265
sum: 640.5550428644743 numCount: 60 movingAverage: 10.675917381074573
popped:  12.307122550101495
polledApi Value: 13.009596489410697
sum: 641.2575168037836 numCount: 60 movingAverage: 10.68762528006306
popped:  13.078199252325842
popped:  12.458113293487468
polledApi Value: 12.08204087822432
sum: 627.8032451361946 numCount: 59 movingAverage: 10.640732968410079
popped:  10.344222940836296
polledApi Value: 10.08052147610145
sum: 627.5395436714598 numCount: 59 movingAverage: 10.63626345205864
popped:  10.103244052644046
polledApi Value: 9.226872754439519
sum: 626.6631723732553 numCount: 59 movingAverage: 10.621409701241616
popped:  8.517363777973337
polledApi Value: 8.307755105005613
sum: 626.4535637002876 numCount: 59 movingAverage: 10.617857011869281
GetAverageHandler called, currentAverage: 10.617857011869281
404 not found.
GetAverageHandler called, currentAverage: 10.617857011869281
popped:  8.546598747845277
GetAverageHandler called, currentAverage: 10.653568361249008
polledApi Value: 8.450210051728892
sum: 626.3571750041713 numCount: 59 movingAverage: 10.616223305155446
GetAverageHandler called, currentAverage: 10.616223305155446
GetAverageHandler called, currentAverage: 10.616223305155446
GetAverageHandler called, currentAverage: 10.616223305155446
GetAverageHandler called, currentAverage: 10.616223305155446
GetAverageHandler called, currentAverage: 10.616223305155446
GetAverageHandler called, currentAverage: 10.616223305155446
popped:  9.08544031110982
GetAverageHandler called, currentAverage: 10.642616115397612
GetAverageHandler called, currentAverage: 10.642616115397612
polledApi Value: 9.314461114910653
sum: 626.5861958079721 numCount: 59 movingAverage: 10.620105013694442
GetAverageHandler called, currentAverage: 10.620105013694442
polledApi Value: 9.753597672198895
sum: 636.339793480171 numCount: 60 movingAverage: 10.605663224669517
popped:  10.519748861174364
polledApi Value: 10.933835677449
sum: 636.7538802964457 numCount: 60 movingAverage: 10.612564671607428
popped:  11.222766020141968
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

The design doc and solution took roughly 3 hours. Though, documenting and error cases was done throughout the following day which add a couple of hours to the take home. If it was used in a protection setting, I would probably implement encryption and authorization more robustly. For example, if we were using certificates I would use certs signed by an actual certificate authority. Perhaps, I would use existing tools and platforms for in memory datastores, and implement a more robust logging system. I would also do more thorough testing of api, and actually test api when deployed to aws, gcp, heroku, or a chosen cloud provider. I would also consider using terraform and docker technologies to manage infrastructure. Also for scaling, I thought about using nginx for load balancing (if the web service was deployed onto multiple containers, and if so we need to have some form of message queue, and a shared dataStore to
sync up data for the moving average). All in all, most of the considerations have to do with managing infrastructure, scaling, and testing.

3. If you used any libraries not in the language’s standard library, why did you use them?

I did not use any libraries outside the standard library for golang for the implementation, since the challenge was not too complicated. However, I did use the assert library to handle assertions for testing.

4. If you have any feedback, feel free to share your thoughts!

Had a lot of fun working on this project! 

## Project Layout

```
:~/avgCableDiameterApi$ tree
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

here is a sample run using make test:

```
$ make test
go test -v ./...
?   	github.com/hashsequence/avgCableDiameterApi/cmd/server	[no test files]
=== RUN   TestNewDataStore
--- PASS: TestNewDataStore (0.00s)
PASS
ok  	github.com/hashsequence/avgCableDiameterApi/pkg/dataStore	(cached)

=== RUN   TestPoller
Started Polling
polledApi Value: 8.796030117897091
sum: 8.796030117897091 numCount: 1 movingAverage: 8.796030117897091
polledApi Value: 8.303582068667177
sum: 17.09961218656427 numCount: 2 movingAverage: 8.549806093282134
polledApi Value: 9.229916059357251
sum: 26.32952824592152 numCount: 3 movingAverage: 8.776509415307173
polledApi Value: 10.038588037806415
sum: 36.368116283727936 numCount: 4 movingAverage: 9.092029070931984
polledApi Value: 10.86249033173242
sum: 47.23060661546036 numCount: 5 movingAverage: 9.446121323092072
polledApi Value: 12.88562169422485
sum: 60.11622830968521 numCount: 6 movingAverage: 10.019371384947535
Stopped Polling
polledApi Value: 12.971654150359377
sum: 73.0878824600446 numCount: 7 movingAverage: 10.441126065720656
polledApi Value: 12.902304151258843
sum: 85.99018661130344 numCount: 8 movingAverage: 10.74877332641293
Started Polling
polledApi Value: 10.442247878357659
sum: 96.4324344896611 numCount: 9 movingAverage: 10.714714943295679
polledApi Value: 8.890039300229873
sum: 105.32247378989098 numCount: 10 movingAverage: 10.532247378989098
polledApi Value: 8.300832580654122
sum: 121.93261288026488 numCount: 12 movingAverage: 10.161051073355408
polledApi Value: 8.309306509719782
sum: 113.63178029961077 numCount: 11 movingAverage: 10.33016184541916
--- PASS: TestPoller (23.00s)
PASS
ok  	github.com/hashsequence/avgCableDiameterApi/pkg/poll	23.007s
=== RUN   TestCableDiameterRouteJsonResponse
Started Polling
polledApi Value: 8.973460112453658
sum: 8.973460112453658 numCount: 1 movingAverage: 8.973460112453658
polledApi Value: 8.342723853185404
sum: 17.316183965639063 numCount: 2 movingAverage: 8.658091982819531
polledApi Value: 9.50913884846703
sum: 26.825322814106094 numCount: 3 movingAverage: 8.941774271368699
GetAverageHandler called, currentAverage: 8.941774271368699
Alloc = 0 MiB	TotalAlloc = 0 MiB	Sys = 70 MiB	NumGC = 0
--- PASS: TestCableDiameterRouteJsonResponse (5.00s)
=== RUN   TestCableDiameterRoutePLainResponse
Started Polling
polledApi Value: 10.879254439331087
sum: 37.70457725343718 numCount: 4 movingAverage: 9.426144313359295
polledApi Value: 11.198157352166028
sum: 48.90273460560321 numCount: 5 movingAverage: 9.780546921120642
polledApi Value: 13.099968640325857
sum: 13.099968640325857 numCount: 1 movingAverage: 13.099968640325857
polledApi Value: 13.048803744015789
sum: 61.951538349619 numCount: 6 movingAverage: 10.325256391603167
polledApi Value: 13.037235141662322
sum: 26.137203781988177 numCount: 2 movingAverage: 13.068601890994088
polledApi Value: 12.765692477534817
sum: 74.71723082715381 numCount: 7 movingAverage: 10.673890118164831
polledApi Value: 12.235361046712564
sum: 38.37256482870074 numCount: 3 movingAverage: 12.790854942900246
polledApi Value: 11.953711002794469
sum: 86.67094182994828 numCount: 8 movingAverage: 10.833867728743535
polledApi Value: 11.19208778537498
sum: 49.56465261407572 numCount: 4 movingAverage: 12.39116315351893
GetAverageHandler called, currentAverage: 12.39116315351893
Alloc = 1 MiB	TotalAlloc = 1 MiB	Sys = 70 MiB	NumGC = 0
plaintext response:  12.39116315351893 type:  float64
--- PASS: TestCableDiameterRoutePLainResponse (5.00s)
PASS
ok  	github.com/hashsequence/avgCableDiameterApi/pkg/routes	10.009s
?   	github.com/hashsequence/avgCableDiameterApi/pkg/utils	[no test files]
```