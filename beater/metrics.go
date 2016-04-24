package beater

import (
	url "net/url"
	http "net/http"
	tls "crypto/tls"
	json "encoding/json"
	logp "github.com/elastic/beats/libbeat/logp"
	pbcommon "github.com/kussj/piwikbeat/common"
)

func (pb *Piwikbeat) GetMetrics(u string, id string, token string, methods []pbcommon.EndPoint) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	for _,m := range methods {
		method := m.GetMethod()
		params := m.GetParameters()
		var Url *url.URL
		Url, err := url.Parse(u)

		if err != nil {
			logp.Err("An error occurred while parsing the provided URL: %v", err)
		}

		parameters := url.Values{}
		parameters.Add("module", "API")
		parameters.Add("method", method)
		parameters.Add("idSite", id)
		parameters.Add("token_auth", token)
		parameters.Add("format", "JSON")
		for k,v := range params {
			if k == "date" {
				v = "today"
			}
			parameters.Add(k,v)
		}
		Url.RawQuery = parameters.Encode()

		var resp interface{}
		err = getJson(Url.String(), &resp)

		if err != nil {
			logp.Err("An error occurred while reading data from Piwik: %v", err)
		} else {
			switch resp.(type) {
			case []string:
				r := resp.([]string)
				for i := 0; i < len(r); i++ {
					key := method
					val := r[i]
					metrics[key] = val
				}
			case []interface{}:
				r := resp.([]interface{})
				for i := 0; i < len(r); i++ {
					switch r[i].(type) {
					case string:
						key := method
						val := r[i]
						metrics[key] = val
					case map[string]interface{}:
						rr := r[i].(map[string]interface{})
						for k,v := range rr{
							key := (method + ":" + k)
							val := v
							metrics[key] = val
						}
					}
				}
			case map[string]interface{}:
				r := resp.(map[string]interface{})
				for k,v := range r {
					key := (method + ":" + k)
					val := v
					metrics[key] = val
				}
			}
		}


	}//for range methods

	return metrics, nil
}//GetMetrics


func getJson(url string, target interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}//Transport

	client := &http.Client{Transport: tr}
	r, err := client.Get(url)

	if err != nil {
		logp.Err("An error occurred while executing HTTP request: %v", err)
		return err
	}//if

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}//getJson

/**
https://analytics.cahcommtech.com/?module=API&\
method=Live.getCounters&\
idSite=7&\
lastMinutes=30&\
format=JSON&\
token_auth=bc5961932268c4d190d5a0893bcbd454

[{"visits":"0","actions":0,"visitors":0,"visitsConverted":0}]
*/