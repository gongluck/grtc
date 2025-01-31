@echo off

rem 设置项目根目录
set root=%cd%
rmdir /s /q "%root%\bin"

rem 编译 signaling_server
cd "%root%\signaling_server\src"
go build -o "..\..\bin\signaling_server\signaling_server.exe" main.go

rem 编译 website_server
cd "%root%\website_server\src"
go build -o "..\..\bin\website_server\website_server.exe" main.go
xcopy /E /I templates "..\..\bin\website_server\templates"

rem 编译 pusher
cd "%root%\pusher\src"
go clean
rem 设置环境变量
set ANDROID_NDK_BIN=%ANDROID_NDK_HOME%\toolchains\llvm\prebuilt\windows-x86_64\bin
set ANDROID_NDK_SYSROOT=%ANDROID_NDK_HOME%\toolchains\llvm\prebuilt\windows-x86_64\sysroot
set GOOS=android
set GOARCH=arm64
set CGO_ENABLED=1
set CC=%ANDROID_NDK_BIN%\aarch64-linux-android21-clang
set CXX=%ANDROID_NDK_BIN%\aarch64-linux-android21-clang++
set CGO_CFLAGS=""
set CGO_CXXFLAGS=""
set CGO_LDFLAGS="%ANDROID_NDK_SYSROOT%\usr\lib\aarch64-linux-android\21\libc++.a"

del webrtc\cpp\c++_android.o
rem 生成 C++ 对象文件
go generate -tags android .\...

rem 编译 Go 代码
go build -tags android -o "..\..\bin\pusher\pusher_android_arm64" main.go
