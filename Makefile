schema: sql/schema.sql
	dropdb authdb || echo no db
	createdb authdb
	psql authdb < sql/schema.sql

data: sql/data.sql
	psql authdb < sql/data.sql

bin/main: src/**/*
	GOPATH="${PWD}" go install main

go: bin/main
	GOMAXPROCS=24 ./bin/main cfg/config.json

java/target/auth-1.0.jar: java/src/main/java/auth/*.java java/pom.xml
	cd java && mvn install

java: java/target/auth-1.0.jar
	java -Xserver -jar java/target/auth-1.0.jar server ./java/auth.yml

all: schema data bin/main java/target/auth-1.0.jar

clean:
	cd java && mvn clean
	GOPATH="${PWD}" go clean -i -x ./...

# ab uses HTTP/1.0 causing an extra header for Go (Connection: close). Java
# violates HTTP by sending a 1.1 response to a 1.0 request!
ab:
	ab -A apikey_value -c 1 -n 1 http://localhost:8080/authenticate

# siege sends an Accept-Encoding: gzip header, which causes Java to gzip the
# response, but Go does not.
siege:
	siege -b -r 1 -c 1 -H 'Authorization: basic apikey_value' \
		http://localhost:8080/authenticate

# httperf causes identical behavior at the HTTP level
httperf:
	httperf --hog --add-header 'Authorization: basic apikey_value\n' \
		--server localhost --port 8080 --uri /authenticate \
		--num-calls 1000 --num-conns 50

# wrk and weighttp are worthy of 2012 benchmarks. They can actually generate
# significant http request loads. wrk is more dynamic, whereas weighttp may be
# more stable. Both provide some valid insights.
#
# N.B. t4 c10 is not a fully balanced workload. this is not an accident, you'll
# want to understand what it causes in the statistics and not just take initial
# results as writ.
#
# N.B. both of these systems utilize keep-alive which reduces connect overhead.
# This is also intentional. It's possible that looking at accept rates can be
# useful, but to do so at anything approaching the performance of these two
# applicaiton platforms, singificant OS tuning is required.
wrk:
	wrk -c 10 -r 100000 -t 4 -H 'Authorization: basic apikey_value' \
		http://localhost:8080/authenticate

weighttp:
	weighttp -c 10 -n 1000000 -t 4 -k -H 'Authorization: basic apikey_value' \
		http://localhost:8080/authenticate
