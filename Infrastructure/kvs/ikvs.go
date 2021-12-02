package ds

type IKVS interface {
	Get() interface{}
	Set() error
}
