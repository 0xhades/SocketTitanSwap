# SocketTitanSwap
an instagram swapper, written in go, uses the tcp client to send requests

# installation
```
git clone https://github.com/0xhades/SocketTitanSwap.git ~/go/src/SocketTitanSwap
```
# libraries
```
go get github.com/fatih/color
go get github.com/Pallinder/go-randomdata
```
# run
```
go run SocketTitanSwap
```
# build
```
#local OS
go build -ldflags "-s -w" -o socketTitanSwap SocketTitanSwap

#Windows
GOOS=windows arch=amd64 go build -ldflags "-s -w" -o socketTitanSwap SocketTitanSwap

#linux
GOOS=linux arch=386 go build -ldflags "-s -w" -o socketTitanSwap SocketTitanSwap
```
# about
instagram: @0xhades @ctpe
