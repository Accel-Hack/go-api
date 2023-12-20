# go-api

api application composed by Golang

## How to run server.

### 1. start MySql server

Start MySQL server with docker compose.

```bash
$ docker compose -f ./compose/compose.yml up mysql
```

> [!TIP]
> ./testdata/mysql/initdb.d direcotry is mounted to mysql container /docker-entrypoint-initdb.d which is loaded in stating MySql server.

### 2. run app

Run go-api application.

```bash
Usage of go-api:
A go-api requires '-mysql.addr' or '-mysql.dsn' (which is prioritized over '-mysql.addr').

  -log.level value
    	Logging level one of [DEBUG INFO WARN ERROR]
  -mysql.addr string
    	MySQL URL. Required if mysql.dsn is empty (default "localhost:3566")
  -mysql.database string
    	Database name (default "YOUR_APPLICATION")
  -mysql.dsn string
    	Data source name format defined as follow: "[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]"
  -mysql.password string
    	Password
  -mysql.table string
    	Table name (default "SAMPLE")
  -mysql.user string
    	Username (default "root")
  -server.host string
    	Host to serve (default "localhost")
  -server.port string
    	Port to serve (default "8080")

$ go run ./cmd/go-api/main.go -mysql.password="root@123"
2023/12/20 17:57:02 connect MySql to "root:root@123@tcp(localhost:3566)/YOUR_APPLICATION?allowNativePasswords=false&checkConnLiveness=false&maxAllowedPacket=0"
2023/12/20 17:57:02 expose GET "/sample"
2023/12/20 17:57:02 expose PUT "/sample"
2023/12/20 17:57:02 expose POST "/sample"
2023/12/20 17:57:02 expose DELETE "/sample"
2023/12/20 17:57:02 expose GET "/samples"
2023/12/20 17:57:02 Linten on localhost:8080
```

### 3. request

Request go-api.

```console
// GET "/samples"
$ curl -s "localhost:8080/samples" | jq
{
  "Total": 5,
  "Samples": [
    {
      "ID": "2e40b651-c32e-4dab-85bd-5a2a81f58c58",
      "Name": "kawamura1",
      "Birthday": "1994-09-14T00:00:00+09:00",
      "IsJapanese": true
    },
    {
      "ID": "7d937a5e-7fa3-4676-949c-6366e988d830",
      "Name": "kawamura2",
      "Birthday": "1994-10-12T00:00:00+09:00",
      "IsJapanese": true
    },
    {
      "ID": "ee4d8f69-7b37-45b2-ba55-08a23e429ec3",
      "Name": "kawamura3",
      "Birthday": "1994-11-08T00:00:00+09:00",
      "IsJapanese": true
    },
    {
      "ID": "f32b76d3-6972-4b62-b19c-1d31bfc88e54",
      "Name": "kawamura5",
      "Birthday": "1994-12-12T00:00:00+09:00",
      "IsJapanese": true
    },
    {
      "ID": "ffda86bf-ee4d-443b-9dcd-5ec9881209b3",
      "Name": "kawamura4",
      "Birthday": "1994-11-08T00:00:00+09:00",
      "IsJapanese": true
    }
  ]
}

// GET "/sample"
$ curl -s "localhost:8080/sample?id=ee4d8f69-7b37-45b2-ba55-08a23e429ec3" | jq
{
  "ID": "ee4d8f69-7b37-45b2-ba55-08a23e429ec3",
  "Name": "kawamura3",
  "Birthday": "1994-11-08T00:00:00+09:00",
  "IsJapanese": true
}

$ curl -s "localhost:8080/samples?limit=2" | jq
{
  "Total": 5,
  "Samples": [
    {
      "ID": "2e40b651-c32e-4dab-85bd-5a2a81f58c58",
      "Name": "kawamura1",
      "Birthday": "1994-09-14T00:00:00+09:00",
      "IsJapanese": true
    },
    {
      "ID": "7d937a5e-7fa3-4676-949c-6366e988d830",
      "Name": "kawamura2",
      "Birthday": "1994-10-12T00:00:00+09:00",
      "IsJapanese": true
    }
  ]
}

$ curl -s "localhost:8080/samples?limit=2&offset=2" | jq
{
  "Total": 5,
  "Samples": [
    {
      "ID": "ee4d8f69-7b37-45b2-ba55-08a23e429ec3",
      "Name": "kawamura3",
      "Birthday": "1994-11-08T00:00:00+09:00",
      "IsJapanese": true
    },
    {
      "ID": "f32b76d3-6972-4b62-b19c-1d31bfc88e54",
      "Name": "kawamura5",
      "Birthday": "1994-12-12T00:00:00+09:00",
      "IsJapanese": true
    }
  ]
}

$ curl -s "localhost:8080/samples?limit=2&offset=4" | jq
{
  "Total": 5,
  "Samples": [
    {
      "ID": "ffda86bf-ee4d-443b-9dcd-5ec9881209b3",
      "Name": "kawamura4",
      "Birthday": "1994-11-08T00:00:00+09:00",
      "IsJapanese": true
    }
  ]
}

// PUT "/sample"
$ curl -i "localhost:8080/sample?id=ee4d8f69-7b37-45b2-ba55-08a23e429ec3&name=mugi&birthday=2022-12-25&is_japanese=false" -XPUT
HTTP/1.1 200 OK
Date: Wed, 20 Dec 2023 09:01:41 GMT
Content-Length: 46
Content-Type: text/plain; charset=utf-8

{"id":"53c33c68-d394-4af9-9776-5b96377ba00b"}

$ curl -s "localhost:8080/sample?id=53c33c68-d394-4af9-9776-5b96377ba00b" | jq
{
  "ID": "53c33c68-d394-4af9-9776-5b96377ba00b",
  "Name": "mugi",
  "Birthday": "2022-12-25T09:00:00+09:00",
  "IsJapanese": false
}

// POST "/sample"
$ curl -i "localhost:8080/sample?id=ee4d8f69-7b37-45b2-ba55-08a23e429ec3&name=ayanodesh&birthday=1994-05-20&is_japanese=false" -XPOST
HTTP/1.1 200 OK
Date: Wed, 20 Dec 2023 09:05:14 GMT
Content-Length: 46
Content-Type: text/plain; charset=utf-8

{"id":"ee4d8f69-7b37-45b2-ba55-08a23e429ec3"}

$ curl -s "localhost:8080/sample?id=ee4d8f69-7b37-45b2-ba55-08a23e429ec3" | jq
{
  "ID": "ee4d8f69-7b37-45b2-ba55-08a23e429ec3",
  "Name": "ayanodesh",
  "Birthday": "1994-05-20T09:00:00+09:00",
  "IsJapanese": false
}

// DELETE "/sample"
$ curl -i "localhost:8080/sample?id=2e40b651-c32e-4dab-85bd-5a2a81f58c58" -XDELETE
HTTP/1.1 200 OK
Date: Wed, 20 Dec 2023 09:06:33 GMT
Content-Length: 0

// check
$ curl -s "localhost:8080/samples" | jq
{
  "Total": 5,
  "Samples": [
    {
      "ID": "53c33c68-d394-4af9-9776-5b96377ba00b",
      "Name": "mugi",
      "Birthday": "2022-12-25T09:00:00+09:00",
      "IsJapanese": false
    },
    {
      "ID": "7d937a5e-7fa3-4676-949c-6366e988d830",
      "Name": "kawamura2",
      "Birthday": "1994-10-12T00:00:00+09:00",
      "IsJapanese": true
    },
    {
      "ID": "ee4d8f69-7b37-45b2-ba55-08a23e429ec3",
      "Name": "ayanodesh",
      "Birthday": "1994-05-20T09:00:00+09:00",
      "IsJapanese": false
    },
    {
      "ID": "f32b76d3-6972-4b62-b19c-1d31bfc88e54",
      "Name": "kawamura5",
      "Birthday": "1994-12-12T00:00:00+09:00",
      "IsJapanese": true
    },
    {
      "ID": "ffda86bf-ee4d-443b-9dcd-5ec9881209b3",
      "Name": "kawamura4",
      "Birthday": "1994-11-08T00:00:00+09:00",
      "IsJapanese": true
    }
  ]
}
$ curl -s "localhost:8080/samples?name=aya" | jq
{
  "Total": 1,
  "Samples": [
    {
      "ID": "ee4d8f69-7b37-45b2-ba55-08a23e429ec3",
      "Name": "ayanodesh",
      "Birthday": "1994-05-20T09:00:00+09:00",
      "IsJapanese": false
    }
  ]
}
```

## How to run tests.

Not yet!!!!!!!
