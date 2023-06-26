# slogt

Test logs from golang slog logger.

SLOG DESIGN SOURCE: https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md


## Observer
// TODO add a method to find a log whicn has an attribute (by its Key)

// TODO move handler arg to second pos in NewObserverHandler


// TODO find a log by its - chronological - index may be difficult?
 
TODO: it may be useful to know if a group attribute has the expected keys ?

### Attributes

Built-in attributes are defined in Logger methods e.g. Warn() Log()...
In output they are prefixed with handler groups joined by a dot.

Shared attributes are defined in the handler.
Logger attributes (e.g. Logger.With() will create a new handler with additional attributes)

### CONTEXT

// TODO check context?