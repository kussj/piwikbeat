// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	pbcommon "github.com/kussj/piwikbeat/common"
)

type Config struct {
	Piwikbeat PiwikbeatConfig
}

type PiwikbeatConfig struct {
	Period 	string `yaml:"period"`
	Url	string `yaml:"url"`
	Token	string `yaml:"token"`
	Methods []pbcommon.EndPoint `yaml:"methods"`
	SiteID	int `yaml:"site_id"`
}
