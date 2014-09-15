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
