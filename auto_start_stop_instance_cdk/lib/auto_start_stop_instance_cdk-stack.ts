import events = require('@aws-cdk/aws-events');
import iam = require('@aws-cdk/aws-iam');
import targets = require('@aws-cdk/aws-events-targets');
import { Function, Runtime, Code } from "@aws-cdk/aws-lambda"
import cdk = require('@aws-cdk/core');
import fs = require('fs');

export class AutoStartStopInstanceCdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const stackConfig = JSON.parse(fs.readFileSync('example.json', {encoding: 'utf-8'}));

    const lambdaFn = new Function(this, 'singleton', {
      functionName: "auto_ec2_startstop", // 関数名
      runtime: Runtime.GO_1_X, // ランタイムの指定
      code: Code.asset("./lambdaSources/source"), // ソースコードのディレクトリ
      handler: "main", // handler の指定
      memorySize: 256, // メモリーの指定
      timeout: cdk.Duration.seconds(10), // タイムアウト時間
    });

    lambdaFn.addToRolePolicy(new iam.PolicyStatement({
      actions: [
        'ec2:DescribeInstances',
        'ec2:StartInstances',
        'ec2:StopInstances',
        "rds:StopDBInstance",
        "rds:StartDBInstance",
        "rds:StartDBCluster",
        "rds:StopDBCluster"
      ],
      resources: ['*']
    }));

    // STOP EC2 instances rule
    const stopRule = new events.Rule(this, 'StopRule', {
      schedule: events.Schedule.expression(`cron(${stackConfig.events.cron.stop})`)
    });

    stopRule.addTarget(new targets.LambdaFunction(lambdaFn, {
      event: events.RuleTargetInput.fromObject({Region: stackConfig.targets.ec2region, Action: 'stop',InstanceNames:stackConfig.targets.instanceNames
        ,DBInstanceIdentifies:stackConfig.targets.dbInstanceIdentifies,DBClusterIdentifies:stackConfig.targets.dbClusterIdentifies
        ,WebHookURL:stackConfig.slack.webHookUrl,SlackChannel:stackConfig.slack.channel})
    }));

    // START EC2 instances rule
    const startRule = new events.Rule(this, 'StartRule', {
      schedule: events.Schedule.expression(`cron(${stackConfig.events.cron.start})`)
    });

    startRule.addTarget(new targets.LambdaFunction(lambdaFn, {
      event: events.RuleTargetInput.fromObject({Region: stackConfig.targets.ec2region, Action: 'start',InstanceNames:stackConfig.targets.instanceNames
        ,DBInstanceIdentifies:stackConfig.targets.dbInstanceIdentifies,DBClusterIdentifies:stackConfig.targets.dbClusterIdentifies
        ,WebHookURL:stackConfig.slack.webHookUrl,SlackChannel:stackConfig.slack.channel})
    }));

  }
}

