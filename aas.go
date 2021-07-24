//go:generate go run github.com/golang/mock/mockgen -package mockaws -destination ./mockaws/aas_mock.go github.com/aws/aws-sdk-go/service/applicationautoscaling/applicationautoscalingiface ApplicationAutoScalingAPI

package ecspresso
