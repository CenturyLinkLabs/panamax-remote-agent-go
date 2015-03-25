package adapter

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type servicesClient struct {
	endpoint  string
	netClient *http.Client
}

func MakeClient(ep string) Client {
	hc := &http.Client{}

	c := servicesClient{
		netClient: hc,
		endpoint:  ep,
	}

	return c
}

func (sc servicesClient) CreateServices(buf *bytes.Buffer) []Service {
	resp, _ := sc.netClient.Post(sc.servicesPath(""), "application/json", buf)

	ars := &[]Service{}
	jd := json.NewDecoder(resp.Body)
	jd.Decode(ars)

	return *ars
}

func (sc servicesClient) GetService(sid string) Service {
	resp, _ := sc.netClient.Get(sc.servicesPath(sid))

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

func (sc servicesClient) DeleteService(sid string) error {
	req, err := http.NewRequest("DELETE", sc.servicesPath(sid), nil)

	if err != nil {
		return err
	}

	_, err = sc.netClient.Do(req)

	return err
}

func (sc servicesClient) FetchMetadata() (interface{}, error) {
	res, err := sc.netClient.Get(sc.endpoint + "/v1/metadata")

	if err != nil {
		//TODO
		// return map[string]string, err
	}

	var r interface{}
	jd := json.NewDecoder(res.Body)
	jd.Decode(&r)

	return r, nil
}

func (sc servicesClient) servicesPath(id string) string {
	return sc.endpoint + "/v1/services/" + id
}
