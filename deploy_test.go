package ecspresso_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"github.com/kayac/ecspresso"
	"github.com/kayac/ecspresso/mockaws"
)

type desiredCountTestCase struct {
	sv       *ecs.Service
	opt      ecspresso.DeployOption
	expected *int64
}

var desiredCountTestSuite = []desiredCountTestCase{
	{
		sv:       &ecs.Service{DesiredCount: nil},
		opt:      ecspresso.DeployOption{DesiredCount: nil},
		expected: nil,
	},
	{
		sv:       &ecs.Service{DesiredCount: nil, SchedulingStrategy: aws.String("DAEMON")},
		opt:      ecspresso.DeployOption{DesiredCount: aws.Int64(10)},
		expected: nil,
	},
	{
		sv:       &ecs.Service{DesiredCount: aws.Int64(2)},
		opt:      ecspresso.DeployOption{DesiredCount: nil},
		expected: nil,
	},
	{
		sv:       &ecs.Service{DesiredCount: aws.Int64(1)},
		opt:      ecspresso.DeployOption{DesiredCount: aws.Int64(3)},
		expected: aws.Int64(3),
	},
	{
		sv:       &ecs.Service{DesiredCount: aws.Int64(1)},
		opt:      ecspresso.DeployOption{DesiredCount: aws.Int64(ecspresso.DefaultDesiredCount)},
		expected: aws.Int64(1),
	},
	{
		sv:       &ecs.Service{DesiredCount: nil},
		opt:      ecspresso.DeployOption{DesiredCount: aws.Int64(5)},
		expected: aws.Int64(5),
	},
	{
		sv:       &ecs.Service{DesiredCount: nil},
		opt:      ecspresso.DeployOption{DesiredCount: aws.Int64(ecspresso.DefaultDesiredCount)},
		expected: nil,
	},
}

func TestCalcDesiredCount(t *testing.T) {
	for n, c := range desiredCountTestSuite {
		count := ecspresso.CalcDesiredCount(c.sv, c.opt)
		if count == nil && c.expected == nil {
			// ok
		} else if count != nil && c.expected == nil {
			t.Errorf("case %d unexpected desired count:%d expected:nil", n, *count)
		} else if count == nil && c.expected != nil {
			t.Errorf("case %d unexpected desired count:nil expected:%d", n, *c.expected)
		} else if *count != *c.expected {
			t.Errorf("case %d unexpected desired count:%d expected:%d", n, *count, *c.expected)
		} else {
			// ok
		}
	}
}

func TestDeploy_ecsDeploy(t *testing.T) {
	c := &ecspresso.Config{
		Region:             "ap-northeast-1",
		Timeout:            60 * time.Second,
		Service:            "test-service",
		Cluster:            "default",
		TaskDefinitionPath: "tests/td.json",
	}
	if err := c.Restrict(); err != nil {
		t.Error(err)
		return
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ecsClient := mockaws.NewMockECSAPI(ctrl)
	aasClient := mockaws.NewMockApplicationAutoScalingAPI(ctrl)
	cdClient := mockaws.NewMockCodeDeployAPI(ctrl)
	cwLogsClient := mockaws.NewMockCloudWatchLogsAPI(ctrl)
	iamClient := mockaws.NewMockIAMAPI(ctrl)
	app, err := ecspresso.NewAppWithAWSAggregate(c, ecspresso.AWSAggregate{ECS: ecsClient, ApplicationAutoScaling: aasClient, CodeDeploy: cdClient, CWLogs: cwLogsClient, IAM: iamClient})
	if err != nil {
		t.Error(err)
	}
	ecsClient.EXPECT().DescribeServicesWithContext(gomock.Any(), gomock.Any()).Times(1).Return(&ecs.DescribeServicesOutput{
		Services: []*ecs.Service{
			{
				ServiceName:    aws.String("test-service"),
				ClusterArn:     aws.String("arn:aws:ecs:ap-northeast-1:123456789012:cluster/default"),
				TaskDefinition: aws.String("arn:aws:ecs:ap-northeast-1:123456789012:taskDefinition/td1:1"),
				LoadBalancers: []*ecs.LoadBalancer{
					{
						ContainerName:    aws.String("app"),
						ContainerPort:    aws.Int64(80),
						LoadBalancerName: aws.String("alb-test"),
						TargetGroupArn:   aws.String("arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/test-tg-1/test-tg-1-id"),
					},
				},
			},
		},
	}, nil)
	ecsClient.EXPECT().ListTaskDefinitionsWithContext(gomock.Any(), gomock.Any()).Times(1).Return(&ecs.ListTaskDefinitionsOutput{
		TaskDefinitionArns: []*string{aws.String("arn:aws:ecs:ap-northeast-1:123456789012:taskDefinition/td1:1")},
	}, nil)
	ecsClient.EXPECT().UpdateServiceWithContext(gomock.Any(), gomock.Any()).Times(1).Return(&ecs.UpdateServiceOutput{}, nil)
	aasClient.EXPECT().DescribeScalableTargets(gomock.Any()).Times(1).Return(&applicationautoscaling.DescribeScalableTargetsOutput{ScalableTargets: []*applicationautoscaling.ScalableTarget{}}, nil)
	opt := ecspresso.DeployOption{
		DryRun:               aws.Bool(false),
		LatestTaskDefinition: aws.Bool(true),
		ForceNewDeployment:   aws.Bool(false),
		NoWait:               aws.Bool(true),
	}
	if err := app.Deploy(opt); err != nil {
		t.Error(err)
		return
	}
}
