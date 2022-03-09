package util

import "os"

var Debug = os.Getenv("DEBUG") == "true"
