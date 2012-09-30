var http = require('http');

var server = http.createServer(function (request, response) {
  response.writeHead(200, {});
  response.end("Hello World\n");
});

server.listen(8080);

/*

% wrk -c 100 -r 100000 -t 4 -H 'Authorization: basic apikey_value' http://mbp.local:8080/
Making 100000 requests to http://mbp.local:8080/
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.19ms  608.45us   9.41ms   87.05%
    Req/Sec     3.97k   167.78     4.00k    97.12%
  100000 requests in 5.14s, 12.40MB read
Requests/sec:  19462.14
Transfer/sec:      2.41MB

*/
