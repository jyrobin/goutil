// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"strings"
)

// Ref: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986

//var src = rand.NewSource(time.Now().UnixNano())
var src = NewCryptoSeededSource()
var rnd = rand.New(src)
var csrc = NewCryptoSource()
var crnd = rand.New(csrc)

func NewCryptoSeededSource() rand.Source {
	var seed int64
	binary.Read(crand.Reader, binary.BigEndian, &seed)
	return rand.NewSource(seed)
}

type cryptoSrc int

func (s cryptoSrc) Seed(seed int64) {}
func (s cryptoSrc) Uint64() uint64 {
	var val uint64
	binary.Read(crand.Reader, binary.BigEndian, &val)
	return val
}
func (s cryptoSrc) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}
func NewCryptoSource() cryptoSrc {
	return 0
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789jy"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		idx := uint8(cache & letterIdxMask)
		sb.WriteByte(letterBytes[idx])
		i--
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}

func RandSouce() rand.Source {
	return src
}
func Rand63() int64 {
	return rnd.Int63()
}
func Rand31() int32 {
	return rnd.Int31()
}
func RandInt() int {
	return rnd.Int()
}

func CrandSouce() rand.Source {
	return csrc
}
func Crand63() int64 {
	return crnd.Int63()
}
func Crand31() int32 {
	return crnd.Int31()
}
func CrandInt() int {
	return crnd.Int()
}

func CrandBytes(sz int) ([]byte, error) {
	b := make([]byte, sz)
	_, err := crand.Read(b)
	return b, err
}
