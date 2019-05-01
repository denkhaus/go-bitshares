package types

//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_assetid.go gen "T1=Asset"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_accountid.go gen "T1=Account"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_forcesettlementid.go gen "T1=ForceSettlement"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_committeememberid.go gen "T1=CommitteeMember"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_witnessid.go gen "T1=Witness"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_limitorderid.go gen "T1=LimitOrder"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_callorderid.go gen "T1=CallOrder"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_customid.go gen "T1=Custom"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_proposalid.go gen "T1=Proposal"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_operationhistoryid.go gen "T1=OperationHistory"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_withdrawpermissionid.go gen "T1=WithdrawPermission"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_vestingbalanceid.go gen "T1=VestingBalance"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_workerid.go gen "T1=Worker"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_balanceid.go gen "T1=Balance"

//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_globalpropertyid.go gen "T1=GlobalProperty"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_dynamicglobalpropertyid.go gen "T1=DynamicGlobalProperty"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_assetdynamicdataid.go gen "T1=AssetDynamicData"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_assetbitassetdataid.go gen "T1=AssetBitAssetData"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_accountbalanceid.go gen "T1=AccountBalance"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_accountstatisticsid.go gen "T1=AccountStatistics"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_transactionid.go gen "T1=Transaction"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_blocksummaryid.go gen "T1=BlockSummary"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_accounttransactionhistoryid.go gen "T1=AccountTransactionHistory"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_blindedbalanceid.go gen "T1=BlindedBalance"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_chainpropertyid.go gen "T1=ChainProperty"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_witnessscheduleid.go gen "T1=WitnessSchedule"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_budgetrecordid.go gen "T1=BudgetRecord"
//go:generate genny -in=../gen/templates/objectid.go.tmpl -out=./gen_specialauthorityid.go gen "T1=SpecialAuthority"

//go:generate stringer -type=OperationType
//go:generate stringer -type=ObjectType
//go:generate stringer -type=AssetType
//go:generate stringer -type=SpaceType
//go:generate stringer -type=AssetPermission
