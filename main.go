package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"transparent-tcp-proxy/proxy"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	logFile, err := os.OpenFile("proxy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	log.Out = logFile

	if len(os.Args) != 3 {
		log.Fatalf("usage: %s <listen ip:port> <target ip:port>", os.Args[0])
	}

	listenAddress := os.Args[1]
	targetAddress := os.Args[2]

	if listenAddress == targetAddress {
		log.Fatalf("listen and target addresses are the same")
	}

	log.Infof("starting proxy on %s, forwarding to %s", listenAddress, targetAddress)

	p, err := proxy.New(log, listenAddress, targetAddress)
	if err != nil {
		log.Fatalf("failed to create proxy: %v", err)
	}

	if err := p.Run(); err != nil {
		log.Fatalf("failed to run proxy: %v", err)
	}
}
