package core

// Validators for the proof of stake consensus

type IValidator interface {


}

type Validator struct {
	Address string
	Staked []Coin
	Weight float64
	StakedWith string

}

type Validators []Validator


