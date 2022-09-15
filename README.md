# discRSS


## TODO

- [ ] cdk stack to host app on lambda
  - [x] DDB tables
  - [x] lambda
  - [ ] lambda endpoints
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