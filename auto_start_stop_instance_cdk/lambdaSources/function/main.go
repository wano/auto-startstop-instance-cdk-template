package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/labstack/gommon/log"
)

func main() {
	lambda.Start(stopOrStartInstancesHandler)
}

const (
	SLACK_MONITORING_ICON = ":ok:"
	SLACK_MONITORING_NAME = "Auto Start-Stop Instance Notice"
)


func stopOrStartInstancesHandler(context context.Context, event map[string]interface{}) (e error) {

	action := event["Action"].(string)
	if action != "start"&& action != "stop"{
		return errors.New("invalid action")
	}
	region := event["Region"].(string)

	switch action {
	case "start":
		if err := startEc2Instances(event["InstanceNames"].([]interface{}),region);err != nil{
			log.Error(err)
			return err
		}

		if err := startRDSInstances(event["DBInstanceIdentifies"].([]interface{}),region);err != nil{
			log.Error(err)
			return err
		}

		if err := startRDSClusters(event["DBClusterIdentifies"].([]interface{}),region);err != nil{
			log.Error(err)
			return err
		}
	case "stop":
		if err := stopEc2Instances(event["InstanceNames"].([]interface{}),region);err != nil{
			log.Error(err)
			return err
		}

		if err := stopRDSInstances(event["DBInstanceIdentifies"].([]interface{}),region);err != nil{
			log.Error(err)
			return err
		}

		if err := stopRDSClusters(event["DBClusterIdentifies"].([]interface{}),region);err != nil{
			log.Error(err)
			return err
		}

	}

	if err := noticeToSlack(event["WebHookURL"].(string),event["SlackChannel"].(string),action);err != nil{
		log.Error(err)
		return err
	}

	return nil
}

func stopEc2Instances(instanceNames []interface{},region string) error {
	if len(instanceNames)==0 {
		return nil
	}


	ec := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	instanceIDs,err := extractTargetEc2InstanceIDs(ec,instanceNames)
	if err != nil {
		return err
	}

	_,err = ec.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})

	if err != nil{
		return err
	}

	return nil

}
func startEc2Instances(instanceNames []interface{},region string) error {
	if len(instanceNames)==0 {
		return nil
	}

	ec := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	instanceIDs,err := extractTargetEc2InstanceIDs(ec,instanceNames)
	if err != nil {
		return err
	}

	_,err = ec.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: instanceIDs,
	})

	if err != nil{
		return err
	}

	return nil

}

func stopRDSInstances(dbInstanceIdentifies []interface{},region string) (e error) {
	if len(dbInstanceIdentifies) == 0 {
		return nil
	}

	rdsClient := rds.New(session.New(), &aws.Config{Region: aws.String(region)})


	for _,dbInstanceIdentify := range dbInstanceIdentifies {
		_, err := rdsClient.StopDBInstance(
			&rds.StopDBInstanceInput{
				DBInstanceIdentifier: aws.String(dbInstanceIdentify.(string)),}, )
		if err != nil {
			return err
		}
	}

	return nil
}

func startRDSInstances(dbInstanceIdentifies []interface{},region string) (e error) {
	if len(dbInstanceIdentifies) == 0 {
		return nil
	}

	rdsClient := rds.New(session.New(), &aws.Config{Region: aws.String(region)})


	for _,dbInstanceIdentify := range dbInstanceIdentifies {
		_, err := rdsClient.StartDBInstance(
			&rds.StartDBInstanceInput{
				DBInstanceIdentifier: aws.String(dbInstanceIdentify.(string)),
			})
		if err != nil {
			return err
		}
	}

	return nil
}

func stopRDSClusters(dbClusterIdentifies []interface{},region string) (e error) {
	if len(dbClusterIdentifies) == 0 {
		return nil
	}

	rdsClient := rds.New(session.New(), &aws.Config{Region: aws.String(region)})

	for _,dbClusterIdentify := range dbClusterIdentifies{
		_, err := rdsClient.StopDBCluster(
			&rds.StopDBClusterInput{
				DBClusterIdentifier: aws.String(dbClusterIdentify.(string)),},)
		if err != nil {
			return err}
	}

	return nil
}

func startRDSClusters(dbClusterIdentifies []interface{},region string) (e error) {
	if len(dbClusterIdentifies) == 0 {
		return nil
	}

	rdsClient := rds.New(session.New(), &aws.Config{Region: aws.String(region)})

	for _,dbClusterIdentify := range dbClusterIdentifies {
		_, err := rdsClient.StartDBCluster(
			&rds.StartDBClusterInput{
				DBClusterIdentifier: aws.String(dbClusterIdentify.(string)),
			})
		if err != nil {
			return err
		}
	}

	return nil
}

func noticeToSlack(webHookURL,channel,action string) error {
	//どっちかが""なら通知しない
	if webHookURL=="" || channel=="" {
		return nil
	}

	var title= "*test/stage環境のサーバー再起動開始通知*"
	if action == "stop" {
		title = "*test/stage環境のサーバー停止開始通知*"
	}

	err := ReportToSlack(webHookURL, SlackRequestBody{
		IconEmoji: aws.String(SLACK_MONITORING_ICON),
		Username:  SLACK_MONITORING_NAME,
		Channel:   channel,
		Text:      title,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func extractTargetEc2InstanceIDs(ec *ec2.EC2,instanceNames []interface{})([]*string,error){
	targetInstances := []*string{}
	for _,instanceName := range instanceNames{
		targetInstances = append(targetInstances,aws.String(instanceName.(string)))
	}
	descOutput,err := ec.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters:[]*ec2.Filter{
			&ec2.Filter{
				Name:aws.String("tag:Name"),
				Values:targetInstances,
			},
		},
	})
	if err != nil {
		log.Error(err)
		return nil,err
	}
	var instanceIDs []*string
	for _,reservation := range descOutput.Reservations{
		instanceIDs = append(instanceIDs,reservation.Instances[0].InstanceId)
	}
	return instanceIDs,nil
}



