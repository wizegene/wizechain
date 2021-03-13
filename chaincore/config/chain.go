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
var LocalPrefix4Byte = [4]byte{0x06, 0x41, 0x02, 0x03}
var LivePrefix4Byte = [4]byte{0x06, 0x41, 0x02, 0x01}
var TestPrefix4Byte = [4]byte{0x06, 0x41, 0x02, 0x02}
var CurrencySymbol4Byte = [4]byte{0x57, 0x49, 0x5a, 0x45}

// Script Tags as Bytes

var Script_4Byte_Amount_Prefix = [4]byte{0x77, 0x41, 0x6d, 0x74}
var Script_Cmd_Prefix = [2]byte{0x77, 0x43}
var Script_Elements_Sep = 0x5f
var Script_Tx_FlagByte = 0x7e
var Script_PrePadding = []byte("______")
var Script_End = []byte("_E_")
