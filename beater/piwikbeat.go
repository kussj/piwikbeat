package beater

import (
	fmt "fmt"
	time "time"
	strconv "strconv"

	beat "github.com/elastic/beats/libbeat/beat"
	cfgfile "github.com/elastic/beats/libbeat/cfgfile"
	common "github.com/elastic/beats/libbeat/common"
	logp "github.com/elastic/beats/libbeat/logp"

	config "github.com/kussj/piwikbeat/config"
	pbcommon "github.com/kussj/piwikbeat/common"
)

type Piwikbeat struct {
	beatConfig	*config.Config
	done		chan struct{}
	period		time.Duration
	url			string
	token		string
	methods		[]pbcommon.EndPoint
	siteid		string
}

// Creates beater
func New() *Piwikbeat {
	return &Piwikbeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Piwikbeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := cfgfile.Read(&bt.beatConfig, "")
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	var url string
	if bt.beatConfig.Piwikbeat.Url != "" {
		bt.url = bt.beatConfig.Piwikbeat.Url
		logp.Debug("Configured URL: %v\n", url)
	} else {
		logp.Err("URL endpoint not configured")
	}

	if bt.beatConfig.Piwikbeat.SiteID != 0 {
		bt.siteid = strconv.Itoa(bt.beatConfig.Piwikbeat.SiteID)
	} else {
		logp.Err("Site ID not configured")
	}

	if bt.beatConfig.Piwikbeat.Token != "" {
		bt.token = bt.beatConfig.Piwikbeat.Token
	} else {
		logp.Err("Auth Token not specified")
	}

	var methods []pbcommon.EndPoint
	if bt.beatConfig.Piwikbeat.Methods != nil {
		methods = bt.beatConfig.Piwikbeat.Methods
	} else {
		logp.Err("No methods configured")
	}

	bt.methods = make([]pbcommon.EndPoint, len(methods))
	for i := 0; i < len(methods); i++ {
		m := methods[i]
		bt.methods[i] = m
	}

	return nil
}

func (bt *Piwikbeat) Setup(b *beat.Beat) error {

	// Setting default period if not set
	if bt.beatConfig.Piwikbeat.Period == "" {
		bt.beatConfig.Piwikbeat.Period = "10s"
	}

	var err error
	bt.period, err = time.ParseDuration(bt.beatConfig.Piwikbeat.Period)
	if err != nil {
		return err
	}

	return nil
}

func (bt *Piwikbeat) Run(b *beat.Beat) error {
	logp.Info("piwikbeat is running! Hit CTRL-C to stop it.")

	ticker := time.NewTicker(bt.period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		metrics, err := bt.GetMetrics(bt.url, bt.siteid, bt.token, bt.methods)

		if err != nil {
			logp.Err("Error reading metrics from %v", err)
		}
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"counter":    counter,
			"metrics":	  metrics,
			"siteid:":	  bt.siteid,
		}
		b.Events.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Piwikbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Piwikbeat) Stop() {
	close(bt.done)
}
