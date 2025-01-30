//go:build !android
// +build !android

/*
 * @Author: gongluck
 * @Date: 2025-01-29 22:41:22
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-30 01:10:55
 */

package webrtc

//go:generate g++ -std=c++11 -c ./cpp/c++.cpp -o ./cpp/c++.o

/*
#cgo CXXFLAGS: -std=c++11
#cgo CFLAGS: -I${SRCDIR}/cpp
#cgo LDFLAGS: -L${SRCDIR}/cpp ${SRCDIR}/cpp/c++.o -lstdc++
#include "c++.h"
*/
import "C"
import "log"

func CallCppFunction() {
	log.Println("Calling CppFunction()...")
	C.CppFunction()
}
