module github.com/ElrondNetwork/notifier-go

go 1.16

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/ElrondNetwork/elrond-go v1.2.4-0.20210625084351-7915cdd77085
	github.com/gin-gonic/gin v1.7.2
	github.com/google/uuid v1.2.0
	github.com/gorilla/websocket v1.4.2
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.7.0
	github.com/vektah/gqlparser/v2 v2.1.0
)

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.19 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.19
