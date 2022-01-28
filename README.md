# Zonnepanelendelen API Golang client

[![Go Report Card](https://goreportcard.com/badge/github.com/skoef/go-zonnepanelendelen)](https://goreportcard.com/report/github.com/skoef/go-zonnepanelendelen) [![Documentation](https://godoc.org/github.com/skoef/go-zonnepanelendelen?status.svg)](http://godoc.org/github.com/skoef/go-zonnepanelendelen)

This is a simple golang library for interfacing with the API of the investment project [Zonnepanelendelen](https://zonnepanelendelen.nl). Zonnepanelendelen allows you to invest in solar panels when you don't have any more room left on your roof or can't put any solar panels on your roof at all.

For simplicity sake, currently only several API endpoints are supported. If you miss specific features in the library, please open an issue!

## Example usage

An example for using the API client is shown below, where the credentials are those you would login to https://mijnstroom.zonnepanelendelen.nl/ with:

```golang
package main

import (
  "fmt"

  zonnepanelendelen "github.com/skoef/go-zonnepanelendelen"
)

func main() {
  api := zonnepanelendelen.New("johndoe", "s3cr3t")
  projects, err := api.GetProjects()
  if err != nil {
    panic(err)
  }

  fmt.Printf("found %d projects in your Zonnepanelendelen account\n", len(projects))
}
```
