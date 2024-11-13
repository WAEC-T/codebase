#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { RDSStack } from '../lib/rds-stack';

const commonTags = { Application: 'WAECT-T' }

const app = new cdk.App();

const RDSSTack = new RDSStack(app, 'waectDATABASE', {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT, // Use environment variables for flexibility
    region: process.env.CDK_DEFAULT_REGION,
  },
});

cdk.Tags.of(RDSSTack).add(commonTags.Application, 'WAECT-T');