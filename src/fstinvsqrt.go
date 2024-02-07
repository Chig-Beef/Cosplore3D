// THIS CODE WAS TAKEN FROM https://github.com/itchyny/fastinvsqrt/blob/main/src/go/fastinvsqrt.go
// The repo is for the fast inverse square algorithm in a variety of languages, included Golang
// I rewrited the code instead of depending on it because I don't know how to depend on it with all the other language examples
// It also doesn't seem to be written to be depended on

package main

import "unsafe"

func fastInvSqrt(x float32) float32 {
	i := *(*int32)(unsafe.Pointer(&x))
	i = 0x5f3759df - i>>1
	y := *(*float32)(unsafe.Pointer(&i))
	return y * (1.5 - 0.5*x*y*y)
}
