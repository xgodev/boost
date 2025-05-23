# Cobra Integration for Boost

## Overview

The Cobra integration in Boost provides a streamlined wrapper around the popular [Cobra](https://github.com/spf13/cobra) command-line interface framework. This package enables developers to quickly build powerful CLI applications while seamlessly integrating with Boost's configuration system.

## Features

- **Command Management**: Easy creation and organization of CLI commands
- **Automatic Flag Generation**: Automatic generation of flags from Boost configuration entries
- **Type Support**: Support for various data types in command flags
- **Configuration Integration**: Seamless integration with Boost's configuration system
- **Command Composition**: Simple composition of command hierarchies

## Installation

```go
import (
    "github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
)
```

## Basic Usage

### Creating a Simple Command

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
    boostCobra "github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Create root command
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "A brief description of your application",
        Long:  `A longer description of your application`,
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello from myapp!")
        },
    }
    
    // Run the command
    if err := boostCobra.Run(rootCmd); err != nil {
        log.Fatal(err)
    }
}
```

### Creating Command Hierarchies

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
    boostCobra "github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Create root command
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "A brief description of your application",
        Long:  `A longer description of your application`,
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello from myapp!")
        },
    }
    
    // Create sub-commands
    versionCmd := &cobra.Command{
        Use:   "version",
        Short: "Print the version number",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("myapp v1.0")
        },
    }
    
    serveCmd := &cobra.Command{
        Use:   "serve",
        Short: "Start the server",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Starting server...")
        },
    }
    
    // Run the command with sub-commands
    if err := boostCobra.Run(rootCmd, versionCmd, serveCmd); err != nil {
        log.Fatal(err)
    }
}
```

### Using NewCommand Helper

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
    boostCobra "github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Create sub-commands
    versionCmd := &cobra.Command{
        Use:   "version",
        Short: "Print the version number",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("myapp v1.0")
        },
    }
    
    serveCmd := &cobra.Command{
        Use:   "serve",
        Short: "Start the server",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Starting server...")
        },
    }
    
    // Create root command with sub-commands
    rootCmd := boostCobra.NewCommand(
        &cobra.Command{
            Use:   "myapp",
            Short: "A brief description of your application",
            Long:  `A longer description of your application`,
            Run: func(cmd *cobra.Command, args []string) {
                fmt.Println("Hello from myapp!")
            },
        },
        versionCmd,
        serveCmd,
    )
    
    // Run the command
    if err := boostCobra.Run(rootCmd); err != nil {
        log.Fatal(err)
    }
}
```

## Integration with Boost Configuration

The Cobra integration automatically generates command-line flags from Boost configuration entries:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
    boostCobra "github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
    "github.com/xgodev/boost/wrapper/config"
    "github.com/xgodev/boost/wrapper/log"
)

func init() {
    // Define configuration entries
    config.Add("app.server.port", 8080, "Server port")
    config.Add("app.server.host", "localhost", "Server host")
    config.Add("app.log.level", "INFO", "Log level")
}

func main() {
    // Create root command
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "A brief description of your application",
        Run: func(cmd *cobra.Command, args []string) {
            // Access configuration values
            port := config.Int("app.server.port")
            host := config.String("app.server.host")
            logLevel := config.String("app.log.level")
            
            fmt.Printf("Server: %s:%d, Log Level: %s\n", host, port, logLevel)
        },
    }
    
    // Run the command
    if err := boostCobra.Run(rootCmd); err != nil {
        log.Fatal(err)
    }
}
```

When running the application, the flags will be automatically available:

```
$ ./myapp --app.server.port=9090 --app.server.host=0.0.0.0 --app.log.level=DEBUG
Server: 0.0.0.0:9090, Log Level: DEBUG
```

## Supported Data Types

The Cobra integration supports various data types for flags:

- `string` and `[]string`
- `bool` and `[]bool`
- `int`, `[]int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint64`
- `float64`
- `time.Duration`
- `[]byte`
- `net.IP`, `net.IPNet`, `net.IPMask`
- `map[string]string`

## Best Practices

1. **Organize Commands Hierarchically**: Group related commands together
   ```go
   rootCmd := boostCobra.NewCommand(appCmd, userCmd, configCmd)
   userCmd := boostCobra.NewCommand(userListCmd, userCreateCmd, userDeleteCmd)
   ```

2. **Use Descriptive Help Text**: Provide clear descriptions for commands and flags
   ```go
   cmd := &cobra.Command{
       Use:   "serve",
       Short: "Start the API server",
       Long:  `Start the API server with the specified configuration.
               The server will listen on the configured port and host.`,
   }
   ```

3. **Leverage Boost Configuration**: Use Boost's configuration system for default values
   ```go
   config.Add("app.timeout", 30, "Connection timeout in seconds")
   ```

4. **Handle Errors Properly**: Always check for errors when running commands
   ```go
   if err := boostCobra.Run(rootCmd); err != nil {
       log.Fatalf("Error: %v", err)
   }
   ```

5. **Use Command Validation**: Implement PreRun or PreRunE for validation
   ```go
   cmd := &cobra.Command{
       PreRunE: func(cmd *cobra.Command, args []string) error {
           // Validate arguments or configuration
           return nil
       },
   }
   ```

## Contributing

Contributions to improve the Cobra integration are welcome. Please follow the Boost project's contribution guidelines.

## License

This package is part of the Boost project and is subject to its license terms.
