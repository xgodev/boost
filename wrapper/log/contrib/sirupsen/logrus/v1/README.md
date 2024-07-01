logrus.v1
=======

Example
--------

```go
package main

import (
	"context"

	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus.v1"
)

func main() {
	ctx := context.Background()

	//example use logrus
	log.Set(logrus.NewLogger())

	logger := log.WithField("main_field", "example")

	logger.Info("main method.")
	//output: INFO[2021/05/14 17:15:04.757] main method. main_field=example

	ctx = logger.ToContext(ctx)

	foo(ctx)

	withoutContext()
}

func foo(ctx context.Context) {
	logger := log.FromContext(ctx)

	logger = logger.WithField("foo_field", "example")
	logger.Infof("%s method.", "foo")
	//output: INFO[2021/05/14 17:15:04.757] foo method. foo_field=example main_field=example

	ctx = logger.ToContext(ctx)
	bar(ctx)
}

func bar(ctx context.Context) {
	logger := log.FromContext(ctx)

	logger = logger.WithField("bar_field", "example")

	logger.Infof("%s method.", "bar")
	//output: INFO[2021/05/14 17:15:04.757] bar method. bar_field=example foo_field=example main_field=example
}

func withoutContext() {
	log.Info("withoutContext method")
	//output: INFO[2021/05/14 17:15:04.757] withoutContext method
}
```

default options:

| option  | value  |
|---|---|
| Formatter | text.New() |
| ConsoleEnabled | true |
| ConsoleLevel | "INFO" |
| FileEnabled | false |
| FileLevel | "INFO" |
| FilePath | "/tmp" |
| FileName | "application.log" |
| FileMaxSize | 100 |
| FileCompress | true |
| FileMaxAge | 28 |
| TimeFormat | "2006/01/02 15:04:05.000" |
| ErrorFieldName | "err" | 

The package accepts a default constructor:
```go
// default constructor
logger := logrus.NewLogger()
```
Or a constructor with Options:
```go
logger := logrus.NewLoggerWithOptions(&logrus.Options{})
```
Or a constructor with multiple parameters using optional pattern:
```go
// multiple optional parameters constructor
logger := logrus.NewLogger(
	logrus.WithFormatter(json.New())
	logrus.WithConsoleEnabled(true),
	logrus.WithFilePath("/tmp"),
	...
)
```

This is the list of all the configuration functions supported by package:

#### WithFormatter
sets output format of the logs. Using TEXT/JSON/CLOUDWATCH.
```go
import (
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus.v1/formatter/text"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus.v1/formatter/json"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus.v1/formatter/cloudwatch"
	...
)

// text formatter
logger := logrus.NewLogger(logrus.WithFormatter(text.New()))

// json formatter
logger := logrus.NewLogger(logrus.WithFormatter(json.New()))

// cloudwatch formatter
logger := logrus.NewLogger(logrus.WithFormatter(cloudwatch.New()))
```

#### WithTimeFormat
sets the format used for marshaling timestamps.
```go
// time format
logger := logrus.NewLogger(logrus.WithTimeFormat("2006/01/02 15:04:05.000"))
```

#### WithConsoleEnabled
sets whether the standard logger output will be in console. Accepts multi writing (console and file).
```go
// console enable true
logger := logrus.NewLogger(logrus.WithConsoleEnabled(true))

// console enable false
logger := logrus.NewLogger(logrus.WithConsoleEnabled(false))
```

#### WithConsoleLevel
sets console logging level to any of these options below on the standard logger.
```go
// log level DEBUG
logger := logrus.NewLogger(logrus.WithConsoleLevel("DEBUG"))

// log level WARN
logger := logrus.NewLogger(logrus.WithConsoleLevel("WARN"))

// log level FATAL
logger := logrus.NewLogger(logrus.WithConsoleLevel("FATAL"))

// log level ERROR
logger := logrus.NewLogger(logrus.WithConsoleLevel("ERROR"))

// log level TRACE
logger := logrus.NewLogger(logrus.WithConsoleLevel("TRACE"))

// log level INFO
logger := logrus.NewLogger(logrus.WithConsoleLevel("INFO"))
```

#### WithHook
sets a hook to be fired when logging on the logging levels.
```go
import (
	"log/syslog"

	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus.v1"
	syshooklg "github.com/sirupsen/logrus/hooks/syslog"
)

func main() {
	hook, _ := syshooklg.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")
	logger := logrus.NewLogger(logrus.WithHook(hook))
	...
}
```

#### WithFileEnabled
sets whether the standard logger output will be in file. Accepts multi writing (file and console).
```go
// file enable true
logger := logrus.NewLogger(logrus.WithFileEnabled(true))

// file enable false
logger := logrus.NewLogger(logrus.WithFileEnabled(false))
```

#### WithFilePath
sets the path where the file will be saved.
```go
// file path
logger := logrus.NewLogger(logrus.WithFilePath("/tmp"))
```

#### WithFileName
sets the name of the file.
```go
// file name
logger := logrus.NewLogger(logrus.WithFileName("application.log"))
```

#### WithFileMaxSize
sets the maximum size in megabytes of the log file. It defaults to 100 megabytes.
```go
// file max size
logger := logrus.NewLogger(logrus.WithFileMaxSize(100))
```

#### WithFileCompress
sets whether the log files should be compressed.
```go
// file compress true
logger := logrus.NewLogger(logrus.WithFileCompress(true))

// file compress false
logger := logrus.NewLogger(logrus.WithFileCompress(false))
```

#### WithFileMaxAge
sets the maximum number of days to retain old log files based on the timestamp encoded in their filename.
```go
// file max age
logger := logrus.NewLogger(logrus.WithFileMaxAge(10))
```

##### WithErrorFieldName
sets the field name used on `WithError`
```go
logger := logrus.NewLogger(logrus.WithErrorFieldName("error"))
```
