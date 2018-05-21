# bitshares

A Bitshares API consuming a websocket connection to an active full node and, if you need wallet functions, an optional RPC connection to your local `cli_wallet`. 
Look for several examples in tests. This is work in progress and may have breaking changes as development goes on. Use it on your own risk. 


## install
```
go get -u github.com/denkhaus/bitshares
```

Install dependencies like [secp256k1](https://github.com/bitcoin-core/secp256k1) with
```
make init
```


This API uses [ffjson](https://github.com/pquerna/ffjson). If you change this code you have to regenerate the required static `MarshalJSON` and `UnmarshalJSON` functions for all API-structures with

```
make generate
```

## code
```
rpcApiUrl    := "http://localhost:8095" 
wsFullApiUrl := "wss://bitshares.openledger.info/ws"

//Note: The RPC endpoint is optional. If you do not need wallet functions
//pass an empty string as second parameter.

api := api.New(wsFullApiUrl, rpcApiUrl)
if err := api.Connect(); err != nil {
	log.Fatal(err)
}

api.OnError(func(err error) {
	log.Fatal(err)
})

UserID   := types.NewGrapheneID("1.2.253") 
AssetBTS := types.NewGrapheneID("1.3.0") 

res, api.GetAccountBalances(UserID, AssetBTS)
if err != nil {
	log.Fatal(err)
}

log.Printf("balances: %v", res)

```

## todo
- add missing operations
- sign transactions by websocket api


Have fun and feel free to contribute.