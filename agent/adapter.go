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

func (ad *Adapter) CreateServices(sIDs []Image) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(sIDs); err != nil {
		panic(err)
	}
	ad.Client.Post(ad.servicesPath(""), "application/json", buf)
}

func (ad *Adapter) servicesPath(id string) string {
	return ad.Endpoint + "/v1/services/" + id
}
