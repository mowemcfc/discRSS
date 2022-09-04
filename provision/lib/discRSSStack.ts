import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as cdk from 'aws-cdk-lib'
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb'

export class DiscRssStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const feedTable = new dynamodb.Table(this, 'FeedTable', {
      tableName: 'discRSS-feeds',
      partitionKey: {
        name: 'feedID',
        type: dynamodb.AttributeType.NUMBER
      }
    })

    const appconfigTable = new dynamodb.Table(this, 'AppconfigTable', {
      tableName: 'discRSS-appConfigs',
      partitionKey: {
        name: 'configID',
        type: dynamodb.AttributeType.NUMBER
      }
    })

    const channelTable = new dynamodb.Table(this, 'ChannelTable', {
      tableName: 'discRSS-channels',
      partitionKey: {
        name: 'channelID',
        type: dynamodb.AttributeType.NUMBER
      }
    })
  }
}
