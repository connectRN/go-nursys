# go-nursys
Nursys API client SDK for Go

API Version
===========
Based on Nursys e-Notify File and API Specifications v3.1.2

Installation
============

To install, use `go get`:

    go get github.com/connectRN/go-nursys

Import the `github.com/connectRN/go-nursys` package into your code.


Usage
=====

```go
package yours

import (
  "github.com/connectRN/go-nursys"
)

func DoSomething() {
	nursysClient := nursys.New(server.URL, "acme", "1234!")

	req := nursys.ChangePasswordSubmitRequestMessage{
		NewPassword: "MyN3wPass!",
	}

	resp, err := nursysClient.ChangePassword(ctx, req)
}
```
