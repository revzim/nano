module github.com/revzim/nano

go 1.16

replace github.com/revzim/nano => ./

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/pingcap/check v0.0.0-20200212061837-5e12011dc712
	github.com/pingcap/errors v0.11.4
	github.com/revzim/go-pomelo-client v0.0.1
	github.com/urfave/cli v1.22.5
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)
