# License BSD  

http://www.linfo.org/bsdlicense.html

# Description

The circleme program will brute force crack the pin code on a Circle from Disney
https://meetcircle.com/. An extensive blog on hacking the Circle from adding a serial
port to disecting the firmeware byte by byte: http://www.hackspot.net/CircleBlog/ 

# Building
## Fedora 23 ##
```
sudo dnf install golang
mkdir -p $HOME/gopath
export GOPATH=$HOME/gopath
go get -u github.com/ndragon70/circleme
cd $HOME/gopath/src/github.com/ndragon70/circleme
go build circleme.go 
./circleme <ip>
```
