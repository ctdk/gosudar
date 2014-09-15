// +build darwin

/*
 * Copyright (c) 2013-2014, Jeremy Bingham (<jbingham@gmail.com>)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
