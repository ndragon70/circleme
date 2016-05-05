
Fedora 23

sudo dnf install golang
mkdir -p $HOME/gopath
export GOPATH=$HOME/gopath
go get -u github.com/ndragon70/circleme
cd $HOME/gopath/src/github.com/ndragon70/circleme
go build circleme.go 

./circleme <ip>
