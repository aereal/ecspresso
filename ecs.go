//go:generate go run github.com/golang/mock/mockgen -package mockaws -destination ./mockaws/ecs_mock.go github.com/aws/aws-sdk-go/service/ecs/ecsiface ECSAPI

package ecspresso
