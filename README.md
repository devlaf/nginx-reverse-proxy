# nginx-reverse-proxy
A little nginx-reverse-proxy docker image with just enough automated setup to be reasonably convenient and nothing more.

I've found variations of something like this to be pretty useful for my own stuff a bunch over the years, and figured I'd formalize some of it here for myself going forward.

### Goals
- Also handle cert creation/renewal for anything in the configured sites via letsencrypt certbot
- Allow for http to https redirection and www redirection on a per-subdomain basis
- Enable configuring resolver/caching for upstreams; proxy need not have dependencies on other containers to start
- IP whitelisting on per-subdomain basis
- Streams

### Using
Dockerfile:
```
FROM devlaf/nginx-reverse-proxy:0.1

COPY config.json /app/config.json
ENV NGINX_SETUP_TOOL_CFG=/app/config.json
```


Example config.json:
```
{
	"letsencrypt_settings": {
		"cert_name": "example.com",
		"email": "me@example.com"
	},
	"sites": [{
			"domain": "example.com",
			"proxy_to_address": "http://website",
			"proxy_target_is_docker_container": true,
			"forward_http_to_https": true,
			"include_and_redirect_www": true
		},
		{
			"domain": "blog.example.com",
			"proxy_to_address": "http://wiki",
			"proxy_target_is_docker_container": true,
			"forward_http_to_https": true,
			"include_and_redirect_www": false,
			"restrict_access_to_ip": "127.0.0.1"
		}
	],
	"streams": [{
		"host_port": 22000,
		"proxy_to_address": "http://blah",
		"proxy_to_port": 22000,
		"proxy_target_is_docker_container": true
	}]
}
```
