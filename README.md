# discRSS


## TODO

- [ ] cdk stack to host app on lambda
  - [x] DDB tables
  - [ ] lambda
  - [ ] lambda endpoints
- [x] make table retention policy DESTROY
- [ ] convert go into lambda
- [x] start sessions in separate functions
  - [x] aws
  - [x] discord
- [x] write go SDK code to fetch data at runtime
- [ ] secrets management
- [ ] deal with onUpdate for AWS CR's
- [x] DB (SQL / nosql) schema
- [x] iterate over multiple channels to post in
- [ ] frontend
  - [ ] hardcode userid=1
  - [ ] display feeds via ddb fetch
  - [ ] display channels via ddb fetch
  - [ ] allow putitem new feeds via frontend form 
  - [ ] putitem new channels via frontend form