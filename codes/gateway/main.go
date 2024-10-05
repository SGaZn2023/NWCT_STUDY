package main

import (
	"flag"
	"nwct_st/common"
	httprouter "nwct_st/gateway/http_router"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "c", "", "config file")
	flag.Parse()

	conf, eerr := ParseConfig(confFile)
	if eerr != nil {
		panic(eerr)
	}

	listenerConfigs, err := ParseListenerConfig(conf.ListenerFile)
	if err != nil {
		panic(err)
	}

	sessionMgr := NewSessionManager()
	httpRouter := httprouter.NewApisixRouter(&httprouter.ApisixConfig{
		Api: "http://127.0.0.1:9180",
		Key: "123",
	})

	for _, listenerConfig := range listenerConfigs {
		listener := NewListener(&common.ProxyProtocol{
			ClientID:         listenerConfig.ClientID,
			PublicProtocol:   listenerConfig.PublicProtocol,
			PublicIP:         listenerConfig.PublicIP,
			PublicPort:       listenerConfig.PublicPort,
			InternalProtocol: listenerConfig.InternalProtocol,
			InternalIP:       listenerConfig.InternalIP,
			InternalPort:     listenerConfig.InternalPort,
		}, sessionMgr, httpRouter)

		go func() {
			defer listener.Close()
			err := listener.ListenAndServer()
			if err != nil {
				panic(err)
			}
		}()
	}

	gw := NewGateway(":5103", sessionMgr)
	err = gw.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
