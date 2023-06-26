package slogt

import "golang.org/x/exp/slog"

/*
The Observer stores captured logs
*/
type Observer struct {
	logs []Log
}

// add a log
func (o *Observer) addLog(record slog.Record, groups []string, attrs []slog.Attr) {
	log := NewLog(record, groups, attrs)
	o.logs = append(o.logs, log)
}

// return the number of captured log(s)
func (o *Observer) CountLogs() int {
	return len(o.logs)
}

// find a log by its chronological index
// if not found a zero-ed Log is returned
func (o *Observer) FindLog(index int) (log Log, found bool) {
	if index <= len(o.logs) {
		found = true
		log = o.logs[index-1]
	}
	return log, found
}
