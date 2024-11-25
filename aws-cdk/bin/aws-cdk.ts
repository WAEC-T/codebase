#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { RDSStack } from '../lib/rds-stack';
import { BucketResourceStack } from '../lib/s3-stack';

const commonTags = { Application: 'WAECT-T' }

const app = new cdk.App();

const rdsStack = new RDSStack(app, 'waectDATABASE', {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
});

cdk.Tags.of(rdsStack).add(commonTags.Application, 'WAECT-T');

const bucketResourceStack = new BucketResourceStack(app, 'waectS3Bucket', {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
})

cdk.Tags.of(bucketResourceStack).add(commonTags.Application, 'WAECT-T');

