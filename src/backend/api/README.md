# Testing this lambda locally

## Method 1 - Lambda emulation (UNFINISHED)

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

## Method 2 - Local Adapter usage

This method disabled the `ginLambda` adapter and instead serves traditional http requests.

Ensure you configure your `.env` file with the correct AWS named profile and region.

```
go run .
```

You can invoke it as follows:

```
curl -XGET http://localhost:9001/user\?userID\=1
```
