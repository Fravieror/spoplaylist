package source_music

import "github.com/newrelic/go-agent/v3/newrelic"

type ISourceMusic interface {
	GetHot100Songs(txn *newrelic.Transaction, date string) ([]string, error)
}

type MockSourceMusic struct {
	MockGetHot100Songs func(txn *newrelic.Transaction, date string) ([]string, error)
}

func (m MockSourceMusic) GetHot100Songs(txn *newrelic.Transaction, date string) ([]string, error) {
	return m.MockGetHot100Songs(txn, date)
}
