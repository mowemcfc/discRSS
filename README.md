# discRSS

## Usage

Rename `provision/local/sample_discord_token.txt` to `discord_token.txt` and replace the text with your own Discord bot token.

Replace the initialisation data in `provision/lib/discRSSStack.ts` with your desired channels and feed subcriptions. (TODO: provide a nicer interface to configure this)

Run the init script: (TODO: think about how to remove profile-based argument)

```sh
make deploy
```

## TODO

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
- [ ] tracing
- [ ] frontend
  - [x] login page
   - [x] make login page redirect to profile if already logged in
  - [x] auth0 integration
  - [x] display feeds via ddb fetch
    - [x] fetch data
    - [x] ui components
  - [x] display channels via ddb fetch
  - [ ] allow putitem new feeds via frontend form  <--
    - [x] ui component
    - [ ] onSubmit handler
    - [ ] API route
  - [ ] putitem new channels via frontend form
- [x] create DNS endpoint for API
  - [ ] update frontend refs
- [x] fix eventbridge call to lambda IOT use /scan endpoint
- [x] add .env placeholders throughout where used
- [ ] make as much of lambda async as possible
- [ ] logging
  - [x] basic
  - [ ] structured
- [ ] simplify timestamp parsing by removing format-specific logic, see time.Parse(layout, str)
- [ ] consolidate .env files and other assorted local txt's
- [ ] lambda JWT auth
  - [x] basic
  - [ ] allow proper scopes, createUser, getFeed
  - [ ] compare provided username against jwt sub, throw error if mismatch