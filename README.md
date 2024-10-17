# dns2ip-templater: A simple tool for resolving DNS records and templating output

**dns2ip-templater** is a command-line tool designed to resolve DNS records and use the
resulting IP addresses to render a user-defined template. This allows you to easily
and dynamically configure configurations, scripts, or other files based on the current
DNS records.

It gives the possibility to run the tool once for cron tasks or as a daemon with an
interval between each rendering.

In my personal use-case, I use `dns2ip-templater` to dynamically update the configuration
of a Traefik middleware ([IPAllowList](https://doc.traefik.io/traefik/master/middlewares/http/ipallowlist/))
because some services are reachable via a dynamic IP.


## Installation

You can install by using the binaries available on the [releases page](https://github.com/Themimitoof/dns2ip-templater/releases).

You can also install it directly on your machine by typing:

```bash
go install github.com/themimitoof/dns2ip-templater
```

## Usage

`dns2ip-templater` requires a configuration file and a template. By default, the config
file path is `config.yml`, the template path is `template.tmpl`. Use the flag `--help`
to check the default value and how to override it:

```
$ dns2ip-templater --help
Usage of dns2ip-templater:
  -conf string
        Configuration file path (default "config.yml")
  -interval duration
        Define the interval to run the routine (by default: executed once)
  -output string
        Rendered file path (default "output.txt")
  -template string
        Template file path (default "template.tmpl")
```

In the configuration file `config.yml`, you need to specify at least a list of
_services_ like `my-house.net`. You can also specify IPs and subnets by using
`ranges`. Here's an example:

```yaml
services:
  - my-house.duckdns.org
  - home.themimitoof.fr

ranges:
  - 8.8.8.8
  - 192.168.1.0/24
```

In a `template.tmpl` file, add your content and make a range like this one in
the following example:

```go
http:
  middlewares:
    no-server-header:
      headers:
        customResponseHeaders:
          Server: ""
    internal_only:
      IPAllowList:
        rejectStatusCode: 200
        sourceRange:
          {{- range $service, $ips := . }}
            # {{ $service }}
            {{- range $ip := $ips}}
            - {{ $ip }}
            {{- end}}
          {{- end }}
```

To run `dns2ip-templater`, simply execute the command `dns2ip-templater` and
take a look on the `output.txt` file.

```toml
$ dns2ip-templater
$ cat output.txt
http:
  middlewares:
    internal_only:
      IPAllowList:
        rejectStatusCode: 200
        sourceRange:
            # Ranges
            - 8.8.8.8
            - 1.1.1.1
            # google.com
            - 2a00:1450:4007:819::200e
            - 142.250.178.142
            # themimitoof.fr
            - 2606:4700:3031::ac43:8201
            - 2606:4700:3032::6815:30a
            - 104.21.3.10
            - 172.67.130.1
```

You can override the path and the name of the output file by using the `-output` flag.

In case you want to run `dns2ip-templater` as a daemon, use the `-interval` flag with
a [Go duration](https://pkg.go.dev/time#ParseDuration) in parameter:

```bash
$ dns2ip-templater -interval 1h -output /traefik/config/dyn-conf/internal-only-middleware.conf
Executing new iteration of dns2ip-templater... Done.
```

## Licenses

This project is released under the MIT license. Feel free to use, contribute, fork and
do what you want with it. Please keep all licenses, copyright notices and mentions in
case you use, re-use, steal, fork code from this repository.
</rewrite_this>

</document>
