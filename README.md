# bitshares

A Bitshares API consuming a RPC connection to your local `cli_wallet` and  a websocket connection to an active full node. 
Look for example code in the tests folder. This is work in progress and may have breaking changes as development goes on. Use it on your own risk. 


## install
```
go get -u github.com/denkhaus/bitshares
```

This API uses [ffjson](https://github.com/pquerna/ffjson). If you change this code you have to regenerate the required static `MarshalJSON` and `UnmarshalJSON` functions for all API-structures with

```
make generate
```
## code
```
rpcApiUrl    := "http://localhost:8095"
wsFullApiUrl := "wss://bitshares.openledger.info/ws"

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