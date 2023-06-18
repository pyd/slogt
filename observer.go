package slogt

import "golang.org/x/exp/slog"

// logs observer
// store captured logs
// implements the RecordCollector interface required by ObserverHandler
type Observer struct {
	records []slog.Record
}

func (o *Observer) addRecord(record slog.Record) {
	o.records = append(o.records, record)
}
