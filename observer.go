package slogt

/*
The Observer stores captured logs
*/
type Observer struct {
	logs []Log
}

// Add a log.
func (o *Observer) addLog(log Log) {
	o.logs = append(o.logs, log)
}

// Count the number of captured log(s).
func (o *Observer) CountLogs() int {
	return len(o.logs)
}

// Find a log by its chronological index.
// If not found a zero-ed Log is returned.
func (o *Observer) FindLog(index int) (log Log, found bool) {
	if index <= len(o.logs) {
		found = true
		log = o.logs[index-1]
	}
	return log, found
}

// Remove all captured logs.
func (o *Observer) ClearLogs() {
	o.logs = []Log{}
}

// Get all captured logs.
func (o *Observer) Logs() []Log {
	return o.logs
}
