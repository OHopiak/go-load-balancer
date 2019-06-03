package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OHopiak/fractal-load-balancer/core"
	"net/http"
	"time"
)

func (w *Worker) Register() error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	request := core.RegisterWorkerRequest{
		Port: w.Port,
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(request)
	if err != nil {
		return err
	}
	r, err := client.Post(fmt.Sprintf("http://%s/worker/register", w.Master), "application/json; charset=utf-8", b)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusCreated {
		return errors.New("failed to connect to master")
	}

	response := core.RegisterWorkerResponse{}
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return err
	}

	w.Server.Logger.Infof("Connected to master. Status: '%s'", response.Status)
	w.Master.Connected = true
	*w.ID = response.WorkerId
	return nil
}
