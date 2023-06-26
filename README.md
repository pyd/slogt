
# slogt

Test logs from golang slog logger.

SLOG DESIGN SOURCE: https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md

WARNING: slogt will fail finding an attribute if its Key has a dot e.g. "user.profil"
as it will consider it a 2 subkeys key !!!!
A solution could be to provide a SetAttributeKeySeparator(separator string)

// TODO add a method to clear logs in observer

// TODO find a log by its - chronological - index may be difficult?
 
TODO: it may be useful to know if a group attribute has the expected keys ?

e.g. log.HasAttributesInGroup(groupAttributeKey, keys ...string) found bool, error
e.g. log.HasAttributesInGroup("user.profile", "age", admin") found bool, error
errors: groupAttribute not found
        subAttributes ["age"] not found

TODO: check file, line when AddSOurce is set in config



### Attributes

Built-in attributes are defined in Logger methods e.g. Warn() Log()...
In output they are prefixed with handler groups joined by a dot.

Shared attributes are defined in the handler.
Logger attributes (e.g. Logger.With() will create a new handler with additional attributes)


// TODO check how Logger.Log() defines attributes (builtin or shared)? 
- func (l *Logger) Log(ctx context.Context, level Level, msg string, args ...any)
    Log emits a log record with the current time and the given level and
    message. The Record's Attrs consist of the Logger's attributes followed by
    the Attrs specified by args.


// TODO check how AddAttrs can affect slogt. It could be used in a custom handler to add attributes
// after the slog.Record has been passed to the Observer!
func (r *Record) AddAttrs(attrs ...Attr) {
	n := copy(r.front[r.nFront:], attrs)
	r.nFront += n
	// Check if a copy was modified by slicing past the end
	// and seeing if the Attr there is non-zero.
	if cap(r.back) > len(r.back) {
		end := r.back[:len(r.back)+1][len(r.back)]
		if !end.isEmpty() {
			panic("copies of a slog.Record were both modified")
		}
	}
	r.back = append(r.back, attrs[n:]...)
}

### Groups

func (l *Logger) WithGroup(name string) *Logger {
	c := l.clone()
	c.handler = l.handler.WithGroup(name)
	return c

}

func (h *JSONHandler) WithGroup(name string) Handler {
	return &JSONHandler{commonHandler: h.commonHandler.withGroup(name)}
}

func (h *TextHandler) WithGroup(name string) Handler {
	return &TextHandler{commonHandler: h.commonHandler.withGroup(name)}
}

### Source

// TODO add method to check source in Log?

type Source struct {
	// Function is the package path-qualified function name containing the
	// source line. If non-empty, this string uniquely identifies a single
	// function in the program. This may be the empty string if not known.
	Function string `json:"function"`
	// File and Line are the file name and line number (1-based) of the source
	// line. These may be the empty string and zero, respectively, if not known.
	File string `json:"file"`
	Line int    `json:"line"`
}

### CONTEXT

// TODO check context?