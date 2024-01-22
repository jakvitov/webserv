# Webserv - lightweight webserver
Webserv is a **lightweight** Go based webserver. Besides classic static page serving and logging capablilities, it also offers reverse proxy functionality. Webserv can be used to **aggregate all HTTP/S based services** under one port and log all access in one point.   

## Configuration
Configuration is done trough a **yaml** config file. Path to this file is passed to webserv as its first parameter. For instance:                                                          
```bash
webserv ./my_web_config.yaml
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
  #Only mandatory
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
    - to: int
    - from: string
    - to: int

security: 
  cert_path: string
  spam_filter: bool

```

The only mandatory field in the config file is content_root, which does specify the target directory with webpage. 

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
  #Mandatory
  content_root: string
  read_timeout_ms: 1000
  write_timeout_ms: 1000
  max_header_bytes: 1 << 20
  cache_enabled: false
  max_cache_bytes: 20 * 1000 * 1000

```

# Todo
- Help man page
- Reverse proxy routing
- Spam filter
- Load testing
- Cache rebalance atomicity
- Write documentation
