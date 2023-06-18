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