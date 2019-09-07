# [resttester]
Simple CLI tool for load testing HTTP endpoints.

# Usage
`resttester -t {duration} -s {strategy} [options...] [URL]`

Where `options` could be:
- `-X`, `--request` change HTTP method (default: `GET`)
- `-H`, `--header` adds header to the request, can be used multiple times (format: `<header>=<value>`)
- `-P`, `--payload` adds payload to the request, valid only for methods that supports it
- `-L`, `--location` if `3xx` responses should be followed (default: `false`)
- `--max-redirs` maximum amount of follows for `3xx` requests
- `-k`, `--insecure` if server certificate should be verified
- `-c`, `--ca-certificate` path to cert file that contains allowed server certificates, should be in PEM format
- `--dry-run` whether only a single batch of requests should be executed, useful for verifying client setup
- `-t`, `--timeout` duration of the performed test (e.g. `2s`, `10ms`, `1m`)
- `-p`, `--plot` path to image where plot for response time over number of parallel requests would be created
- `-s`, `--strategy` description of a function that will be used to increment the amount of parallel requests.
Currently only linear (a*x+b) and exponential (a*b^x) functions are supported and should be passed as follows:
`lin[a,b]` or `exp[a,b]`

