# go-gluons

A grab bag of personal go packages, and templates to hold things together

### Stable gluons:

- **log**: A super simple, traditional logging with levels
- **fileserver:** https://github.com/prasannavl/go-gluons/tree/master/http/fileserver - Reimplemented Go's http file server that properly returns errors instead of having it's logic inter-mingled. This allows nice directory listing handling, and error handling with ease.
- **hostrouter:** https://github.com/prasannavl/go-gluons/tree/master/http/hostrouter - A router that handles hosts switching between the most efficient representations on the fly.
- **ansicode:** https://github.com/prasannavl/go-gluons/tree/master/ansicode - Common ansi-code as simple constants, and helpers. No mutexes like other color libraries. Minimal form for high-performance usage areas.  
- **cert:** https://github.com/prasannavl/go-gluons/tree/master/cert - Simple helper methods to create X509 certs, including self-signed certs on the fly with just a simple call. 
- **writer:** https://github.com/prasannavl/go-gluons/tree/master/http/writer - A http `ResponseWriter` impl that's simple and handles reliable lazy status writes and with other helpers like `Tee`, `IsHijacked`, `IsStatusWritten`, `IsStatusSet` etc.

### Ever-changing gluons:
- **handlerutils:** https://github.com/prasannavl/go-gluons/tree/master/http/handlerutils - Handler helpers that ease a lot of boiler plate for common cases.
- **chainutils:** https://github.com/prasannavl/go-gluons/tree/master/http/chainutils - Middleware chaining helpers that ease boilerplate.
- **middleware:** https://github.com/prasannavl/go-gluons/tree/master/http/middleware - Some middlewares that are helpful.  

### Other useful packages:

- **go-errors**: https://github.com/prasannavl/go-errors - Go error handling helpers.
- **mchain**: https://github.com/prasannavl/mchain - Go http middlewares and chaining helpers with idiomatic error handling.
- **mroute**: https://github.com/prasannavl/mroute - A fork of goji router for mchain with addons.
- **mrouter**: https://github.com/prasannavl/mrouter - A fork of httprouter for mchain.

