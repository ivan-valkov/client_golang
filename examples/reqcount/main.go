// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A minimal example of how to include Prometheus instrumentation.
package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	gzipRequests := 0
	nonGZIPRequests := 0
	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", http.HandlerFunc(func(rsp http.ResponseWriter, req *http.Request) {
		if promhttp.GzipAccepted(req.Header) {
			gzipRequests++
		} else {
			nonGZIPRequests++
		}
		h := promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				// Opt into OpenMetrics to support exemplars.
				EnableOpenMetrics: true,
			},
		)
		h.ServeHTTP(rsp, req)
	}))
	http.Handle("/req-count", http.HandlerFunc(func(rsp http.ResponseWriter, req *http.Request) {
		respBody := fmt.Sprintf("GZIP requests:%d\n" +
			"Non-GZIP requests:%d\n" +
			"Total requests:%d\n", gzipRequests, nonGZIPRequests, gzipRequests+nonGZIPRequests)
		rsp.Write([]byte(respBody))
	}))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
