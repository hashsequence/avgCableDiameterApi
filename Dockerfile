
FROM golang

RUN mkdir ./avgCableDiameterApi

WORKDIR ./avgCableDiameterApi

COPY . .

RUN make build

ENTRYPOINT ["./server"]

EXPOSE 8080