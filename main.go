/**

**squat**

  - it's like cat for sqs
  - reads from stdin and emits aws sqs messages
**/

package main

import (
	"fmt"
	"os"

	squat "github.com/nod/squat/squat"
)

func main() {
	cfg := squat.BuildRuntimeConfig()
	sq, err := squat.NewSquat(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR %s", err)
		os.Exit(1)
	}
	sq.Run()
}
