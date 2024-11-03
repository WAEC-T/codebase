import * as cdk from 'aws-cdk-lib';
import { Duration, RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import {
  InstanceClass,
  InstanceSize,
  InstanceType,
  Peer,
  Port,
  SecurityGroup,
  SubnetType,
  Vpc,
} from "aws-cdk-lib/aws-ec2";
import {
  Credentials,
  DatabaseInstance,
  DatabaseInstanceEngine,
  PostgresEngineVersion
} from "aws-cdk-lib/aws-rds";
import { Construct } from "constructs";
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';

export class RDSStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const engine: cdk.aws_rds.IInstanceEngine = DatabaseInstanceEngine.postgres({ version: PostgresEngineVersion.VER_16_4 })
    const instanceType: cdk.aws_ec2.InstanceType = InstanceType.of(InstanceClass.T3, InstanceSize.MICRO);
    const port: number = 5432;
    const dbName: string = "waect"; 

    // create a VPC (Virtual Private Cloud)
    const vpc = Vpc.fromLookup(this, 'waectVPC', { vpcId: "vpc-08d145893bcbe80f7" });

    // create a security group
    const dbSg: cdk.aws_ec2.SecurityGroup = new SecurityGroup(this, "waectSG", {
      securityGroupName: "waectSG",
      vpc: vpc
    });

    // Allow incoming traffic from any IP address on port 5432
    dbSg.addIngressRule(
      Peer.anyIpv4(),
      Port.tcp(5432),
      'Allow PostgreSQL access from port 5432'
    );

    // Create a custom password
    const customPassword: cdk.SecretValue = cdk.SecretValue.unsafePlainText(process.env.AWS_DATABASE_PASSWORD || ''); // Replace with a secure password

    // Use the custom password in credentials
    const credentials: cdk.aws_rds.Credentials = Credentials.fromPassword('waect', customPassword);

    // create RDS instance (PostgreSQL)
    const dbInstance: cdk.aws_rds.DatabaseInstance = new DatabaseInstance(this, "waectDB", {
      vpc: vpc,
      vpcSubnets: { subnetType: SubnetType.PUBLIC },
      instanceType,
      engine,
      port,
      securityGroups: [dbSg],
      databaseName: dbName,
      credentials: credentials,
      backupRetention: Duration.days(7),
      deleteAutomatedBackups: false,
      removalPolicy: RemovalPolicy.DESTROY,
      allocatedStorage: 20
    });
    // Capture the endpoint
    const hostname: string = dbInstance.instanceEndpoint.hostname;

    // Create a new Secrets Manager secret with the endpoint included
    new secretsmanager.Secret(this, 'waectSECRET', {
      secretName: 'waectSECRET',
      secretStringValue: cdk.SecretValue.unsafePlainText(JSON.stringify({
          username: 'waect',
          password: process.env.AWS_DATABASE_PASSWORD || '', 
          hostname: hostname,
          port: port,
          database: dbName,
      })),
    });

    // Output the database endpoint for reference
    new cdk.CfnOutput(this, 'hostname', {
      value: hostname,
    });
  }
}
