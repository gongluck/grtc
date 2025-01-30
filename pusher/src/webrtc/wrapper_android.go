//go:build android
// +build android

/*
 * @Author: gongluck
 * @Date: 2025-01-29 22:41:22
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-31 01:29:25
 */

package webrtc

//go:generate rm -f ./cpp/c++_android.o
//go:generate ${CC} -c ./cpp/c++.cpp -o ./cpp/c++_android.o

/*
#cgo CFLAGS: -I${SRCDIR}/cpp
#cgo CXXFLAGS: -I${SRCDIR}/cpp
#cgo LDFLAGS: -L${SRCDIR}/cpp ${SRCDIR}/cpp/c++_android.o
#include "c++.h"
*/
import "C"
import "log"

func CallCppFunction() {
	log.Println("Calling android CppFunction()...")
	C.CppFunction()
}
