package libs

import "golang.org/x/crypto/sha3"

func GetSha3512(from []byte) []byte {

	hash := sha3.New512()
	hash.Write(from)
	return hash.Sum(nil)

}

// GetKeccak512 only for legacy keccak512 compatibility
func GetKeccak512(from []byte) []byte {
	hash := sha3.NewLegacyKeccak512()
	hash.Write(from)
	return hash.Sum(nil)
}

// GetKeccak256 only for legacy keccak512 compatibility
func GetKeccak256(from []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(from)
	return hash.Sum(nil)
}

// Faster but recommend SHA3-512 when possible for higher security
func GetSha3256(from []byte) []byte {
	hash := sha3.New256()
	hash.Write(from)
	return hash.Sum(nil)
}
