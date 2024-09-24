# transparent-tcp-proxy

A simple TCP proxy written in Go

## Usage

```
./ttp <listen IP:port> <target IP:port>
```
```
./ttp 0.0.0.0:8080 192.168.52.1:8080
```

## Compile

```
go build -o ttp main.go
```