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
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"github.com/ctdk/gosudar/plugin"
	"os/signal"
	"syscall"
)

func main() {
	// System information is, of course, system specific. That's built in
	// a separate file depending on what platform this is running on.

	info, err := sysInfo()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Plugins would run here and be added to the rest of the output, once 
	// that's figured out.

	pDir, err := os.Open(plugin.PluginDir)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	pRun, err := pDir.Readdirnames(0)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	
	stopch := make(chan struct{}, 1)
	donech := make(chan struct{}, 1)
	readych := make(chan struct{}, 1)
	go startPluginServer(stopch, readych, donech)

	go func() {
		for {
			pi := <- plugin.InfoCh
			// TODO: make a merge function. Also, a mutex for the
			// info hash.
			for k, v := range pi {
				info[k] = v
			}
		}
	}()
	<-readych

	for _, v := range pRun {
		cmdStr := plugin.PluginDir + "/" + v
		cmd := exec.Command(cmdStr)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			log.Println("hmm")
			log.Println(stderr.String())
			log.Fatal(err)
			os.Exit(1)
		}
	}

	stopch <- struct{}{}
	<-donech

	o, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var out bytes.Buffer
	json.Indent(&out, o, "", "\t")
	out.WriteTo(os.Stdout)
}

func startPluginServer(stopch <-chan struct{}, readych, donech chan<- struct{}) {
	uaddr, _ := net.ResolveUnixAddr("unix", "/tmp/gosudar.sock")
	l, err := net.ListenUnix("unix", uaddr)
	readych <- struct{}{}
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal){
		<-c
		l.Close()
		os.Exit(0)
	}(sigch)

	if err != nil {
		log.Printf("Failed to start socket for plugins: %s\n", err.Error())
		os.Exit(1)
	}
	rpc.Register(new(plugin.Info))
	done := false
	go func(){
		for {
			log.Printf("Waiting for plugins...")
			if conn, err := l.AcceptUnix(); err == nil {
				log.Println("reading data from plugin...")
				go jsonrpc.ServeConn(conn)
			} else {
				if !done {
					log.Printf("Plugin connection failed: %s ", err.Error())
					os.Exit(1)
				} else {
					return
				}
			}
			
		}
	}()
	<-stopch
	done = true
	err = l.Close()
	donech <- struct{}{}
	if err != nil {
		log.Println("err closing ", err)
	}
	return
}
