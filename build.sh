#!/bin/bash

# 设置项目根目录
root=$(pwd)
rm -r $root/bin

# 编译 signaling_server
cd $root/signaling_server/src
go build -o ../../bin/signaling_server/signaling_server main.go

# 编译 website_server
cd $root/website_server/src
go build -o ../../bin/website_server/website_server main.go
cp -r templates ../../bin/website_server/

# 编译 pusher
cd $root/pusher/src
go clean
# 设置环境变量
ANDROID_NDK_BIN=$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/darwin-x86_64/bin
ANDROID_NDK_SYSROOT=$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/darwin-x86_64/sysroot
export GOOS=android
export GOARCH=arm64
export CGO_ENABLED=1
export CC=$ANDROID_NDK_BIN/aarch64-linux-android21-clang
export CXX=$ANDROID_NDK_BIN/aarch64-linux-android21-clang++
export CGO_CFLAGS=""
export CGO_CXXFLAGS=""
export CGO_LDFLAGS="$ANDROID_NDK_SYSROOT/usr/lib/aarch64-linux-android/21/libc++.a"
# 生成 C++ 对象文件
go generate -tags android ./...
# 编译 Go 代码
go build -tags android -o ../../bin/pusher/pusher_android_arm64 main.go
