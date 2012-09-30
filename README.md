# Go and Java

This code is a downstream derivation of what was started at this article from
Collin at Boundary in this [blog post](http://boundary.com/blog/2012/09/17/comparing-go-and-java-part-2-performance/).

## How to

 * `make`       # sets up general items, the database, builds the apps

 * `make go`    # starts the go server, run this in it's own shell
 * `make wrk`   # runs the http benchmarking tool wrk

 * `make java`  # starts the java server, run this in it's own shell
 * `make wrk`   # runs the http benchmarking tool wrk

There are other makefile targets, read the Makefile.

## Observations

 * Go 1.0.3 was out of the box faster than in the boundary article on my OSX.
 * The patch mentioned below fixes the most major bug on the Go side
 * The JIT costs a lot in the JVM to start with, then settles down
 * The OSX kernel starts to freak under this load, 50+% system load at 10kr/s+
 * wieghttp and wrk are the only valid http load tools for this performance
 * Go's build tools rock
 * Both systems are breaching the 10kr/s boundary on a laptop. that's excellent
 * 10mb of ram vs 150mb of ram is a significant difference
 * net/http has far less features than the dropwizard stack
 * database/sql has far less features than jdbc
 * It may be worth / fair trying Gmx and Gorilla against Dropwizard
 * It might be interesting to add some HTML rendering tests
 * I had trouble with the Java apps accept backlog when going for 200+ conns
 * shmmax and shmall are set pretty high for the concurrent tests, need db host
 * Both systems now saturate CPU

## Additional Items

### Patch database/sql bug in go

 * database/sql in 1.0.3 has a hardcoded idle connection count of 2. it's not
   really an idle connection count though, it's more of a "retained connection
   count", bad for postgres and benchmarks / concurrency in general.
 * See patch/go-sql.diff - this patch isn't perfect, there are real world
   subtleties due for discussion on the Go ML.
 * TODO - get that patch submitted

### Postgres tuning

 * postgres max_connections set high, and shmmax and shmall set high

### Network Tuning / Notes

 * ulimit -n 4096 (requires tuning prior tuning on OSX)
 * my network runs "jumbo frames" (mtu 9000)
 * somaxconn at default of 128, go uses this for listen backlog
 * Mountain Lion has higher values for recvspace etc than older versions
 * Moutnain Lion also runs a shorter msl timeout
 * My system runs with very high maxfiles, this might be old tuning, but I can't
   find it - sysctl.conf is missing. It's 65535 anyway. ulimit -n works.
 * FTR this system was Lion Server, then upgraded to ML non-server. *shrug*

 * TODO - most of these are TODO Linux essentially
 * TODO - setrlimit nofile
 * TODO - limits.conf / nofile
 * TODO - /proc/sys/fs/file-max
 * TODO - sysctl kern.maxfiles
 * TODO - sysctl kern.maxfilesperproc
 * TODO - firewall bypass (if multi-host)
 * TODO - jumbo frames
 * TODO - window scaling
 * TODO - msl tuning if we get to real high concurrency & lossy

## Example Results

### Go (patched database/sql) - 12MB of ram - 2 hosts:

```
@min ~ % wrk -c 200 -r 100000 -t 4 -H 'Authorization: basic apikey_value' http://mbp.local:8080/authenticate
Making 100000 requests to http://mbp.local:8080/authenticate
  4 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    11.44ms   13.86ms  93.02ms   84.38%
    Req/Sec     3.96k   190.58     4.00k    96.25%
  100000 requests in 6.15s, 24.99MB read
Requests/sec:  16270.02
Transfer/sec:      4.07MB
@min ~ % wrk -c 200 -r 100000 -t 4 -H 'Authorization: basic apikey_value' http://mbp.local:8080/authenticate
Making 100000 requests to http://mbp.local:8080/authenticate
  4 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    14.76ms   15.96ms  81.73ms   90.00%
    Req/Sec     3.94k   231.13     4.00k    94.38%
  100000 requests in 6.21s, 24.99MB read
Requests/sec:  16102.53
Transfer/sec:      4.02MB
```

### Java (after a prior run of 1M requests) - 228MB of ram:

```
@min ~ % wrk -c 200 -r 100000 -t 4 -H 'Authorization: basic apikey_value' http://mbp.local:8080/authenticate
Making 100000 requests to http://mbp.local:8080/authenticate
  4 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    19.75ms    6.06ms  61.24ms   90.15%
    Req/Sec     2.00k     0.00     2.00k   100.00%
  100000 requests in 10.08s, 24.99MB read
  Socket errors: connect 0, read 0, write 0, timeout 4
Requests/sec:   9916.83
Transfer/sec:      2.48MB
```
