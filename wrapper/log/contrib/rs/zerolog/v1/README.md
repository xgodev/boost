zerolog.v1
=======

Example
--------

```go
package main

import (
	"context"

	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/log/contrib/rs/zerolog.v1"
)

func main() {
	ctx := context.Background()

	//example use zerolog
	logger := zerolog.NewLogger()

	logger = logger.WithField("main_field", "example")

	logger.Info("main method.")
	//output: 2:30PM INF main method. main_field=example

	ctx = logger.ToContext(ctx)

	foo(ctx)

	withoutContext()
}

func foo(ctx context.Context) {
	logger := log.FromContext(ctx)

	logger = logger.WithField("foo_field", "example")
	logger.Infof("%s method.", "foo")
	//output: 2:30PM INF foo method. foo_field=example main_field=example

	ctx = logger.ToContext(ctx)
	bar(ctx)
}

func bar(ctx context.Context) {
	logger := log.FromContext(ctx)

	logger = logger.WithField("bar_field", "example")

	logger.Infof("%s method.", "bar")
	//output: 2:30PM INF bar method. bar_field=example foo_field=example main_field=example
}

func withoutContext() {
	log.Info("withoutContext method")
	//output: 2:30PM INF withoutContext method
}
```

default options:

| option  | value  |
|---|---|
| Formatter  | "TEXT"  |
| Level  | "INFO"  |
| ConsoleEnabled  | true  |
| FileEnabled  | false  |
| FilePath  | "/tmp"  |
| FileName  | "application.log"  |
| FileMaxSize  | 100  |
| FileCompress  | true  |
| FileMaxAge  | 28  |
| ErrorFieldName | "err" | 

The package accepts a default constructor:
```go
logger := zerolog.NewLogger()
```
Or a constructor with Options:
```go
logger := zerolog.NewLoggerWithOptions(&zerolog.Options{})
```
Or a constructor with multiple parameters using options functions pattern:
```go
logger := zerolog.NewLogger(
	zerolog.WithFormatter("TEXT"),
	zerolog.WithConsoleEnabled(true),
	zerolog.WithFilePath("/tmp"),
	...
)
```

This is the list of all the configuration functions supported by package:

#### WithFormatter
sets output format of the logs. Using TEXT/JSON.
```go
// text formatter
logger := zerolog.NewLogger(zerolog.WithFormatter("TEXT"))

// json formatter
logger := zerolog.NewLogger(zerolog.WithFormatter("JSON"))
```

#### WithLevel
sets logging level to any of these options below on the standard logger.
```go
// log level DEBUG
logger := zerolog.NewLogger(zerolog.WithLevel("DEBUG"))

// log level WARN
logger := zerolog.NewLogger(zerolog.WithLevel("WARN"))

// log level FATAL
logger := zerolog.NewLogger(zerolog.WithLevel("FATAL"))

// log level ERROR
logger := zerolog.NewLogger(zerolog.WithLevel("ERROR"))

// log level TRACE
logger := zerolog.NewLogger(zerolog.WithLevel("TRACE"))

// log level INFO
logger := zerolog.NewLogger(zerolog.WithLevel("INFO"))
```

#### WithConsoleEnabled
sets whether the standard logger output will be in console. Accepts multi writing (console and file).
##### Enabled
```go
// console enable true
logger := zerolog.NewLogger(zerolog.WithConsoleEnabled(true))

// console enable false
logger := zerolog.NewLogger(zerolog.WithConsoleEnabled(false))
```

#### WithFileEnabled
sets whether the standard logger output will be in file. Accepts multi writing (file and console).
##### Enabled
```go
// file enable true
logger := zerolog.NewLogger(zerolog.WithFileEnabled(true))

// file enable false
logger := zerolog.NewLogger(zerolog.WithFileEnabled(false))
```

##### WithFilePath
sets the path where the file will be saved.
```go
// file path
logger := zerolog.NewLogger(zerolog.WithFilePath("/tmp"))
```

##### WithFileName
sets the name of the file.
```go
// file name
logger := zerolog.NewLogger(zerolog.WithFileName("application.log"))
```

##### WithFileMaxSize
sets the maximum size in megabytes of the log file. It defaults to 100 megabytes.
```go
// file max size
logger := zerolog.NewLogger(zerolog.WithFileMaxSize(100))
```

##### WithFileCompress
sets whether the log files should be compressed.
```go
// file compress true
logger := zerolog.NewLogger(zerolog.WithFileCompress(true))

// file compress false
logger := zerolog.NewLogger(zerolog.WithFileCompress(false))
```

##### WithFileMaxAge
sets the maximum number of days to retain old log files based on the timestamp encoded in their filename.
```go
// file max age
logger := zerolog.NewLogger(zerolog.WithFileMaxAge(10))
```

##### WithErrorFieldName
sets the field name used on `WithError`
```go
logger := zerolog.NewLogger(zerolog.WithErrorFieldName("error"))
```