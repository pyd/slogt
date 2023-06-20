# Slogt.

Helper for checking logs from slog in tests.

## Logs capture.

A log can be captured when it is written (with a custom io.Writer passed to the handler)
or when iut is passed - as a slog.record - to the handler by the logger (via handler.Handle(ctx context.Context, record slog.Record)).

In the first case we will have to parse the log string to extract time, message...
The slog.record already implements an api to get time, message, level and attributes of the log. 

Let's try the handler way.

## Handler

First let's create a custom handler which implements the slog.Handler interface.
Handler will embed an observer and pass it each slog.Record from slog.Logger.
Observer must implement a RecordCollector interface to receive those records.

## Observer

Must be a pointer (adding records to it).
It's API should provide the number of captured logs and a getter for a log
by its index in the collection (chronological order of logs capture)

## Log

Represents a log
The log.go.archive file has a different implementation of the FindAttribute() method.
Less readable than the final version, I keep it for the recursive findAttributeInGroup() method
that may be useful for future evolutions.
 
TODO: it may be useful to know if a group attribute has the expected keys ?

e.g. log.HasAttributesInGroup(groupAttributeKey, keys ...string) found bool, error
e.g. log.HasAttributesInGroup("user.profile", "age", admin") found bool, error
errors: groupAttribute not found
        subAttributes ["age"] not found

TODO: remove methods implementations in ObserverHandler