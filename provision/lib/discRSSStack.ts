import { Duration, RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cr from 'aws-cdk-lib/custom-resources';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb'
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager'
import * as lambda from 'aws-cdk-lib/aws-lambda'
import * as apigateway from 'aws-cdk-lib/aws-apigateway'
import * as eventbridge from 'aws-cdk-lib/aws-events'
import * as eventtargets from 'aws-cdk-lib/aws-events-targets'
import * as acm from 'aws-cdk-lib/aws-certificatemanager'
import * as route53 from 'aws-cdk-lib/aws-route53'
import * as targets from 'aws-cdk-lib/aws-route53-targets'

import * as fs from 'fs';
import * as path from 'path';
import { ServicePrincipal } from 'aws-cdk-lib/aws-iam';

import 'datejs'
import { TargetTrackingScalingPolicy } from 'aws-cdk-lib/aws-applicationautoscaling';

export class DiscRssStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const scanLambda = new lambda.Function(this, 'DiscRSS-ScanLambda', {
      functionName: 'discRSS-scanLambda',
      architecture: lambda.Architecture.X86_64,
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset(path.join(__dirname, '../../backend/scanner.zip')),
      handler: 'main',
      timeout: Duration.seconds(60)
    })

    const lambdaScheduledExecution = new eventbridge.Rule(this, 'DiscRSS-ScanLambdaScheduledExecution', {
      schedule: eventbridge.Schedule.cron({ minute: '0/30' })
    })

    lambdaScheduledExecution.addTarget(
      new eventtargets.LambdaFunction(scanLambda, {
        event: eventbridge.RuleTargetInput.fromObject({userID: "1"})
      })
    )


    const apiLambda = new lambda.Function(this, 'DiscRSS-ApiLambda', {
      functionName: 'discRSS-apiLambda',
      architecture: lambda.Architecture.X86_64,
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset(path.join(__dirname, '../../backend/api.zip')),
      handler: 'main',
      timeout: Duration.seconds(60)
    })

    const discRSSHostedZone = route53.HostedZone.fromHostedZoneAttributes(this, 'discRSS-HostedZone', {
      zoneName: 'discrss.com',
      hostedZoneId: 'Z01872172J0T33M526LB9'
    })

    const apiCertificate = new acm.Certificate(this, 'discRSS-Certificate', {
      domainName: 'discrss.com',
      subjectAlternativeNames: [ '*.discrss.com' ],
      validation: acm.CertificateValidation.fromDns(discRSSHostedZone)
    })


    const discRSSApi = new apigateway.LambdaRestApi(this, 'DiscRSS-API', {
      handler: apiLambda,
      deploy: true,
      proxy: true,
      domainName: {
        domainName: 'api.discrss.com',
        certificate: apiCertificate
      }
    })

    const apiDnsRecord = new route53.ARecord(this, 'discRSS-ApiDnsRecord', {
      zone: discRSSHostedZone,
      recordName: 'api.discrss.com',
      target: route53.RecordTarget.fromAlias(new targets.ApiGateway(discRSSApi))
    })

    apiLambda.addPermission('DiscRSS-AllowAPIGWInvocation', {
      principal: new ServicePrincipal('apigateway.amazonaws.com'),
      sourceArn: discRSSApi.arnForExecuteApi('*'),
    })

    const userTable = new dynamodb.Table(this, 'DiscRSS-UserTable', {
      tableName: 'discRSS-UserRecords',
      partitionKey: {
        name: 'userId',
        type: dynamodb.AttributeType.NUMBER
      }
    })
    userTable.applyRemovalPolicy(RemovalPolicy.DESTROY)
    userTable.grantReadWriteData(apiLambda.role!.grantPrincipal)
    userTable.grantReadWriteData(scanLambda.role!.grantPrincipal)
    
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
    appConfigTable.grantReadWriteData(scanLambda.role!.grantPrincipal)
    appConfigTable.grantReadWriteData(apiLambda.role!.grantPrincipal)

    // Use for debugging, if you want to test your feed update over a longer period
    //const currentTime = "2022-08-02T15:04:05Z"
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

    const auth0ClientSecret = new secretsmanager.Secret(this, 'DiscRSS-Auth0ClientSecret', {
      secretName: 'discRSS/auth0-client-secret'
    })
    auth0ClientSecret.grantRead(apiLambda.role!.grantPrincipal)

    const putAuth0ClientSecretAction: cr.AwsSdkCall = {
        service: 'SecretsManager',
        action: 'putSecretValue',
        parameters: {
          SecretId: auth0ClientSecret.secretName,
          SecretString: fs.readFileSync(path.join(__dirname, '../local/auth0_client_secret.txt'), { encoding: 'utf-8' })
        },
        physicalResourceId: cr.PhysicalResourceId.of(auth0ClientSecret.secretName + '_initialisation')
    }

    const auth0ClientSecretUpdateCr = new cr.AwsCustomResource(this, 'DiscRSS-Auth0ClientSecretUpdate', {
      onCreate: putAuth0ClientSecretAction,
      onUpdate: putAuth0ClientSecretAction,
      policy: cr.AwsCustomResourcePolicy.fromSdkCalls({ resources: cr.AwsCustomResourcePolicy.ANY_RESOURCE }),
    }) 

    const discordBotSecret = new secretsmanager.Secret(this, 'DiscRSS-DiscordBotSecret', {
      secretName: 'discRSS/discord-bot-secret',
    })
    discordBotSecret.grantRead(scanLambda.role!.grantPrincipal)

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
