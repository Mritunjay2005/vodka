# Vodka

```text
__/\\\________/\\\_______/\\\\\_______/\\\\\\\\\\\\_____/\\\________/\\\_____/\\\\\\\\\____        
 _\/\\\_______\/\\\_____/\\\///\\\____\/\\\////////\\\__\/\\\_____/\\\//____/\\\\\\\\\\\\\__       
  _\//\\\______/\\\____/\\\/__\///\\\__\/\\\______\//\\\_\/\\\__/\\\//______/\\\/////////\\\_      
   __\//\\\____/\\\____/\\\______\//\\\_\/\\\_______\/\\\_\/\\\\\\//\\\_____\/\\\_______\/\\\_     
    ___\//\\\__/\\\____\/\\\_______\/\\\_\/\\\_______\/\\\_\/\\\//_\//\\\____\/\\\\\\\\\\\\\\\_    
     ____\//\\\/\\\_____\//\\\______/\\\__\/\\\_______\/\\\_\/\\\____\//\\\___\/\\\/////////\\\_   
      _____\//\\\\\_______\///\\\__/\\\____\/\\\_______/\\\__\/\\\_____\//\\\__\/\\\_______\/\\\_  
       ______\//\\\__________\///\\\\\/_____\/\\\\\\\\\\\\/___\/\\\______\//\\\_\/\\\_______\/\\\_ 
        _______\///_____________\/////_______\////////////_____\///________\///__\///________\///__
```

**The ultra-fast, no-botanical-nonsense HTTP framework for Go.**

Vodka is a lightweight, high-performance web framework built to make Go HTTP servers ergonomic and incredibly fast. It strips away the bloat and provides exactly what you need: a blazing-fast radix tree router, a clean context wrapper, seamless middleware chaining, and a built-in hot-reloading development environment.

## Installation

Vodka comes in two parts: the framework itself, and a powerful CLI tool for local development.

1. Install the Framework

Run this inside your Go project to add Vodka to your go.mod:

```bash
go get github.com/DevanshuTripathi/vodka
```

2. Install the Hot-Reload CLI

The Vodka CLI acts as a supervisor process, watching your .go files and automatically rebuilding and restarting your server whenever you save.

```bash
go install github.com/DevanshuTripathi/vodka/cmd/vodka@latest
```

Important Note: Ensure your Go bin directory is added to your system's PATH, otherwise your terminal won't recognize the vodka command.
Mac/Linux: Add export PATH=$PATH:$(go env GOPATH)/bin to your ~/.bashrc or ~/.zshrc.
Windows: Ensure %USERPROFILE%\go\bin is added to your Environment Variables.

## Quickstart
Here is a minimal, fully functional Vodka application:

```Go
package main

import (
	"log"
	"[github.com/DevanshuTripathi/vodka](https://github.com/DevanshuTripathi/vodka)"
)

func main() {
	// Initialize a new Vodka engine
	app := vodka.New()

	// Define a simple GET route
	app.GET("/ping", func(c *vodka.Context) {
		c.String(200, "pong! 🏓")
	})

	// Start the server on port 8080
	if err := app.Run(":8080"); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}
```

To run this in development with hot-reloading, simply open your terminal in the project directory and type:

```bash
vodka
```

## What Everything Does (Architecture)

### The Engine (Router)

At its core, vodka.Engine is a highly optimized router powered by httprouter under the hood. Instead of scanning linearly through routes, it uses a Radix Tree (Prefix Tree). This means the time it takes to find a route is proportional only to the length of the URL path, not the number of routes in your API. It handles highly complex routing logic in under a microsecond without allocating new memory.

### The Context

vodka.Context is the central struct passed to all your handlers. It wraps Go's standard http.Request and http.ResponseWriter into a single, clean object and provides quality-of-life helper methods.

1. c.JSON(statusCode, obj): Automatically sets the content type and encodes your Go structs into JSON.

2. c.String(statusCode, text): Sends plain text responses instantly.

### Mixers (Middleware)

Vodka uses the Chain of Responsibility pattern to handle middleware sequentially. Middlewares are essentially an array of handlers.

You can inject global middleware using app.Use(). Inside your middleware, calling c.Next() pauses the current function, executes the next handler in the chain, and then returns control back to the middleware.

**Example: A Custom Logger Middleware**

```Go
func Logger() vodka.HandlerFunc {
	return func(c *vodka.Context) {
		// 1. Mark the start time (Before Request)
		t := time.Now()

		// 2. Pass control to the next handler
		c.Next()

		// 3. Calculate latency (After Request)
		latency := time.Since(t)
		log.Printf("[%s] %s %v\n", c.Request.Method, c.Request.URL.Path, latency)
	}
}

// Applying it globally:
app.Use(Logger())
```

### The CLI(vodka dev)

Go is a compiled language, meaning you can't normally hot-reload code in memory. The vodka CLI solves this by tapping into the OS file system (fsnotify).
When you run vodka, it:
1. Compiles your app into a temporary binary (tmp/vodka-build).

2. Starts the server.

3. Listens for changes to any .go files.

4. If a change is detected, it gracefully kills the process, rebuilds the binary, and starts it back up in milliseconds.