package api

// func (p *bitsharesAPI) Broadcast(wifKeys []string, feeAsset objects.GrapheneObject, ops ...objects.Operation) (string, error) {

// 	operations := objects.Operations(ops)
// 	fees, err := p.GetRequiredFees(operations, feeAsset)
// 	if err != nil {
// 		return "", errors.Annotate(err, "GetRequiredFees")
// 	}

// 	if err := operations.ApplyFees(fees); err != nil {
// 		return "", errors.Annotate(err, "ApplyFees")
// 	}

// 	props, err := p.GetDynamicGlobalProperties()
// 	if err != nil {
// 		return "", errors.Annotate(err, "GetDynamicGlobalProperties")
// 	}

// 	tx, err := objects.NewTransactionWithBlockData(props)
// 	if err != nil {
// 		return "", errors.Annotate(err, "NewTransaction")
// 	}

// 	tx.Operations = operations

// 	pubKeys, err := p.GetPotentialSignatures(tx)
// 	if err != nil {
// 		return "", errors.Annotate(err, "GetPotentialSignatures")
// 	}

// 	util.DumpJSON("potential pubkeys >", pubKeys)

// 	if err := tx.Sign(wifKeys, p.chainConfig.Id()); err != nil {
// 		return "", errors.Annotate(err, "Sign")
// 	}

// 	util.DumpJSON("tx >", tx)

// 	resp, err := p.BroadcastTransaction(tx)
// 	if err != nil {
// 		return "", errors.Annotate(err, "BroadcastTransaction")
// 	}

// 	return resp, err
// }
