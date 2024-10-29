# Bigip

## Description

`bigip` is a tool designed to help you convert an F5 configuration to
HAProxy. It's important to note that it won't automatically convert
the entire configuration for you. Configuration conversion is an
iterative process.


## Build

To build `bigip` you need to have golang 1.20 (or a newer version)
installed. Run:

```
go build .
```


## Tutorial

### `test_bigip`

To convert `test_bigip` to an HAProxy configuration, we need to create
templates that will be used by the `convert` action.

```
bigip gen-templates -v -o /tmp/test_bigip/templates examples/test_bigip.conf
```

The `-v` flag can be added several times to increase the verbosity level.

Some F5 objects have now been converted into templates in the
`/tmp/test_bigip/templates` directory. You can find them in the
`examples/test_bigip` directory within this repository.


The available templates include:

- `monitor.tpl.cfg`
- `node.tpl.cfg`
- `persistence.tpl.cfg`
- `policy.tpl.cfg`
- `pool.tpl.cfg`
- `profile.tpl.cfg`
- `rule.tpl.cfg`
- `virtual.tpl.cfg`


All templates, except for those in `virtual` and `pool`, have been
intentionally disabled. You will need to review and enable each one.


Here’s an example of a `virtual`:

```
{{ define "virtual:/Common/app1_t443_vs" }}
### virtual: /Common/app1_t443_vs
## file: examples/test_bigip.conf, 18 lines: 316-333
#
#F5# ---8<---8<---8<---8<---8<---
#F5#    ltm virtual /Common/app1_t443_vs {
#F5#        destination /Common/192.168.1.21:443
#F5#        ip-protocol tcp
#F5#        last-modified-time 2020-09-18:10:05:54
#F5#        mask 255.255.255.255
#F5#        pool /Common/app1_t80_pool
#F5#        profiles {
#F5#            /Common/http { }
#F5#            /Common/tcp { }
#F5#        }
#F5#        serverssl-use-sni disabled
#F5#        source 0.0.0.0/0
#F5#        source-address-translation {
#F5#            type automap
#F5#        }
#F5#        translate-address enabled
#F5#        translate-port enabled
#F5#    }
#F5# ---8<---8<---8<---8<---8<---
frontend Common::app1_t443_vs
    bind 192.168.1.21:443
    # profiles
{{   templateIndent 4 "profile:/Common/http" "" }}
{{   templateIndent 4 "profile:/Common/tcp" "" }}
    default_backend Common::app1_t80_pool

#
### /virtual: /Common/app1_t443_vs
{{ end }}
```


Each template includes the original F5 configuration snippet as a
reference for future use. In this example, we also need to convert the
`Common::app1_t80_pool` backend, which can be located in the
`pool.tpl.cfg` file.

```
{{ define "pool:/Common/app1_t80_pool" }}
### pool: /Common/app1_t80_pool
## file: examples/test_bigip.conf, 11 lines: 232-242
#
#F5# ---8<---8<---8<---8<---8<---
#F5#    ltm pool /Common/app1_t80_pool {
#F5#        members {
#F5#            /Common/app1_Node1:80 {
#F5#                address 192.168.1.22
#F5#            }
#F5#            /Common/app1_Node2:80 {
#F5#                address 192.168.1.23
#F5#            }
#F5#        }
#F5#        monitor /Common/http and /Common/tcp
#F5#    }
#F5# ---8<---8<---8<---8<---8<---
backend Common::app1_t80_pool
    # TODO: change to "tcp" if this backend is not for HTTP
    mode http
{{   templateIndent 4 "monitor:/Common/http" "" }}
{{   templateIndent 4 "monitor:/Common/tcp" "" }}
{{   templateIndent 4 "node:/Common/app1_Node1" "" }}
    server Common::app1_Node1 192.168.1.22:80 check
{{   templateIndent 4 "node:/Common/app1_Node2" "" }}
    server Common::app1_Node2 192.168.1.23:80 check
#
### /pool: /Common/app1_t80_pool
{{ end }}
```

Now we can proceed with converting the `/Common/app1_t443_vs` virtual
server into an HAProxy configuration:

```
bigip convert  -v -t /tmp/test_bigip/templates examples/test_bigip.conf  -V '/Common/app1_t443_vs'
```

The `-V` option limits the conversion to the `/Common/app1_t443_vs`
virtual server and all its related objects. If additional pools are
needed (e.g., due to a complex iRule), the `-P` option can be used to
include those pools as well.

In the logs, we can see that both `node:/Common/app1_Node1` and
`node:/Common/app1_Node2` templates are not found. Remember that all
templates, except those for virtuals and pools, are disabled by
default.

We now need to enable them in the `node.tpl.cfg` file:


```
{{ define "-node:/Common/app1_Node1" }}
### node: /Common/app1_Node1
## file: examples/test_bigip.conf, 3 lines: 206-208
#
#F5# ---8<---8<---8<---8<---8<---
#F5#    ltm node /Common/app1_Node1 {
#F5#        address 192.168.1.22
#F5#    }
#F5# ---8<---8<---8<---8<---8<---
#
### /node: /Common/app1_Node1
{{ end }}
```

The `node:/Common/app1_Node1` entry has been created as
`-node:/Common/app1_Node1` (with a leading `-` to disable it). All we
need to do is remove the `-` from both templates and regenerate the
configuration:

```
{{ define "node:/Common/app1_Node1" }}
### node: /Common/app1_Node1
## file: examples/test_bigip.conf, 3 lines: 206-208
#
#F5# ---8<---8<---8<---8<---8<---
#F5#    ltm node /Common/app1_Node1 {
#F5#        address 192.168.1.22
#F5#    }
#F5# ---8<---8<---8<---8<---8<---
#
### /node: /Common/app1_Node1
{{ end }}
```

We can see that `Common::app1_t443_vs` operates over HTTPS. However,
this configuration does not reference any TLS options. In the
`virtual:/Common/app1_t443_vs`, we can add a certificate:

```
frontend Common::app1_t443_vs
    bind 192.168.1.21:443 ssl crt app1.pem
```

You'll need to provide the appropriate SSL certificate in `app1.pem`.


The `app1` configuration also includes an HTTP to HTTPS upgrade in
`virtual:/Common/app1_t80_vs`. We can either leave it as is and
convert both virtuals at once:

```
bigip convert -t /tmp/test_bigip/templates examples/test_bigip.conf  -V '/Common/app1_t80_vs,/Common/app1_t443_vs'
```

Or we can optimize the HAProxy configuration by merging both
frontends:

```
frontend Common::app1_t443_vs
    bind 192.168.1.21:443 ssl crt app1.pem
    bind 192.168.1.21:80
    # profiles
{{   templateIndent 4 "profile:/Common/http" "" }}
{{   templateIndent 4 "profile:/Common/tcp" "" }}
{{   templateIndent 4 "rule:/Common/_sys_https_redirect" "" }}
    default_backend Common::app1_t80_pool
```

Now, the frontend listens for both HTTP and HTTPS traffic and performs
the HTTPS upgrade using a redirection.

After removing all comments, we end up with the following HAProxy
configuration:


```
frontend Common::app1_t443_vs
    bind 192.168.1.21:443 ssl crt app1.pem
    bind 192.168.1.21:80
    mode http
    http-request redirect scheme https code 302
    default_backend Common::app1_t80_pool

backend Common::app1_t80_pool
    mode http
    option httpchk GET / HTTP/1.1
    server Common::app1_Node1 192.168.1.22:80 check
    server Common::app1_Node2 192.168.1.23:80 check
```

Used command:

```
bigip convert -t /tmp/test_bigip/templates examples/test_bigip.conf  -V '/Common/app1_t80_vs,/Common/app1_t443_vs' | grep -v '[[:space:]]*#' | grep -v '^[[:space:]]*$' | grep .
```

## Tips


- Convert services in a logical way.
- Split the generated configuration into multiple files.
- Use `-V` to select only specific virtuals for conversion. The pools
  from the `pool` directive will be automatically included in the
  conversion. Use `-P` to add additional pools for conversion.
- Use the following default file as a starting point:

```
defaults
    mode http
    timeout client 60s
    timeout server 60s
    timeout connect 60s
```

- Use HAProxy check option to validate the generated configuration:

```
haproxy -c -- default.cfg generated-config.cfg
```

- check F5 ltm documentation: https://clouddocs.f5.com/cli/tmsh-reference/v16/modules/ltm/

- Check HAProxy documentation: https://docs.haproxy.org/


## License

Copyright 2021-2024 (c) Sébastien Gross.

Released under GNU Affero General Public License. See the LICENSE file.
