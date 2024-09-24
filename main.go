package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"transparent-tcp-proxy/proxy"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	if len(os.Args) != 3 {
		log.Fatalf("usage: %s <listen ip:port> <target ip:port>", os.Args[0])
	}

	listenAddress := os.Args[1]
	targetAddress := os.Args[2]

	if listenAddress == targetAddress {
		log.Fatalf("listen and target addresses are the same")
	}

	log.Infof("starting proxy on %s, forwarding to %s", listenAddress, targetAddress)

	p, err := proxy.New(listenAddress, targetAddress)
	if err != nil {
		log.Fatalf("failed to create proxy: %v", err)
	}

	if err := p.Run(); err != nil {
		log.Fatalf("failed to run proxy: %v", err)
	}
}
