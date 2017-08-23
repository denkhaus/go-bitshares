package api

//Sell places a limit order attempting to sell one asset for another.
//This API call abstracts away some of the details of the sell_asset call to be more user friendly.
//All orders placed with sell never timeout and will not be killed if they cannot be filled immediately.
//If you wish for one of these parameters to be different, then sell_asset should be used instead.
/* func (p *bitsharesAPI) Sell(sellerAccount string, base, quote objects.GrapheneObject, rate float64, amount float64, broadcast bool) (*objects.SignedTransaction, error) {

	resp, err := p.rpcClient.CallAPI("sell", sellerAccount, base.Id(), quote.Id(), rate, amount, broadcast)
	if err != nil {
		return nil, errors.Annotate(err, "sell")
	}

	ret := objects.SignedTransaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal SignedTransaction")
	}

	return &ret, nil
}
*/
