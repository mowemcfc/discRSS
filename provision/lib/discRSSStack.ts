import { Duration, RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cr from 'aws-cdk-lib/custom-resources';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb'
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager'
import * as lambda from 'aws-cdk-lib/aws-lambda'
import * as apigateway from 'aws-cdk-lib/aws-apigateway'
import * as eventbridge from 'aws-cdk-lib/aws-events'
import * as eventtargets from 'aws-cdk-lib/aws-events-targets'
import * as fs from 'fs';
import * as path from 'path';
import { ServicePrincipal } from 'aws-cdk-lib/aws-iam';

import 'datejs'

export class DiscRssStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const discRSSLambda = new lambda.Function(this, 'DiscRSSLambda', {
      functionName: 'discRSSLambda',
      architecture: lambda.Architecture.X86_64,
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset(path.join(__dirname, '../../src/backend/main.zip')),
      handler: 'main',
      timeout: Duration.seconds(60)
    })

    const userApi = new apigateway.RestApi(this, 'DiscRSS-UserAPI', {
      restApiName: 'discRSS-UserAPI',
      deploy: true,
      deployOptions: {
        stageName: 'v1'
      },
      defaultCorsPreflightOptions: {
        allowHeaders: [
          '*',
          'Authorization'
        ],
        allowMethods: apigateway.Cors.ALL_METHODS,
        allowCredentials: true,
        allowOrigins: ['*']
      },
    })
    
    const userApiLambdaIntegration = new apigateway.LambdaIntegration(
      discRSSLambda, {
        contentHandling: apigateway.ContentHandling.CONVERT_TO_TEXT,
      }
    )
    userApi.root.addProxy({ defaultIntegration: userApiLambdaIntegration })

    discRSSLambda.addPermission('DiscRSS-AllowAPIGWInvocation', {
      principal: new ServicePrincipal('apigateway.amazonaws.com'),
      sourceArn: userApi.arnForExecuteApi('*'),
    })

    const lambdaScheduledExecution = new eventbridge.Rule(this, 'DiscRSS-LambdaScheduledExecution', {
      schedule: eventbridge.Schedule.cron({ minute: '0/30' })
    })

    lambdaScheduledExecution.addTarget(
      new eventtargets.ApiGateway(userApi, {
        method: 'GET',
        path: '/scan',
        queryStringParameters: {
          'userID': '1'
        }
      })
    )

    const userTable = new dynamodb.Table(this, 'DiscRSS-UserTable', {
      tableName: 'discRSS-UserRecords',
      partitionKey: {
        name: 'userID',
        type: dynamodb.AttributeType.NUMBER
      }
    })
    userTable.applyRemovalPolicy(RemovalPolicy.DESTROY)
    userTable.grantReadWriteData(discRSSLambda.role!.grantPrincipal)
    
    // Load initialised feed/channel data from file. Feel free to replace this data with your own if you like
    const initData = fs.readFileSync(path.join(__dirname, '../local/init_data.json'), { encoding: 'utf-8' })

    const userTableInitAction: cr.AwsSdkCall = {
      service: 'DynamoDB',
      action: 'batchWriteItem',
      parameters: {
        RequestItems: JSON.parse(initData)
      },
      physicalResourceId: cr.PhysicalResourceId.of(userTable.tableName + '_initialization')
    }
    const userTableInit = new cr.AwsCustomResource(this, 'initTable', {
      onCreate: userTableInitAction,
      onUpdate: userTableInitAction,
      policy: cr.AwsCustomResourcePolicy.fromSdkCalls({ resources: cr.AwsCustomResourcePolicy.ANY_RESOURCE }),
    });

    const appConfigTable = new dynamodb.Table(this, 'DiscRSS-AppConfigTable', {
      tableName: 'discRSS-AppConfigs',
      partitionKey: {
        name: 'appName',
        type: dynamodb.AttributeType.STRING
      }
    })
    appConfigTable.applyRemovalPolicy(RemovalPolicy.DESTROY)
    appConfigTable.grantReadWriteData(discRSSLambda.role!.grantPrincipal)

    const currentTime = new Date().toISOString()
    const appConfigTableInitAction: cr.AwsSdkCall = {
      service: 'DynamoDB',
      action: 'putItem',
      parameters: {
        TableName: appConfigTable.tableName,
        Item: { 
          appName: { S: "discRSS" },
          lastCheckedTime: { S: currentTime }, 
          lastCheckedTimeFormat: { S: "2006-01-02T15:04:05Z07:00" }
        }
      },
      physicalResourceId: cr.PhysicalResourceId.of(appConfigTable.tableName + '_initialization')
    }

    const appConfigTableInit = new cr.AwsCustomResource(this, 'DiscRSS-InitAppConfigTable', {
      onCreate: appConfigTableInitAction,
      onUpdate: appConfigTableInitAction,
      policy: cr.AwsCustomResourcePolicy.fromSdkCalls({ resources: cr.AwsCustomResourcePolicy.ANY_RESOURCE })
    })

    const discordBotSecret = new secretsmanager.Secret(this, 'DiscRSS-DiscordBotSecret', {
      secretName: 'discRSS/discord-bot-secret',
    })
    discordBotSecret.grantRead(discRSSLambda.role!.grantPrincipal)

    const putDiscordBotSecretAction: cr.AwsSdkCall = {
        service: 'SecretsManager',
        action: 'putSecretValue',
        parameters: {
          SecretId: discordBotSecret.secretName,
          SecretString: fs.readFileSync(path.join(__dirname, '../local/discord_token.txt'), { encoding: 'utf-8' })
        },
        physicalResourceId: cr.PhysicalResourceId.of(discordBotSecret.secretName + '_initialisation')
    }

    const discordBotTokenUpdateCr = new cr.AwsCustomResource(this, 'DiscRSS-DiscordBotSecretUpdate', {
      onCreate: putDiscordBotSecretAction,
      onUpdate: putDiscordBotSecretAction,
      policy: cr.AwsCustomResourcePolicy.fromSdkCalls({ resources: cr.AwsCustomResourcePolicy.ANY_RESOURCE }),
    }) 
  }
}
