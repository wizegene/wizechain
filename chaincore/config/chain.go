package config

var (
	Chaincode            = []byte("6A2322D4C28884FCEC5B242E85282BE16D67D139087D219A7BFFE76E5D4E010F")
	Prefix               = []byte("6A23")
	ChainName            = "Wizechain"
	ChainVersion         = "v0.0.1"
	Currency             = "WIZECOIN"
	Symbol               = "WIZE"
	BadgerDBDir          = "./db"
	MinValidators        = 3
	MaxValidators        = 4999
	MinTestingValidators = 3
	Corum                = 0.33
	// 0.01 wize as reward per winning vote 100 winning votes = 1 wize
	StartRewardPerValidator  = 0.0001
	StartRewardAsWize        = 0.01
	RewardInflationPerSecond = 0.00000201873385

	// MaxBlockPerYear 10 million
	MaxBlockPerYear            = 10000000
	CoinsPerBlock              = 100
	MaxCoinsPerYear            = MaxBlockPerYear * CoinsPerBlock
	SecondsBetweenBlocks       = 31536000 / MaxBlockPerYear
	StartingValuePerWizeInFiat = 0.01
	WebsiteDomain              = "wizechain.com"
)
