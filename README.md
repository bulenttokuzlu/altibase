# GORM Altibase Driver



## Description

GORM Oracle driver for connect Altibase DB and Manage Altibase DB, Based on
，not recommended for use in a production environment

## Required dependency Install

- Altibase
- Golang 1.13+
- see [ODPI-C Installation.](https://oracle.github.io/odpi/doc/installation.html)

## Quick Start
### how to install 
```bash
go get github.com/bulenttokuzlu/altibase
```
###  usage

```go
import (
	"fmt"
	"github.com/bulenttokuzlu/altibase"
	"gorm.io/gorm"
	"log"
)

func main() {
    db, err := gorm.Open(oracle.Open("sys/altibase@172.20.1.80:31114/mydb"), &gorm.Config{})
    if err != nil {
        // panic error or log error info
    } 
    
    // do somethings
}
```





go clean -modcache
go mod tidy





go mod init github.com/bulenttokuzlu/altibase