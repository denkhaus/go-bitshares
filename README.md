# bitshares

A Bitshares API consuming a websocket connection to an active full node or a RPC connection to your `cli_wallet`. 
Look for several examples in [examples](/examples) and [tests](/tests) folder. This is work in progress and may have breaking changes. 
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
If you change types or operations you have to regenerate the required static `MarshalJSON` and `UnmarshalJSON` functions for the new/changed code.

```
make generate
```
If you encounter any errors, try: 

```
make generate_new
```
to generate ffjson helpers from scratch.


## generate operation samples
To generate op samples for testing, go to [gen](/gen) package.
Generated operation samples get injected automatically while running operation tests.

## testing

To test this stuff I use a combined Docker based MainNet/TestNet wallet, you can find [here](https://github.com/denkhaus/bitshares-docker).
Operations testing uses generated real blockchain sample code by [gen](/gen) package. To test run:
// func (suite *operationsAPITest) Test_WithdrawPermissionUpdateOperation() {
// 	suite.OpSamplesTest(&WithdrawPermissionUpdateOperation{})
// }
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
wsFullApiUrl := "wss://bitshares.openledger.info/ws"

api := bitshares.NewWebsocketAPI(wsFullApiUrl)
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

If you need wallet functions, use:
```
rpcApiUrl    := "http://localhost:8095" 
api := bitshares.NewWalletAPI(rpcApiUrl)

if err := api.connect(); err != nil{
	log.Fatal(err)
}

...
```

For a long application lifecycle, you can use an API instance with latency tester that connects to the most reliable node.
Note: Because the tester takes time to unleash its magic, use the above-mentioned constructor for quick in and out.

```
wsFullApiUrl := "wss://bitshares.openledger.info/ws"

//wsFullApiUrl serves as "quick startup" fallback endpoint here, 
//until the latency tester provides the first results.

api, err := bitshares.NewWithAutoEndpoint(wsFullApiUrl)
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
- [x] OperationTypeFillOrder (virtual)
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
- [x] OperationTypeWithdrawPermissionCreate              
- [ ] OperationTypeWithdrawPermissionUpdate              
- [ ] OperationTypeWithdrawPermissionClaim               
- [x] OperationTypeWithdrawPermissionDelete              
- [ ] OperationTypeCommitteeMemberCreate                 
- [ ] OperationTypeCommitteeMemberUpdate                 
- [x] OperationTypeCommitteeMemberUpdateGlobalParameters 
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