# Rotate Elasticsearch Index
A tool to delete indices in Elasticsearch older than a given number of days. Index names must contain date in YYYY.MM.DD format.

## Building from Source
Clone repo into your go path under `$GOPATH/src`:

```sh
$ git clone https://github.com/kairen/rotate-elasticsearch-index.git $GOPATH/src/github.com/kairen/rotate-elasticsearch-index
$ cd $GOPATH/src/github.com/kairen/rotate-elasticsearch-index
$ make
```

## Usage

```sh 
$ ./out/rotate-index
Usage of ./out/rotate-index:
      --alsologtostderr                  log to standard error as well as files
      --date-format string               Format template for parsing date. (default "2006.1.2")
  -d, --days int                         Days to keep. (default 90)
      --endpoints strings                Endpoints of elasticsearch. (default [http://elasticsearch:9200])
      --index-regex-patterns strings     Index's regex pattern.
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --password string                  Password for basic auth.
      --retry-count int                  The number of retry for deleting request. (default 5)
      --sniffer                          Enable client to use a sniffing process for finding all nodes of your cluster.
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
      --tls                              Enable that servers are TLS.
      --tls.ca string                    SSL Certificate Authority file used to secure elasticsearch communication.
      --tls.cert string                  SSL certification file used to secure elasticsearch communication.
      --tls.key string                   SSL key file used to secure elasticsearch communication.
      --tls.skip-host-verify             (insecure) Skip server's certificate chain and host name verification
      --username string                  Username for basic auth.
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```