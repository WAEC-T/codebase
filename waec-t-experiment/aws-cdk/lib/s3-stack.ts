import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { Construct } from "constructs";

export class BucketResourceStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        new s3.Bucket(this, 'waectbucket', {
            bucketName: 'waectbucket',
            publicReadAccess: true, // Allow public read access
            blockPublicAccess: s3.BlockPublicAccess.BLOCK_ACLS, // Allow public access by disabling ACL blocking
            removalPolicy: cdk.RemovalPolicy.RETAIN, // Retain bucket when stack is destroyed
        });
    }
}