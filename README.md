# ported
Ported - Self Hosted HTTP Tunnel for in-house devs  \
An open source self-hosted (in-house) alternative to ngrok

## Server Side
### Porter service
This is responsible for creating external side of the tunnel (other side being the developers machine). It also configures services in the external gateway server to expose the tunnel to an accessible endpoint.
```sh
./porter -rootDomain ported.mydomain.com
```

## Developer side
### Exposing local HTTP service
```sh
./ported -using http://ported.mydomain.com -to localhost:8080
```
