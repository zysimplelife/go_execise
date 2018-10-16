package stringutil

import (
	"runtime"
)

func GetOsInfo() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "OS X."
	case "linux":
		return "Linux."
	default:
		return "%s."
	}
}
