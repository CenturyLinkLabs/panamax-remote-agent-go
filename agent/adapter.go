package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Adapter struct {
	Endpoint string
	Client   *http.Client
}

func NewAdapter(ep string) Adapter {
	client := &http.Client{}

	ad := Adapter{
		Client:   client,
		Endpoint: ep,
	}

	return ad
}

func (ad *Adapter) CreateServices(sIDs []Image) AdapterResponses {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(sIDs); err != nil {
		panic(err)
	}
	resp, _ := ad.Client.Post(ad.servicesPath(""), "application/json", buf)

	ars := &AdapterResponses{}
	jd := json.NewDecoder(resp.Body)
	jd.Decode(ars)

	return *ars
}

func (ad *Adapter) GetService(sid string) Service {
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

func (ad *Adapter) DeleteService(sid string) error {
	req, err := http.NewRequest("DELETE", ad.servicesPath(sid), nil)

	if err != nil {
		return err
	}

	_, err = ad.Client.Do(req)

	return err
}

func (ad *Adapter) FetchMetadata() (interface{}, error) {
	res, err := ad.Client.Get(ad.Endpoint + "/v1/metadata")

	if err != nil {
		// return map[string]string, err
	}

	var r interface{}
	jd := json.NewDecoder(res.Body)
	jd.Decode(&r)

	return r, nil
}

func (ad *Adapter) servicesPath(id string) string {
	return ad.Endpoint + "/v1/services/" + id
}

type AdapterResponses []AdapterResponse

type AdapterResponse struct {
	ID          string `json:"id"`
	ActualState string `json:"actualState,omitempty"`
}
