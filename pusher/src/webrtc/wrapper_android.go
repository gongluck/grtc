//go:build android
// +build android

/*
 * @Author: gongluck
 * @Date: 2025-01-29 22:41:22
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-30 01:24:45
 */

package webrtc

//go:generate $ANDROID_NDK_HOME/toolchains/llvm/prebuilt/darwin-x86_64/bin/clang++ -target aarch64-linux-android21 -std=c++11 -c ./cpp/c++.cpp -o ./cpp/c++.o

/*
#cgo CXXFLAGS: -std=c++11
#cgo CFLAGS: -I${SRCDIR}/cpp
#cgo LDFLAGS: ${SRCDIR}/cpp/c++.o -lstdc++ -llog -landroid
#include "c++.h"
*/
import "C"
import "log"

func CallCppFunction() {
	log.Println("Calling CppFunction()...")
	C.CppFunction()
}
