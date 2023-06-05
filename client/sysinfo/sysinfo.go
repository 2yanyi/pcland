// Source code file, created by Developer@YANYINGSONG.

package sysinfo

import (
	"github.com/klauspost/cpuid/v2"
	"library/generic/chars"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

const __unknown = "unknown"

func wmicValue(script string) string {
	ls := strings.Split(chars.CmdC("", script), "\n")
	for i := 0; i < len(ls); i++ {
		if i == 1 {
			value := strings.TrimSpace(ls[i])
			if value != "" {
				return value
			}
			break
		}
	}

	return __unknown
}

func unixValue(fp string) string {
	if chars.FileExist(fp) {
		value := strings.TrimSpace(chars.CatString(fp))
		if value != "" {
			return value
		}
	}

	return __unknown
}

func macValueCmd(cmd string) string {
	out, _ := exec.Command("bash", "-c", cmd).Output()
	values := strings.Split(string(out), ":")
	if len(values) >= 2 {
		return strings.TrimSpace(values[1])
	}
	value := strings.ReplaceAll(string(out), "\n", " ")

	return strings.TrimSpace(value)
}

func GetProductName() string {
	switch runtime.GOOS {
	case "darwin":
		return macValueCmd(`system_profiler SPHardwareDataType | grep "Model Name"`)
	case "windows":
		return wmicValue("wmic csproduct get name")
	}

	return unixValue("/sys/class/dmi/id/product_name")
}

func GetLanIPv4() string {
	ls, err := net.InterfaceAddrs()
	if err != nil {
		return __unknown
	}
	address := ""
	for _, addr := range ls {
		switch {
		case strings.HasSuffix(addr.String(), "/24"):
		case strings.HasSuffix(addr.String(), "/20"):
		default:
			continue
		}
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() == nil {
				continue
			}
			if strings.HasSuffix(ipNet.IP.String(), ".1") {
				continue
			}
			if address == "" {
				address = ipNet.IP.String()
				break
			}
		}
	}

	return address
}

// GetTarget return Processor.System
func GetTarget() (string, string) {
	title := systemTitle()
	arch := runtime.GOARCH
	goos := runtime.GOOS
	cpu := cpuid.CPU
	vendor := cpu.VendorString

	switch goos {
	case "windows":
		goos = systemTitleWindows(title)
	case "linux":
		goos = systemTitleLinux(title)
	case "darwin":
		goos = "Apple"
	}

	switch arch {
	case "amd64":
		switch vendor {
		case "AuthenticAMD":
			arch = "A64"
		case "GenuineIntel":
			arch = "I64"
		}
		if !cpu.Supports(cpuid.AVX2) {
			arch = "old-" + arch
		}
	case "arm64":
		if goos == "Apple" {
			arch = macValueCmd(`sysctl -a | grep machdep.cpu.brand_string`)
			arch = strings.TrimPrefix(arch, "Apple ")
			arch = strings.ReplaceAll(arch, " ", "_")
		}
	case "riscv64":
		arch = "RISC-V"
	case "loong64":
		arch = "龙芯"
	}

	target := strings.Join([]string{arch, goos}, ".")

	return target, title
}

func systemTitle() string {
	switch runtime.GOOS {
	case "darwin":
		return macValueCmd("sw_vers --productName; sw_vers --productVersion")
	case "windows":
		return wmicValue("wmic os get caption")
	}

	value := unixValue("/etc/issue.net")
	if value == __unknown {
		return unixValue("/etc/system-releases")
	}

	return value
}

func systemTitleWindows(title string) string {
	goos := "Windows"

	switch {
	case strings.Contains(title, "Windows 11"):
		goos += "11"
	case strings.Contains(title, "Windows 10"):
		goos += "10"
	case strings.Contains(title, "Windows 8"):
		goos += "8"
	case strings.Contains(title, "Windows 7"):
		goos += "7"
	case strings.Contains(title, "Windows Embedded"):
		goos += "E"
	}

	return goos
}

func systemTitleLinux(title string) string {
	goos := "Linux"

	switch {
	case strings.Contains(title, "Arch"):
		goos = "Arch"
	case strings.Contains(title, "Debian"):
		goos = "Debian"
	case strings.Contains(title, "Raspbian"):
		goos = "Raspbian"
	case strings.Contains(title, "Kali"):
		goos = "Kali"
	case strings.Contains(title, "Ubuntu"):
		goos = "Ubuntu"
	case strings.Contains(title, "Pop!"):
		goos = "Pop!"
	case strings.Contains(title, "Red Hat"):
		goos = "RedHat"
	case strings.Contains(title, "Fedora"):
		goos = "Fedora"
	case strings.Contains(title, "CentOS"):
		goos = "CentOS"
	case strings.Contains(title, "openSUSE"):
		goos = "openSUSE"
	case strings.Contains(title, "Manjaro"):
		goos = "Manjaro"
	case strings.Contains(title, "Amazon"):
		goos = "Amazon"
	case strings.Contains(title, "Alpine"):
		goos = "Alpine"
	case strings.Contains(title, "MX Linux"):
		goos = "MX"
	case strings.Contains(title, "Oracle"):
		goos = "Oracle"
	case strings.Contains(title, "Slackware"):
		goos = "Slackware"
	}

	return goos
}
