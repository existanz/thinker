package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"hash"
)

type Hasher struct {
	hash.Hash
}

func MD5(data []byte) []byte {
	hasher := md5.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func H256(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func IsHashEqual(a, b []byte) bool {
	return string(a) == string(b)
}
