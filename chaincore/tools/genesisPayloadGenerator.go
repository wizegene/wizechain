package tools

import (
	"io/ioutil"
)

const crazyGenesisSeed = `46454252554152592032332C2032303231333A303420504D55504441544544203131204d494E555445532041474FA4C61737420737461747565206f66206469637461746f72204672616E636F2072656d6F76656420696e2027686973746F726963206461792720666f7220537061696E`

/*func GeneratePayload() {

	alpha := core.WizegeneAlphabet
	encoded := core.Encode([]byte(crazyGenesisSeed), alpha)
	_ = ioutil.WriteFile("./genesispayload", []byte(encoded), 0644)


}*/

func GetGenesisPayload() []byte {
	data, _ := ioutil.ReadFile("./genesispayload")
	return data
}
