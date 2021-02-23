package tools

import (
	"log"

)

func To4byte(content... byte) [4]byte {

	if len(content) < 4 {
		log.Fatalf("you need 4 element of 1 byte to make a 4 byte structure")
	}

	return [4]byte{byte(content[0]), byte(content[1]), byte(content[2]), byte(content[3])}
}