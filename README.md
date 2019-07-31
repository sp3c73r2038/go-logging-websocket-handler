inspired by https://github.com/hhkbp2/go-logging

```go
handler := NewWebsocketHandler(
	"My",  // Name
    "127.0.0.1:5577", // host
	"/logs/channel/attendance",  // path
	logging.LevelTrace,  // level
)

format := "%(asctime)s %(levelname)s (%(filename)s:%(lineno)d) " +
	"%(name)s %(message)s"
dateFormat := "%Y-%m-%d %H:%M:%S.%3n"
// create a formatter(which controls how log messages are formatted)
formatter := logging.NewStandardFormatter(format, dateFormat)
// set formatter for handler
handler.SetFormatter(formatter)

stderr, err := logging.NewFileHandler("/dev/stderr", os.O_APPEND, -1)
if err != nil {
	log.Fatal(err)
}
stderr.SetFormatter(formatter)
stderr.SetLevel(logging.LevelTrace)

// create a logger(which represents a log message source)
logger := logging.GetLogger("a.b.c")
logger.SetLevel(logging.LevelTrace)
logger.AddHandler(handler)
logger.AddHandler(stderr)

// ensure all log messages are flushed to disk before program exits.
defer logging.Shutdown()
logger.Tracef("message: %s %d", "Hello", 2015)
logger.Debugf("message: %s %d", "Hello", 2015)
logger.Infof("message: %s %d", "Hello", 2015)
logger.Warnf("message: %s %d", "Hello", 2015)
logger.Errorf("message: %s %d", "Hello", 2015)
logger.Fatalf("message: %s %d", "Hello", 2015)
```
