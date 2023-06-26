
# slogt

Test logs from golang slog logger.

WARNING: slogt will fail finding an attribute if its Key has a dot e.g. "user.profil"
as it will consider it a 2 subkeys key !!!!
A solution could be to provide a SetAttributeKeySeparator(separator string)

// TODO refactor: add FindSharedAttribute() method in Handler

// TODO add a method to clear logs in observcer

## usage:

Instanciate the observer.  
It will receive and store logs from the handler (you need a pointer) 

    observer := new(slogt.Observer)
    observer = &slogt.Observer{}

Instantiate the handler which embed the observer

    // an error is returned if the observer arg is nil
    handler, _ := slogt.NewDefaultObserverHandler(observer)

Instantiate the logger

Wahoou, l'instantiation du logger ct pour les tests de slogt lui meme
Avec zap on passait l'instance du logger comme une dépendance, ici on a un global Logger :/ :/ :/ PFFFFF
###############################################################################################################################  
POUR TESTER UN AUTRE CODE IL FAUT FAIRE UN slog.SetDefault() !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! ???????????????????????????????????????
WTF QUE CE PASSE T'IL S'IL Y A UN DEJA UN SetDefault ??? EXISTE t IL UN slog.GetDefault() pour récupérer l'instance/la config globale ???
###############################################################################################################################

    // Alternatively, if you want to pass your own embedded handler
// var handler slogt.ObserverHandler
// var embeddedHandler = slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})
// var handler, _ = slogt.NewObserverHandler(embeddedHandler, observer)
var logger = slog.New(handler)

// Now you are ready to go
// in tested code
func TestedCode() {
	// ...
	slog.Warn("warn message")
	// ...
	slog.Info("info message", )
	// ...
	slog.Error("user can't order ...", slog.Group("user", slog.Group("profile", slog.Int("age", user.Age))))
}

Expect(observer.CountLogs()).To(Equal(3))

var log = observer.FindLog(3)

Expect(log.Message()).To(Equal("error message"))
Expect(log.Level()).To(Equal(slog.LevelInfo))
Expect(log.Time()).To(BeTemporally("~", time.Now(), time.Millisecond*500))

attribute, attributeFound = log.FindBuiltInAttribute("user.profile.age")

Expect(attributeFound).To(BeTrue())
Expect(attribute.Value.Int64()).To(Equal(int64(user.Age)))


le Handler doit contenir groups

pour pouvoir les passer à chaque appel de   addRecord(record, groups)


Log doit maintenant contenir:
record slog.Record
groups Handler.groups []string
 
TODO: it may be useful to know if a group attribute has the expected keys ?

e.g. log.HasAttributesInGroup(groupAttributeKey, keys ...string) found bool, error
e.g. log.HasAttributesInGroup("user.profile", "age", admin") found bool, error
errors: groupAttribute not found
        subAttributes ["age"] not found

TODO: remove methods implementations in ObserverHandler

TODO: enhance ObserverHandler.NewDefaultObserverHandler: explain use of TextHandler with io.Discard Writer*

TODO: RecordCollector addRecord() must be exported
    : also rename this interface LoggerObserver because we also check for group with WithGroup

TODO: check file, line when AddSOurce is set in config

SLOG DESIGN SOURCE: https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md

### Attributes

built-in attributs du record (qualifiés avec les groups)

attributs logger et handler (non qualifiés)

// quel peut être la source d'un attribut:
- les arguments d'une méthode de log e.g. slog.Info("msg", args)
- le logger: func (l *Logger) With(args ...any) *Logger
    la bonne nouvelle c'est que ces attrs sont passés au handler via l'appel à handler.WithAttrs (ON LE GERE)
        The new Logger's handler is the result of calling WithAttrs on the receiver's handler.


- func (l *Logger) Log(ctx context.Context, level Level, msg string, args ...any)
    Log emits a log record with the current time and the given level and
    message. The Record's Attrs consist of the Logger's attributes followed by
    the Attrs specified by args.

// With calls Logger.With on the default logger.
func With(args ...any) *Logger {
	return Default().With(args...)
}

// LogAttrs is a more efficient version of [Logger.Log] that accepts only Attrs.
func (l *Logger) LogAttrs(ctx context.Context, level Level, msg string, attrs ...Attr) {
	l.logAttrs(ctx, level, msg, attrs...)
}

// With returns a new Logger that includes the given arguments, converted to
// Attrs as in [Logger.Log].
// The Attrs will be added to each output from the Logger.
// The new Logger shares the old Logger's context.
// The new Logger's handler is the result of calling WithAttrs on the receiver's
// handler.
func (l *Logger) With(args ...any) *Logger {
	c := l.clone()
	c.handler = l.handler.WithAttrs(argsToAttrSlice(args))
	return c
}

************************************************************************************************
ATTENTION ICI, ON POURRAIT IMAGINER UN CUSTOM HANDLER QUI AJOUTE DES ATTRIBUTS À UN RECORD !!!
apres qu'on l'ai transmi à l'observer: workaround, lui passer un pointer du record 

// AddAttrs appends the given Attrs to the Record's list of Attrs.
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
************************************************************************************************



### GROUP

LogValuer interface n'aura pas d'incidence sur les group des attributs
    LogValue() slog.Value 


// WithGroup returns a new Logger that starts a group. The keys of all
// attributes added to the Logger will be qualified by the given name.
// (How that qualification happens depends on the [Handler.WithGroup]
// method of the Logger's Handler.)
// The new Logger shares the old Logger's context.
//
// The new Logger's handler is the result of calling WithGroup on the receiver's
// handler.
func (l *Logger) WithGroup(name string) *Logger {
	c := l.clone()
	c.handler = l.handler.WithGroup(name)
	return c

}

### Source

On pourrait passer une custom func pour récupérer file, line

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

### LOGGER

logger := slog.With(slog.Int("userId", 99))

### CONTEXT

### LOG

FindAttribute(key) returns built-in or handler attribute:

	- if log has group names & key is qualified it will search in built-in attributes
	- if log has group names & key is not qualified it will search in handler attributes






