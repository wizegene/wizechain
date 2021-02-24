package consensus

type IConsensus interface {

	StakeCoins(validatorAddress string, toAddress string, coins []Coin) (error, bool)
	UnstakeCoins(validatorAddress string, fromAddress string, coins []Coin) (error, bool)
	WithdrawReward(validatorAddress string, fromAddress string) (error, bool)
	GetRewardBalance(validatorAddress string, fromAddress string) (error, []Coin)
	calculateWeight()
}
