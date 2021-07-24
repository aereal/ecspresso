//go:generate go run github.com/golang/mock/mockgen -package mockaws -destination ./mockaws/cwlogs_mock.go github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface CloudWatchLogsAPI

package ecspresso
