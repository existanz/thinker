package hash

import (
	"crypto/md5"
	"io"
)

func MD5(data []byte) []byte {
	hasher := md5.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func GetHash(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return MD5(data), nil
}

func IsHashEqual(a, b []byte) bool {
	return string(a) == string(b)
}
