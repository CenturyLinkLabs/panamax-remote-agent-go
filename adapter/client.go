package adapter

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type client struct {
	endpoint string
	client   *http.Client
}

func MakeClient(ep string) AdapterClient {
	hc := &http.Client{}

	c := client{
		client:   hc,
		endpoint: ep,
	}

	return c
}

func (ad client) CreateServices(buf *bytes.Buffer) []Service {
	resp, _ := ad.client.Post(ad.servicesPath(""), "application/json", buf)

	ars := &[]Service{}
	jd := json.NewDecoder(resp.Body)
	jd.Decode(ars)

	return *ars
}

func (ad client) GetService(sid string) Service {
	resp, _ := ad.client.Get(ad.servicesPath(sid))

	if resp.StatusCode == http.StatusNotFound {
		return Service{ID: sid, ActualState: "not found"}
	} else if resp.StatusCode != http.StatusOK {
		return Service{ID: sid, ActualState: "error"}
	}

	srvc := &Service{}
	jd := json.NewDecoder(resp.Body)
	jd.Decode(srvc)

	return *srvc
}

func (ad client) DeleteService(sid string) error {
	req, err := http.NewRequest("DELETE", ad.servicesPath(sid), nil)

	if err != nil {
		return err
	}

	_, err = ad.client.Do(req)

	return err
}

func (ad client) FetchMetadata() (interface{}, error) {
	res, err := ad.client.Get(ad.endpoint + "/v1/metadata")

	if err != nil {
		//TODO
		// return map[string]string, err
	}

	var r interface{}
	jd := json.NewDecoder(res.Body)
	jd.Decode(&r)

	return r, nil
}

func (ad client) servicesPath(id string) string {
	return ad.endpoint + "/v1/services/" + id
}
