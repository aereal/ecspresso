//go:generate go run github.com/golang/mock/mockgen -package mockaws -destination ./mockaws/iam_mock.go github.com/aws/aws-sdk-go/service/iam/iamiface IAMAPI

package ecspresso
