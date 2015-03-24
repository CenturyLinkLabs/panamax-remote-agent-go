package adapter

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	Endpoint string
	Client   *http.Client
}

func NewClient(ep string) Client {
	client := &http.Client{}

	ad := Client{
		Client:   client,
		Endpoint: ep,
	}

	return ad
}

func (ad Client) CreateServices(buf *bytes.Buffer) []Service {
	resp, _ := ad.Client.Post(ad.servicesPath(""), "application/json", buf)

	ars := &[]Service{}
	jd := json.NewDecoder(resp.Body)
	jd.Decode(ars)

	return *ars
}

func (ad Client) GetService(sid string) Service {
	resp, _ := ad.Client.Get(ad.servicesPath(sid))

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

func (ad Client) DeleteService(sid string) error {
	req, err := http.NewRequest("DELETE", ad.servicesPath(sid), nil)

	if err != nil {
		return err
	}

	_, err = ad.Client.Do(req)

	return err
}

func (ad Client) FetchMetadata() (interface{}, error) {
	res, err := ad.Client.Get(ad.Endpoint + "/v1/metadata")

	if err != nil {
		// return map[string]string, err
	}

	var r interface{}
	jd := json.NewDecoder(res.Body)
	jd.Decode(&r)

	return r, nil
}

func (ad Client) servicesPath(id string) string {
	return ad.Endpoint + "/v1/services/" + id
}
