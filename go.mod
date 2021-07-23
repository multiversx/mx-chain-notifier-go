module github.com/ElrondNetwork/notifier-go

go 1.16

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/ElrondNetwork/elrond-go v1.2.5-0.20210722133055-8c8ab1de09f2
	github.com/ElrondNetwork/elrond-go-core v1.0.1-0.20210721121720-f02fb03b2e1a
	github.com/ElrondNetwork/elrond-go-logger v1.0.5
	github.com/gin-contrib/cors v0.0.0-20190301062745-f9e10995c85a
	github.com/gin-gonic/gin v1.7.2
	github.com/google/uuid v1.2.0
	github.com/gorilla/websocket v1.4.2
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.5
	github.com/vektah/gqlparser/v2 v2.1.0
)

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.27 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.27

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.24 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.24

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.19 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.19

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.26 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.26

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.28 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.28

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_4 v1.4.3 => github.com/ElrondNetwork/arwen-wasm-vm v1.4.3
