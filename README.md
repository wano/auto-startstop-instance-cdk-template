# Auto-startstop-instance-cdk-template
## What is this
- Auto-start-stop-instance(EC2,RDS)-cdk-template


## What can you do with this template??
- you can reduce costs by stopping instances(EC2,RDS) while not in use
- you can choice configure to suit your case
    - only ec2 instance,rds cluster etc
    - multiple ec2 and rds instances support
- you can notice to slack optionaly

## How to use
### 1.Clone this project  
`git clone https://github.com/wano/auto-startstop-instance-cdk-template.git`
### 2.Modify auto_start_stop_instance_cdk/example.json to suit your case
- Note that the crone time is set in UTC

| |  ec2|  rds(cluster) | rds(instance)|Slack Notification|
|---|----|---|---|---|
|case1| ◯ |　◯ | ◯ | ◯ |
|case2| ◯ |　  | ◯ | ◯ |
|case3|　 |　  | ◯　|  |

- case1

```json
{
  "events": {
    "cron": {
      "start": "00 00 ? * MON-FRI *", 
      "stop": "00 11 ? * * *"
    }
  },
  "targets": {
    "ec2region": "your_region",
    "instanceNames": ["your_instance_name_1","your_instance_name_ 2"],
    "dbClusterIdentifies": ["your_cluster_name_1"],
    "dbInstanceIdentifies": ["your_db_instance_name_1","your_db_instance_name_ 2"]
  },
  "slack": {
    "webHookUrl":"https://hooks.slack.com/services/xxxxxxxxx/xxxxxxxxxxxxx/xxxxxxxxxxxxxxxxxx",
    "channel":"xxxxxxxxxxxxxx"
  }
}
```

- case2

```json
{
  "events": {
    "cron": {
      "start": "00 00 ? * MON-FRI *", 
      "stop": "00 11 ? * * *"
    }
  },
  "targets": {
    "ec2region": "your_region",
    "instanceNames": ["your_instance_name_1","your_instance_name_ 2"],
    "dbClusterIdentifies": [""],
    "dbInstanceIdentifies": ["your_db_instance_name_1","your_db_instance_name_ 2"]
  },
  "slack": {
    "webHookUrl":"https://hooks.slack.com/services/xxxxxxxxx/xxxxxxxxxxxxx/xxxxxxxxxxxxxxxxxx",
    "channel":"xxxxxxxxxxxxxx"
  }
}
```

- case3

```json
{
  "events": {
    "cron": {
      "start": "00 00 ? * MON-FRI *",
      "stop": "00 11 ? * * *"
    }
  },
  "targets": {
    "ec2region": "ap-northeast-1",
    "instanceNames": [""],
     "dbClusterIdentifies": ["your_cluster_name_1"],
    "dbInstanceIdentifies": [""]
  },
  "slack": {
    "webHookUrl":"",
    "channel":""
  }
}
```

### 3.Install CDK and deploy with following commands
```$xslt
npm install -g aws-cdk
cd ./auto_start_stop_instance_cdk
npm install @aws-cdk/core 
npm install @aws-cdk/aws-lambda @aws-cdk/aws-iam @aws-cdk/aws-events @aws-cdk/aws-events-targets
npm run build
export AWS_ACCESS_KEY_ID="xxxxxxxxxxx"
export AWS_SECRET_ACCESS_KEY="xxxxxxxxxxx"
cdk bootstrap
cdk deploy
```
## Test Coverage 
Following events already tested
- All start(EC2,RDS)
- All stop(EC2,RDS)
- EC2Instance and DBInstance start
- Only EC2Instance start and not notice to slack
## Environment I developed in
- MacOS Mojave version 10.14.6
- Go 1.12.4
- TypeScript 3.6.2
- CDK z1.6.0


