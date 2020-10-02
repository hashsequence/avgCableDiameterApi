# AvgCableDiameterApi

## Challenge

At Oden we deal with a lot of time-series data. We process millions of IoT metrics every day and

have to make sure our end users can visualize the data in a timely manner. One common use

case for our users is to monitor metrics that are indicators of product quality in real-time. For our

cable extrusion customers, this happens to be the cable diameter.

http://takehome-backend.oden.network/?metric=cable-diameter

The API returns the current value for the metric.

Your task is to assemble an application that polls the Oden API, calculates a one minute moving

average of the Cable Diameter, and exposes that moving average via an HTTP service when

issued a GET request to localhost:8080/cable-diameter.

Example:

```
curl localhost:8080/cable-diameter
11.24
```

Your moving average should be updated, if possible, once per second and, after each update,

the new moving average is logged to STDOUT.