package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World\n"))
	})

	http.ListenAndServe(":8080", nil)
}

/*
% wrk -c 100 -r 100000 -t 4 -H 'Authorization: basic apikey_value' http://mbp.local:8080/
Making 100000 requests to http://mbp.local:8080/
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.55ms    7.50ms  86.50ms   98.56%
    Req/Sec     4.07k   259.33     5.00k    92.81%
  100000 requests in 5.06s, 14.02MB read
Requests/sec:  19763.79
Transfer/sec:      2.77MB


MAXPROCS=20

% wrk -c 100 -r 100000 -t 4 -H 'Authorization: basic apikey_value' http://mbp.local:8080/
Making 100000 requests to http://mbp.local:8080/
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    11.71ms   19.87ms 112.16ms   88.24%
    Req/Sec     9.06k     1.24k   13.00k    83.82%
  100034 requests in 2.62s, 14.02MB read
Requests/sec:  38241.17
Transfer/sec:      5.36MB

*/
