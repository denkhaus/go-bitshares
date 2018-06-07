# bitshares

A Bitshares API consuming a websocket connection to an active full node. If you need wallet features, specify an optional RPC connection to your local `cli_wallet`. 
Look for several examples in tests. This is work in progress and may have breaking changes. 
No additional cgo dependencies for transaction signing required. 
Use it at your own risk. 

## install

```
go get -u github.com/denkhaus/bitshares
```

Install dev-dependencies with
```
make init
```

This API uses [ffjson](https://github.com/pquerna/ffjson). 
If you change this code you have to regenerate the required static `MarshalJSON` and `UnmarshalJSON` functions for all API-structures with

```
make generate
```

## testing

To test this stuff I use a combined Docker based MainNet/TestNet wallet, you can find [here](https://github.com/denkhaus/bitshares-docker).
Operations testing uses generated real blockchain sample code by gen package. To test run:

```
make test_operations
make test_api
```

or a long running block (deserialize/serialize/compare) range test.

```
make test_blocks
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

For a long application lifecycle, you can use an API instance with latency tester that connects to the most reliable node.
Note: Because the tester takes time to unleash its magic, use the above-mentioned constructor for quick in and out.

```
rpcApiUrl    := "http://localhost:8095" 
wsFullApiUrl := "wss://bitshares.openledger.info/ws"

//Note: The RPC endpoint is optional. If you do not need wallet functions
//pass an empty string as second parameter.

//wsFullApiUrl serves as "quick startup" fallback endpoint here, until the latency tester provides the first results.
api, err := api.NewWithAutoEndpoint(wsFullApiUrl, rpcApiUrl)
if err != nil {
	log.Fatal(err)
}

if err := api.Connect(); err != nil {
	log.Fatal(err)
}

api.OnError(func(err error) {
	log.Fatal(err)
})

...

```

## implemented and tested (serialize/unserialize) operations

- [x] OperationTypeTransfer OperationType
- [x] OperationTypeLimitOrderCreate
- [x] OperationTypeLimitOrderCancel
- [x] OperationTypeCallOrderUpdate
- [x] OperationTypeFillOrder (test failing)
- [x] OperationTypeAccountCreate
- [x] OperationTypeAccountUpdate
- [x] OperationTypeAccountWhitelist
- [x] OperationTypeAccountUpgrade
- [ ] OperationTypeAccountTransfer 
- [x] OperationTypeAssetCreate
- [x] OperationTypeAssetUpdate
- [x] OperationTypeAssetUpdateBitasset
- [x] OperationTypeAssetUpdateFeedProducers
- [x] OperationTypeAssetIssue
- [x] OperationTypeAssetReserve
- [x] OperationTypeAssetFundFeePool
- [x] OperationTypeAssetSettle
- [ ] OperationTypeAssetGlobalSettle 
- [x] OperationTypeAssetPublishFeed
- [x] OperationTypeWitnessCreate
- [x] OperationTypeWitnessUpdate
- [x] OperationTypeProposalCreate
- [x] OperationTypeProposalUpdate
- [x] OperationTypeProposalDelete
- [ ] OperationTypeWithdrawPermissionCreate              
- [ ] OperationTypeWithdrawPermissionUpdate              
- [ ] OperationTypeWithdrawPermissionClaim               
- [ ] OperationTypeWithdrawPermissionDelete              
- [ ] OperationTypeCommitteeMemberCreate                 
- [ ] OperationTypeCommitteeMemberUpdate                 
- [ ] OperationTypeCommitteeMemberUpdateGlobalParameters 
- [x] OperationTypeVestingBalanceCreate
- [x] OperationTypeVestingBalanceWithdraw
- [x] OperationTypeWorkerCreate
- [ ] OperationTypeCustom 
- [ ] OperationTypeAssert 
- [x] OperationTypeBalanceClaim
- [x] OperationTypeOverrideTransfer
- [ ] OperationTypeTransferToBlind   
- [ ] OperationTypeBlindTransfer     
- [ ] OperationTypeTransferFromBlind 
- [ ] OperationTypeAssetSettleCancel 
- [ ] OperationTypeAssetClaimFees    
- [ ] OperationTypeFBADistribute     
- [x] OperationTypeBidColatteral
- [ ] OperationTypeExecuteBid 

## todo
- add missing operations
- add convenience functions 


Have fun and feel free to contribute needed operations and tests.