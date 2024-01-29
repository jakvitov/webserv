![Build](https://github.com/jakvitov/webserv/actions/workflows/build.yml/badge.svg)![Tests](https://github.com/jakvitov/webserv/actions/workflows/tests.yml/badge.svg) ![Audit](https://github.com/jakvitov/webserv/actions/workflows/audit.yml/badge.svg)

```text
 __  __  _ _ _ _____ _____ _____ _____ _____ _____           
____    | | | |   __| __  |   __|   __| __  |  |  |          
   ___  | | | |   __| __ -|__   |   __|    -|  |  |          
 ___  _ |_____|_____|_____|_____|_____|__|__|\___/           

```
# Webserv - lightweight webserver
Webserv is a **lightweight** Go based webserver. Besides classic static page serving and logging capablilities, it also offers reverse proxy functionality. Webserv can be used to **aggregate all HTTP/S based services** under one port and log all access in one point.   

## Running Webserv
Download binaries for your target platform and start the server with configuration according to the *Configuration* section. The following screen should appear:
```log
[INFO];[YOUR_DATE_TIME];Creating http server for port [TARGET_PORT]
[INFO];[YOUR_DATE_TIME];Created server cache with max size of [MAX_CACHE_SET] bytes.
[INFO];[YOUR_DATE_TIME];                                                                                                                          
 __  __  _ _ _ _____ _____ _____ _____ _____ _____           
____    | | | |   __| __  |   __|   __| __  |  |  |          
   ___  | | | |   __| __ -|__   |   __|    -|  |  |          
 ___  _ |_____|_____|_____|_____|_____|__|__|\___/           

                                                             
                                                             
[INFO];[YOUR_DATE_TIME];Starting listener on port [:TARGET_PORT]

```

## Configuration
Configuration is done trough a **yaml** config file. Path to this file is passed to webserv as its first parameter. For instance:   
- On Unix based systems:                                                        
```bash
webserv ./my_web_config.yaml
```
- On Windows:
```bash
webserv.exe .\my_web_config.yaml
```


## Config structure
- Structure of the yaml config file is as follows:
```yaml
ports:
  http_port: int
  https_port: int

logger:
  level: [INFO, WARN, ERROR, FATAL]
  output_to_file: bool
  output_file: string
  append_output: bool

handler:
  content_root: string
  read_timeout_ms: int
  write_timeout_ms: int
  max_header_bytes: int
  cache_enabled: bool
  //If cache is enabled and this is 0 or not set, we default to 20MB
  max_cache_bytes: int64

reverse_proxy:
  routes:
    - from: string
      to: int
    - from: string
      to: int

security: 
  cert_path: string
  priv_key_path: string 

```

### Default config values
If any of the config values is not present, following values are used instead. Values not listed here are ignored, when not set meaning, that their feature is turned off. 

```yaml
ports:
  http_port: 8080

logger:
  level: "INFO"
  output_to_file: true
  output_file: "./websersv_log.log"
  append_output: false

handler:
  #Mandatory content root or proxy settings
  content_root: string
  read_timeout_ms: 1000
  write_timeout_ms: 1000
  max_header_bytes: 1 << 20
  cache_enabled: false
  max_cache_bytes: 20 * 1000 * 1000

```

### Security TLS setup
For correct HTTPS setup, a TLS certificate and private key is required. You can self sign one, or get a CA issued pem. Note, that self signing certificates is recommended only for testing purposes. 
Webserv requires you to have **both certificates in the PEM format and with .pem suffix!** Provide path to both private key and certificate in the *security* yaml key. 
- For example:

 ```yaml
 security: 
  cert_path: ./my_certificate.pem
  priv_key_path: ./my_private_key.pem
 ```

# Todo
- Help man page
- Load testing
- Cache rebalance atomicity
- Write documentation
