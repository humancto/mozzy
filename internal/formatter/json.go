package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func PrintJSONOrText(b []byte, jqQuery string) error {
	// v1: pretty-print JSON if possible; ignore jqQuery for now (hook to replace)
	var out bytes.Buffer
	if json.Indent(&out, b, "", "  ") == nil {
		_, _ = out.WriteTo(os.Stdout)
		fmt.Println()
		return nil
	}
	// Not JSON — print raw
	fmt.Println(string(b))
	return nil
}
