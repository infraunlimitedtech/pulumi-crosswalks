package infra

type Infra interface {
	GetNodes() map[string]map[string]interface{}
}
