# GitHub Webhook Server Module

This is a small module you can use to create a process (or embed in an existing process) that can process Webhook
pushes from GitHub. This module is useful for CI/CD automations, static websites/blogs, etc.

It handles all the plumbing needed to listen for the webhook pushes, and provides a simple handler interface
for the application to provide the specific logic needed.

Most things are configurable when creating the server. The configuration options are:

| Name                         | Type     | Description                                                                                              |
|------------------------------|----------|----------------------------------------------------------------------------------------------------------|
| bindAddr                     | string   | The address to which the server will bind, e.g. `0.0.0.0:3000`                                           |
| urlPath                      | string   | The endpoint URL path to accept pushes, e.g. `/payload`                                                  |
| secretEnvVarName<sup>1</sup> | string   | The name of the environment variable that contains the secret key to verify the payload, e.g. `GHWH_KEY` |
| maxPayloadSize               | int      | The maximum payload size the server will accept in bytes                                                 |
| handlerFunc<sup>2</sup>      | function | The function to be called when a push is received, and validated                                         |

Notes:

<sup>1</sup> If the `secretEnvVarName` is the empty string, signature verification will not take place.

<sup>2</sup> the handler function is passed a `map[string]interface{}` and is expected to return an `error` or `nil`

## Example Usage

```
package main

import "github.com/luciddev/github_webhook"

func main() {
	server, err := github_webhook.NewServer("0.0.0.0:3000", "/payload, "GHWH_SECRET, 50*1024, handler)
	if err != nil {
		log.Printf("failed to create github webhook server, reason: %v", err)
		os.Exit(-1)
	}
	if err := server.Serve(); err != nil {
		log.Printf("failed to start github webhook server, reason: %v", err)
		os.Exit(-2)
	}
}

func handler(payload map[string]interface{}) error {
	fmt.Printf("received: %+v\n", payload)
	return nil
}

```