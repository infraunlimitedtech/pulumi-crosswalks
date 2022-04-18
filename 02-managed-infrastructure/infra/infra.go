package infra

type ComputeInfra interface {
	GetNodes() map[string]map[string]interface{}
}

type S3Infra interface{
	GetStorage() map[string]map[string]interface{}
}
