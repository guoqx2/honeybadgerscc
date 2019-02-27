docker run -i --rm  -v $GOPATH/src/github.com/hybridNeo/honeybadgerscc:/opt/gopath/src/github.com/hybridNeo/honeybadgerscc -w /opt/gopath/src/github.com/hybridNeo/honeybadgerscc \
		-v /Users/rahulshivumahadev/go/src/github.com/hyperledger/fabric:/opt/gopath/src/github.com/hyperledger/fabric \
                -v /Users/rahulshivumahadev/go/src/github.com/hybridNeo/honeybadgerscc/.build/docker/bin:/opt/gopath/bin \
                -v /Users/rahulshivumahadev/go/src/github.com/hybridNeo/honeybadgerscc/.build/docker/pkg:/opt/gopath/pkg \
                -v /Users/rahulshivumahadev/go/src/github.com/syndtr/goleveldb:/opt/gopath/src/github.com/syndtr/goleveldb \
                -v /Users/rahulshivumahadev/go/src/github.com/golang/snappy:/opt/gopath/src/github.com/golang/snappy \
                hyperledger/fabric-baseimage:amd64-0.4.14 \
                go build -buildmode=plugin -tags ""
