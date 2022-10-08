# Testing this lambda locally

The following will setup your lambda in "listen" mode on localhost:9000

```
docker run --rm \
  -e DOCKER_LAMBDA_STAY_OPEN=1 \ 
  -p 9001:9001 \
  -v /Users/jessecarter/pers/discRSS/src/backend:/var/task:ro,delegated \
  lambci/lambda:go1.x \
  main
```

You can invoke it as follows:

```
curl -vv -d '{}' http://localhost:9001/2015-03-31/functions/myfunction/invocations
```

or

```
aws lambda invoke --endpoint http://localhost:9001 --no-sign-request \
  --function-name myfunction --cli-binary-format raw-in-base64-out --payload '{}' output.json
```
