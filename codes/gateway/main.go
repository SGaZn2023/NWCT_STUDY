package main

import (
	"flag"
	"nwct_st/common"
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

	for _, listenerConfig := range listenerConfigs {
		listener := NewListener(&common.ProxyProtocol{
			ClientID:         listenerConfig.ClientID,
			PublicProtocol:   listenerConfig.PublicProtocol,
			PublicIP:         listenerConfig.PublicIP,
			PublicPort:       listenerConfig.PublicPort,
			InternalProtocol: listenerConfig.InternalProtocol,
			InternalIP:       listenerConfig.InternalIP,
			InternalPort:     listenerConfig.InternalPort,
		}, sessionMgr)

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
