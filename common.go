package main

import (
	"os"
)

// Common, cross platform information. Calls platform specific functions built
// for each platform in other files.

func sysInfo() (map[string]interface{}, error) {
	var err error

	info := make(map[string]interface{})
	info["hostname"], err = os.Hostname()
	if err != nil {
		return nil, err
	}
	info["cpu"], err = cpuInfo()
	if err != nil {
		return nil, err
	}

	info["kernel"], err = kernelInfo()
	if err != nil {
		return nil, err
	}

	return info, nil
}
