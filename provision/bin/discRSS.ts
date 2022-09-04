#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { DiscRssStack } from '../lib/discRSSStack';

const app = new cdk.App();
new DiscRssStack(app, 'DiscRssStack', {
  env: { account: process.env.CDK_DEFAULT_ACCOUNT, region: 'ap-southeast-2' },
});