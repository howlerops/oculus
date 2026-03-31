package platform

import "runtime"

func GetPlatform() string { return runtime.GOOS }
func IsMacOS() bool       { return runtime.GOOS == "darwin" }
func IsLinux() bool       { return runtime.GOOS == "linux" }
func IsWindows() bool     { return runtime.GOOS == "windows" }
func GetArch() string     { return runtime.GOARCH }
