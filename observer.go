package slogt

import "golang.org/x/exp/slog"

// logs observer
// store captured logs
// implements the RecordCollector interface required by ObserverHandler
type Observer struct {
	records []slog.Record
}

// add a log in slog.Record format
func (o *Observer) addRecord(record slog.Record) {
	o.records = append(o.records, record)
}

// return the number of captured log(s)
func (o *Observer) CountLogs() int {
	return len(o.records)
}

// find a log by its chronological index
// if not found a zero-ed Log is returned
func (o *Observer) FindLog(index int) (log Log, found bool) {
	if index <= len(o.records) {
		found = true
		log = Log{record: o.records[index-1]}
	}
	return log, found
}
