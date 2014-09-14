package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
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

	o, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var out bytes.Buffer
	json.Indent(&out, o, "", "\t")
	out.WriteTo(os.Stdout)
}
