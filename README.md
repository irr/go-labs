go-labs
-----------

**go-labs**  is a set of sample codes whose main purpose is to experiment and test [Go] programming language

Libraries
-----------

* [cryptobox]: Go port of Cryptobox
* [gocql]: A database/sql driver for CQL, the Cassandra query language for Go
* [go-iconv]: iconv binding for Go
* [go-simplejson]: a Go package to interact with arbitrary JSON
* [go-options]: a Go package to structure and resolve options
* [Go-MySQL-Driver]: A MySQL-Driver for Go
* [redigo]: Go client for Redis
* [redis]: Go client for Redis (Cluster/Sentinel)


```shell
go get -v github.com/garyburd/redigo/redis
go get -v github.com/go-redis/redis
# code in directory /opt/golib/src/github.com/go-redis/redis expects import "gopkg.in/redis.v3"
go get -v github.com/gocql/gocql
go get -v github.com/go-sql-driver/mysql
go get -v github.com/sloonz/go-iconv
go get -v github.com/cryptobox/gocryptobox/box
go get -v github.com/bitly/go-simplejson
go get -v github.com/mreiferson/go-options
```

Copyright and License
---------------------
Copyright 2013 Ivan Ribeiro Rocha

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[Go]: http://golang.org/
[Go-MySQL-Driver]: https://github.com/go-sql-driver/mysql
[gocql]: https://github.com/gocql/gocql
[redigo]: https://github.com/garyburd/redigo
[redis]: https://github.com/go-redis/redis
[go-iconv]: https://github.com/sloonz/go-iconv
[cryptobox]: https://github.com/cryptobox/gocryptobox
[go-simplejson]: https://github.com/bitly/go-simplejson
[go-options]: https://github.com/mreiferson/go-options
