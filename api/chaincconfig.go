package api

type ChainConfig map[string]interface{}

var (
	knownNetworks = []ChainConfig{
		ChainConfig{
			"name":           "Unknown",
			"core_asset":     "n/a",
			"address_prefix": "n/a",
			"chain_id":       "n/a",
		},
		ChainConfig{
			"name":           "BitShares",
			"core_asset":     "BTS",
			"address_prefix": "BTS",
			"chain_id":       "4018d7844c78f6a6c41c6a552b898022310fc5dec06da467ee7905a8dad512c8",
		},
		ChainConfig{
			"name":           "Muse",
			"core_asset":     "MUSE",
			"address_prefix": "MUSE",
			"chain_id":       "45ad2d3f9ef92a49b55c2227eb06123f613bb35dd08bd876f2aea21925a67a67",
		},
		ChainConfig{
			"name":           "Test",
			"core_asset":     "TEST",
			"address_prefix": "TEST",
			"chain_id":       "39f5e2ede1f8bc1a3a54a7914414e3779e33193f1f5693510e73cb7a87617447",
		},
		ChainConfig{
			"name":           "Obelisk",
			"core_asset":     "GOV",
			"address_prefix": "FEW",
			"chain_id":       "1cfde7c388b9e8ac06462d68aadbd966b58f88797637d9af805b4560b0e9661e",
		},
	}
)

func (p *bitsharesAPI) GetChainConfig(chainID string) (*ChainConfig, error) {
	for _, cnf := range knownNetworks {
		if cnf["chain_id"] == chainID {
			return &cnf, nil
		}
	}

	cnf := knownNetworks[0]
	cnf["chain_id"] = chainID

	return &cnf, nil
}
