## Silk.web

Webserver for logging user interaction with the client's website

## Logging
All saved data is stored across different files in <root path>/filedb

Logged errors are stored in <root path>/log

## Event schema
`{uuid: "", event_type: "...", user_id: "...", page_url:"...", payload: "..."}`

## Protocol
#### GET /events/{:uuid}
Always returns 200 and {"saved":true} if the event persisted, {"saved":false} - otherwise

Tracks last 10000 events

#### POST /events
Requirements:
- JSON string with params in request body
- uuid field is present and unique (across the entire JSON string)
- timestamp, ip address and user agent are added on the server side, **don't include** these fields into request body 

Always returns 200, on error - saves to log


## Auth

- Secret is passed in X-Secret header
- Domain-secret pairs are stores in redis "silk:auth:keys" hash
- For staging use 12 database, for production - 13 database in Redis, e.g. `SELECT 9` in redis-cli
- Auth is skipped for the /script route, as this is where the credentials are stored

## Script

The client script is stored in a redis queue silk:js:script (each element corresponds to deploy version), 
the client config json is stored in silk:js:config hash (the key is domain, the value is json config).
 
The script is served based on the origin header, it merges config silk:js:config with silk:js:script, wraps
it into immediate function call and serves as a js script. 
Inside the JS function the config is available via the config variable.

The route is /script

## Deploy

- Deploy is сonfigured for ubuntu only
- To deploy to other OS - tweak daemon, sudoers and ARCH and OS for build
- Usage - `deploy/deploy <staging|production>`
- Deploy (root) path is /home/delploy/apps/silk.web/<production|staging>/current
- To export data use `service silk.web export`, the data is stored in <root path>/filedb/export