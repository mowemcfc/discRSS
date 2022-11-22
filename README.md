# discRSS

## Usage

Rename `provision/local/sample_discord_token.txt` to `discord_token.txt` and replace the text with your own Discord bot token.

Replace the initialisation data in `provision/lib/discRSSStack.ts` with your desired channels and feed subcriptions. (TODO: provide a nicer interface to configure this)

Run the init script: (TODO: think about how to remove profile-based argument)

```sh
make deploy
```

## TODO

### Priority

- [ ] introduce more granular routes to API handlers e.g. user -> user/feed, user/channel
- [ ] add feed button
  - [x] ui component
  - [x] onSubmit handler
    - [x] local data
    - [x] post request to API route
  - [x] backend handler
- [ ] remove feed button
  - [ ] ui component
  - [ ] onsubmit handler
    - [ ] post request to API route
    - [ ] set local state
  - [ ] backend handler
- [ ] simplify timestamp parsing by removing format-specific logic, see time.Parse(layout, str)
- [ ] remove timestamp from frontend

### Backlog 

- [ ] validate submitted feed are indeed valid RSS
- [ ] tracing
- [ ] make as much of lambda async as possible
- [ ] logging
  - [x] basic
  - [ ] structured
- [ ] consolidate .env files and other assorted local txt's
- [ ] lambda JWT auth
  - [x] basic
  - [ ] allow proper scopes, createUser, getFeed
  - [ ] prevent IDOR by comparing userID's against JWT `sub`

### Completed

- [x] cdk stack to host app on lambda
  - [x] DDB tables
  - [x] lambda
  - [x] apigw endpoints
    - [x] user
- [x] make table retention policy DESTROY
- [x] convert go into lambda
- [x] start sessions in separate functions
  - [x] aws
  - [x] discord
- [x] write go SDK code to fetch data at runtime
- [x] secrets management
- [x] deal with onUpdate for AWS CR's
  - [x] discord secret
  - [x] ddb table
- [x] DB (SQL / nosql) schema
- [x] iterate over multiple channels to post in
- [x] eventbridge schedule for lambda
- [x] update lastChecked from code
- [x] interface init data, consumed by cdk
- [x] split user fetch call into separate lambda
  - [x] call this lambda from cronned lambda
- [x] make aws session global
- [x] init lastCheckedTime in CDK as time.now()
- [x] fix autoscan
- [x] figure out calling scan endpoint internally (probably just separate it into separate lambda)
- [x] clean up newly-separated go modules
- [x] login page
  - [x] make login page redirect to profile if already logged in
- [x] auth0 integration
- [x] display feeds via ddb fetch
  - [x] fetch data
  - [x] ui components
- [x] display channels via ddb fetch
- [x] create DNS endpoint for API
  - [x] update frontend refs
- [x] fix eventbridge call to lambda IOT use /scan endpoint
- [x] add .env placeholders throughout where used
- [x] user POST handler