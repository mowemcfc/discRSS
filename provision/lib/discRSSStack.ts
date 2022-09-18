import { Duration, RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cr from 'aws-cdk-lib/custom-resources';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb'
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager'
import * as lambda from 'aws-cdk-lib/aws-lambda'
import * as eventbridge from 'aws-cdk-lib/aws-events'
import * as eventtargets from 'aws-cdk-lib/aws-events-targets'
import * as fs from 'fs';
import * as path from 'path';

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

    const lambdaScheduledExecution = new eventbridge.Rule(this, 'lambdaScheduledExecution', {
      schedule: eventbridge.Schedule.cron({ minute: '0/10' })
    })

    lambdaScheduledExecution.addTarget(
      new eventtargets.LambdaFunction(discRSSLambda, {})
    )

    const userTable = new dynamodb.Table(this, 'UserTable', {
      tableName: 'discRSS-UserRecords',
      partitionKey: {
        name: 'userID',
        type: dynamodb.AttributeType.NUMBER
      }
    })
    userTable.applyRemovalPolicy(RemovalPolicy.DESTROY)
    userTable.grantReadWriteData(discRSSLambda.role!.grantPrincipal)
    
    const userTableInitAction: cr.AwsSdkCall = {
      service: 'DynamoDB',
      action: 'putItem',
      parameters: {
        TableName: userTable.tableName,
        Item: { 
          userID: { N: "1" },
          username: { S: "mowemcfc" }, 
          feedList: { L: [
            {
              M: { 
                feedID: { N: "1" }, 
                title: { S: "The Future Does Not Fit In The Containers Of The Past" },
                url: {S: "https://rishad.substack.com/feed" }, 
                timeFormat: { S: "Mon, 02 Jan 2006 15:04:05 MST" },
              }
            },
            {
              M: { 
                feedID: { N: "2" }, 
                title: { S: "Dan Luu" },
                url: {S: "https://danluu.com/atom.xml" }, 
                timeFormat: { S: "Mon, 02 Jan 2006 15:04:05 -0700" },
              }
            },
          ]},
          channelList: { L: [
            {
              M: { 
                channelID: { N: "985831956203851786" }, 
                channelName: { S: "mowes mate" },
                serverName: {S: "mines" }, 
              }
            },
            {
              M: { 
                channelID: { N: "958948046606053406" }, 
                channelName: { S: "rss" },
                serverName: {S: "klnkn (pers)" }, 
              }
            },
          ]}
        },
      },
      physicalResourceId: cr.PhysicalResourceId.of(userTable.tableName + '_initialization')
    }
    const userTableInit = new cr.AwsCustomResource(this, 'initTable', {
      onCreate: userTableInitAction,
      onUpdate: userTableInitAction,
      policy: cr.AwsCustomResourcePolicy.fromSdkCalls({ resources: cr.AwsCustomResourcePolicy.ANY_RESOURCE }),
    });

    const appConfigTable = new dynamodb.Table(this, 'AppConfigTable', {
      tableName: 'discRSS-AppConfigs',
      partitionKey: {
        name: 'appName',
        type: dynamodb.AttributeType.STRING
      }
    })
    appConfigTable.applyRemovalPolicy(RemovalPolicy.DESTROY)
    appConfigTable.grantReadWriteData(discRSSLambda.role!.grantPrincipal)

    const appConfigTableInitAction: cr.AwsSdkCall = {
      service: 'DynamoDB',
      action: 'putItem',
      parameters: {
        TableName: appConfigTable.tableName,
        Item: { 
          appName: { S: "discRSS" },
          lastCheckedTime: { S: "2022-08-30T00:00:00+10:00" }, 
          lastCheckedTimeFormat: { S: "2006-01-02T15:04:05Z07:00" }
        }
      },
      physicalResourceId: cr.PhysicalResourceId.of(appConfigTable.tableName + '_initialization')
    }

    const appConfigTableInit = new cr.AwsCustomResource(this, 'initAppConfigTable', {
      onCreate: appConfigTableInitAction,
      onUpdate: appConfigTableInitAction,
      policy: cr.AwsCustomResourcePolicy.fromSdkCalls({ resources: cr.AwsCustomResourcePolicy.ANY_RESOURCE })
    })

    const discordBotSecret = new secretsmanager.Secret(this, 'DiscordBotSecret', {
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

    const discordBotTokenUpdateCr = new cr.AwsCustomResource(this, 'DiscordBotSecretUpdate', {
      onCreate: putDiscordBotSecretAction,
      onUpdate: putDiscordBotSecretAction,
      policy: cr.AwsCustomResourcePolicy.fromSdkCalls({ resources: cr.AwsCustomResourcePolicy.ANY_RESOURCE }),
    }) 
  }
}
