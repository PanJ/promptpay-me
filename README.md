# PromptPay Me

A web server written in Go to generate PromptPay QR code image based on url

## Usage

### Server

Just run the server with environment variable named `PROMPTPAY_RECIPIENT_ID` indicating recipient address.

  `$ go run cmd/promptpay-server/main.go`

### Library

```go
package main

import (
	"fmt"

	promptpay "github.com/PanJ/promptpay-me"
)

func main() {
	payload := promptpay.Generate("0800000000", 1000.0)
	fmt.Println(payload)
}
```

## Reference

- [promptpay-qr](https://github.com/dtinth/promptpay-qr) by dtinth
