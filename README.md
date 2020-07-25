# ported
Ported - Self Hosted HTTP Tunnel for in-house devs  \
An open source self-hosted (in-house) alternative to ngrok

## Server Side
### Porter service
This is responsible for configuring services in the external gateway server (which will tunnel to the developers local services)
```sh
./porter -rootDomain ported.mydomain.com
```

## Developer side
### Exposing local HTTP service
```sh
./ported -using http://ported.mydomain.com -to localhost:8080
```
