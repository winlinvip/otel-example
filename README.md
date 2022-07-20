# otel-example

Example for https://opentelemetry.io/docs/instrumentation/go/getting-started/

## Usage: otel-example

Example for [otel](https://opentelemetry.io/docs/instrumentation/go/getting-started/#creating-a-console-exporter).

Build and run:

```bash
(cd otel-example && go build . && ./otel-example)
```

## Usage: otel-http

Example for [http](https://opentelemetry.io/docs/instrumentation/go/libraries/)

Build and run:

```bash
(cd otel-http && go build . && ./otel-http)
```

Test by curl:

```bash
curl http://127.0.0.1:8094
```

## Usage: otel-mysql

Example for [mysql](https://opentelemetry.io/registry/?language=go&component=instrumentation)

Run mysql in docker:

```bash
docker run --rm -e MYSQL_ROOT_HOST=% -e MYSQL_ROOT_PASSWORD=12345678 -p 13306:3306 mysql/mysql-server:latest
```

Please test by:

```bash
mysql -uroot --host=127.0.0.1 --port=13306 -p12345678 -e 'select CURRENT_TIMESTAMP'
```

Build and run:

```bash
(cd otel-mysql && go build . && ./otel-mysql)
```

Test by curl:

```bash
curl http://127.0.0.1:8095
```

## Usage: otel-baggage

Example for [baggage](https://opentelemetry.io/docs/concepts/signals/baggage/)

Build and run:

```bash
(cd otel-baggage-server && go build . && ./otel-baggage-server)
(cd otel-baggage && go build . && ./otel-baggage)
```

Or run simple server for baggage:

```bash
(cd otel-baggage-server2 && go build . && ./otel-baggage-server2)
(cd otel-baggage && go build . && ./otel-baggage)
```

