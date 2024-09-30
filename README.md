# transparent-tcp-proxy

A simple TCP proxy written in Go

## Usage

```
./transparent-tcp-proxy <listen IP:port> <target IP:port>
```
```
./transparent-tcp-proxy 0.0.0.0:8080 192.168.52.1:8080
```

## Compile

```
go build -o transparent-tcp-proxy main.go
```
