//go:generate go run github.com/golang/mock/mockgen -package mockaws -destination ./mockaws/codedeploy_mock.go github.com/aws/aws-sdk-go/service/codedeploy/codedeployiface CodeDeployAPI

package ecspresso
