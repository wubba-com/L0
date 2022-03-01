package publisher

import (
	"encoding/json"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/nats"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const(
	channel   = "test"
	file = "model.json"
)

func Run()  {
	order := &domain.Order{}
	wd, err := os.Getwd()
	log.Println(wd)

	if err != nil {
		log.Printf("err: %s\n", err.Error())
		return
	}

	f, err := os.Open(filepath.Join(wd, file))
	if err != nil {
		log.Printf("err: %s\n", err.Error())
		return
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("err: %s\n", err.Error())
		return
	}
	err = json.Unmarshal(b, order)
	if err != nil {
		log.Printf("err: %s\n", err.Error())
		return
	}

	sc := nats.NewStanConn("ClusterID", "clientID")
	err = sc.Publish(channel, b)
	if err != nil {
		log.Printf("err: %s\n", err.Error())
		return
	}

	err = sc.Close()
	if err != nil {
		log.Printf("err: %s\n", err.Error())
		return
	}
}
