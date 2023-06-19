package slogt

import "golang.org/x/exp/slog"

// a log
type Log struct {
	// record is the data source
	record slog.Record
}

// get the message of the log
func (l Log) Message() string {
	return l.record.Message
}
