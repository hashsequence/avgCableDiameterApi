package main 

import (
	routes "github.com/hashsequence/avgCableDiameterApi/pkg/routes"
	utils "github.com/hashsequence/avgCableDiameterApi/pkg/utils"
	ds "github.com/hashsequence/avgCableDiameterApi/pkg/dataStore"
	poll "github.com/hashsequence/avgCableDiameterApi/pkg/poll"
	"time"
	"net/http"
)

//sample web server
func main() {
	//load default configuration for web server
	config := utils.LoadConfig("")
	//instantiate dataStore for moving average
	dataStore := ds.NewDataStore()
	//create a central logger
	logger := utils.CreateLogger(config.File)
	//instantiate handler for "/cable-diameter" route
	avgCableDiameterApiHandler := routes.NewGetAverageHandler(dataStore, logger, config.ResponseType)
	//instantiate poller
	poller := poll.NewPoll(config.PollApi, time.Duration(config.TimeWindow) * time.Second, dataStore, logger)
	//since I only have one route "/cable-diameter" I don't need a router
	//however if I were to add more routes I woul replace the Handler with a *http.ServeMux to handle multiple routes
	//there are also popular router libraries like gorilla/mux and go-chi, but the simplicity of this challenge does not
	//require it
	//I have left the commented out use Go's default router/multiplexer if one chooses to add more routes
	//mux := http.NewServeMux()
	s := &http.Server{
		Addr:           config.Address,
		Handler : routes.RecoveryMiddleWare(avgCableDiameterApiHandler), //uncomment and replace value if want to use multiplexer //mux
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	//uncomment if want to use multiplexer 
	//mux.Handle("/", routes.RecoveryMiddleWare(&routes.IndexHandler{}))
	//mux.Handle("/cable-diameter", routes.RecoveryMiddleWare(avgCableDiameterApiHandler))
	utils.InitMsg("AvgCableDiameter Web Service", config.Address)
	//start polling
	poller.Start()
	//listen and server on tcp based on custom server configurations
	s.ListenAndServe()
	//sample use of tls with pre-generated self-signed certificates
    //run curl -k <host>:<port> to call api if hosted locally
    //this.Server.ListenAndServeTLS("ssl/server-cert.pem","ssl/server-key.pem")
}