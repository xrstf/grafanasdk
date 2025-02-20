// This is a simple example of usage of Grafana sdk
// for importing dashboards from a bunch of JSON files (current dir used).
// This example uses SetRawDashboard() and doesn't unmarshal input JSONs.
//
// You are can export dashboards with backup-dashboards utitity.
// NOTE: old dashboards with same names will be silently overrided!
//
// Usage:
//   import-dashboards http://grafana.host:3000 api-key-string-here
//
// You need get API key with Admin rights from your Grafana!
package main

/*
   Copyright 2016 Alexander I.Grafov <grafov@gmail.com>
   Copyright 2016-2019 The Grafana SDK authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

	   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	sdk "github.com/kubermatic/grafanasdk"
)

func main() {
	var (
		filesInDir []os.FileInfo
		rawBoard   []byte
		err        error
	)
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "Usage: import-dashboards http://grafana.host:3000 api-key-string-here\n")
		os.Exit(0)
	}
	ctx := context.Background()
	c, err := sdk.NewClient(os.Args[1], os.Args[2], sdk.DefaultHTTPClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}
	filesInDir, err = ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawBoard, err = ioutil.ReadFile(file.Name()); err != nil {
				log.Println(err)
				continue
			}
			_, err := c.SetRawDashboard(ctx, rawBoard)
			if err != nil {
				log.Printf("error on importing dashboard from %s", file.Name())
				continue
			}
		}
	}
}
