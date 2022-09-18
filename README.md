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
  - [ ] split user fetch call into separate lambda
    - [ ] call this lambda from cronned lambda
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
- [ ] make aws session global
- [ ] init lastCheckedTime in CDK as time.now()
- [x] eventbridge schedule for lambda
- [x] update lastChecked from code
- [ ] make as much of lambda async as possible
- [ ] logging
  - [x] basic
  - [ ] structured
- [ ] tracing
- [ ] frontend
  - [ ] hardcode userid=1
  - [ ] display feeds via ddb fetch
  - [ ] display channels via ddb fetch
  - [ ] allow putitem new feeds via frontend form 
  - [ ] putitem new channels via frontend form