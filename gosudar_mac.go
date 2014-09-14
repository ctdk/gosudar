// +build darwin

// Mac specific information gathering functions

package main

import (
	"strings"
	"syscall"
)

func cpuInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	cpuKeys := map[string]string{ "vendor": "vendor_id", "brand_string": "model_name" }
	for sysctlName, gohaiName := range cpuKeys {
		k, err := syscall.Sysctl("machdep.cpu."+sysctlName)
		if err != nil {
			return nil, err
		}
		info[gohaiName] = k
	}
	cpuKeyInts := map[string]string{ "model": "model", "family": "family", "stepping": "stepping" }
	for sysctlName, gohaiName := range cpuKeyInts {
		k, err := syscall.SysctlUint32("machdep.cpu."+sysctlName)
		if err != nil {
			return nil, err
		}
		info[gohaiName] = k
	}
	hwKeyInts := map[string]string{ "physicalcpu":"real", "logicalcpu":"total", "cpufrequency":"mhz" }
	for sysctlName, gohaiName := range hwKeyInts {
		k, err := syscall.SysctlUint32("hw."+sysctlName)
		if err != nil {
			return nil, err
		}
		info[gohaiName] = k
	}
	info["mhz"] = info["mhz"].(uint32) / 1000000

	cpuFlags, err := syscall.Sysctl("machdep.cpu.features")
	if err != nil {
		return nil, err
	}
	info["flags"] = strings.Split(strings.ToLower(cpuFlags), " ")

	return info, nil
}

func kernelInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})
	// still need modules
	kernKeys := map[string]string{ "ostype": "name", "osrelease": "release", "version":"version" }
	for sysctlName, gohaiName := range kernKeys {
		k, err := syscall.Sysctl("kern."+sysctlName)
		if err != nil {
			return nil, err
		}
		info[gohaiName] = k
	}
	info["os"] = info["name"]
	return info, nil
}
