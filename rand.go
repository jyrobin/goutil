// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import (
	"math/rand"
	"strings"
	"time"
)

// Ref: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986

var src = rand.NewSource(time.Now().UnixNano())

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

