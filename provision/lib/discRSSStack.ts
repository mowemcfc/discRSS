import { RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cr from 'aws-cdk-lib/custom-resources';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb'

export class DiscRssStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const userTable = new dynamodb.Table(this, 'UserTable', {
      tableName: 'discRSS-UserRecords',
      partitionKey: {
        name: 'userID',
        type: dynamodb.AttributeType.NUMBER
      }
    })

    userTable.applyRemovalPolicy(RemovalPolicy.DESTROY)
    
    const userTableInit = new cr.AwsCustomResource(this, 'initTable', {
      onCreate: {
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
            }
        },
        physicalResourceId: cr.PhysicalResourceId.of(userTable.tableName + '_initialization')
      },
      policy: cr.AwsCustomResourcePolicy.fromSdkCalls({ resources: cr.AwsCustomResourcePolicy.ANY_RESOURCE }),
    });

    const appconfigTable = new dynamodb.Table(this, 'AppconfigTable', {
      tableName: 'discRSS-AppConfigs',
      partitionKey: {
        name: 'configID',
        type: dynamodb.AttributeType.NUMBER
      }
    })
    appconfigTable.applyRemovalPolicy(RemovalPolicy.DESTROY)
  }
}
