package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/kussj/piwikbeat/beater"
)

func main() {
	err := beat.Run("piwikbeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
