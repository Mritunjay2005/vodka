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

Vodka is a lightweight, high-performance web framework built to make Go HTTP servers ergonomic and incredibly fast. It strips away the bloat and provides exactly what you need: a blazing-fast radix tree router, a clean context wrapper, seamless middleware chaining, and a **powerful built-in CLI for instantly generating and running full-stack React applications.**

## Installation: The Vodka CLI

To get the full Vodka experience, you first need to install the CLI tool globally. This tool acts as your project generator and dev-server supervisor.

```bash
go install https://github.com/DevanshuTripathi/vodka/cmd/vodka@latest
```

**Important Note:** Ensure your Go bin directory is added to your system's PATH, otherwise your terminal won't recognize the vodka command.
- **Mac/Linux:** Add export PATH=$PATH:$(go env GOPATH)/bin to your ~/.bashrc or ~/.zshrc.
- **Windows:** Ensure %USERPROFILE%\go\bin is added to your Environment Variables.

## Building Full-Stack Apps

Vodka is designed to bridge the gap between Go backends and modern single-page applications. You can spin up a complete Full-Stack environment (Vodka + React Vite) in seconds.

1. **Create a Project**

Run the following command to scaffold a new project:

```bash
vodka create my-app
```

This instantly generates:

- A ready-to-use Go backend with Vodka installed.

- A lightning-fast Vite + React frontend.

- Pre-configured routing to seamlessly serve your React app from the Go backend in production.

2. **Install Dependencies**

Navigate into your new project and install the frontend packages:

```bash
cd my-app
cd frontend && npm install
cd ..
```

2. **Run the Dev Environment**

Forget running multiple terminals. The Vodka CLI manages everything for you:

```bash
vodka run dev
```

This command concurrently starts the Vite Hot-Module-Replacement (HMR) server for your React frontend and the Vodka hot-reload watcher for your Go backend. Edit a .jsx file, the browser updates instantly. Edit a .go file, the backend rebuilds in milliseconds.

## Using Vodka as a Standalone API

If you are just building a REST API or microservice without a React frontend, Vodka is still the perfect tool.

**Installing**

Install the vodka CLI

```bash
go install github.com/DevanshuTripathi/vodka/cmd/vodka@latest
```

Init a Go module and install vodka

```bash
mkdir backend-app
cd backend-app
go mod init app
go get github.com/DevanshuTripathi/vodka
```

**Hot-Reloading your API**

Inside any standard Go module using Vodka, simply type:

```bash
vodka
```

The CLI will compile your app into a temporary binary, start the server, and automatically rebuild and restart it whenever you save a .go file.

**Minimal Backend Example**

Here is a fully functional Vodka application:

```Go
package main

import (
    "log"
    "https://github.com/DevanshuTripathi/vodka"
)

func main() {
    // Initialize a new Vodka engine with default logging/recovery middleware
    app := vodka.Default()

    // Define a simple GET route
    app.GET("/ping", func(c *vodka.Context) {
        c.JSON(200, vodka.M{"message": "pong! 🏓"})
    })

    // Start the server on port 8080
    if err := app.Run(":8080"); err != nil {
        log.Fatalf("Server crashed: %v", err)
    }
}
```

## What Everything Does (Architecture)

### The Engine (Router)

At its core, vodka.Engine is a highly optimized router powered by httprouter under the hood. Instead of scanning linearly through routes, it uses a Radix Tree (Prefix Tree). This means the time it takes to find a route is proportional only to the length of the URL path, not the number of routes in your API. It handles highly complex routing logic in under a microsecond without allocating new memory.

### The Context

vodka.Context is the central struct passed to all your handlers. It wraps Go's standard http.Request and http.ResponseWriter into a single, clean object and provides quality-of-life helper methods.

- c.JSON(statusCode, obj): Automatically sets the content type and encodes your Go structs into JSON.

- c.String(statusCode, text): Sends plain text responses instantly.

- c.BindJSON(&obj): Instantly parses incoming request bodies into your custom Go structs.

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

**Full-Stack Routung(ServeSPA)**

When you use vodka create, your backend comes with app.ServeSPA("./frontend/dist") pre-configured. This acts as a smart fallback handler: if a user requests a route that isn't a registered API endpoint, Vodka will automatically serve your React index.html, allowing React Router to handle client-side navigation flawlessly without throwing 404s on page refresh.