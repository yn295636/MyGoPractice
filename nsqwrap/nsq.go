package nsqwrap

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
)

const (
	NSQShardFormat = "%s_%08d"
)

var (
	NsqShardCount = uint32(100)
)

func NSQQueueIdx(dataId uint64) uint64 {
	return dataId % uint64(NsqShardCount)
}

var p *nsq.Producer

type NsqConfig struct {
	Topic              string
	Channel            string
	LookupdAddrs       string
	Handler            nsq.HandlerFunc
	MsgTimeout         time.Duration
	ConcurrentHandlers int
}

type NsqLogger struct {
}

func (logger *NsqLogger) Output(calldepth int, s string) error {
	log.Printf("nsqlog[%v]", s)
	return nil
}

func InitConsumer(configList []*NsqConfig) error {
	for _, v := range configList {
		cfg := nsq.NewConfig()
		if v.MsgTimeout != 0 {
			cfg.MsgTimeout = v.MsgTimeout
		}
		c, err := nsq.NewConsumer(v.Topic, v.Channel, cfg)
		if err != nil {
			log.Printf("nsq.NewConsumer error=%v", err)
			return err
		}
		c.SetLogger(&NsqLogger{}, nsq.LogLevelWarning)
		c.ChangeMaxInFlight(200)
		//c.AddHandler(nsq.HandlerFunc(v.Handler))
		if v.ConcurrentHandlers != 0 {
			c.AddConcurrentHandlers(nsq.HandlerFunc(v.Handler), v.ConcurrentHandlers)
		} else {
			c.AddConcurrentHandlers(nsq.HandlerFunc(v.Handler), 1)
		}

		go func(c *nsq.Consumer, lookupdAddrs, topic, channel string) {
			err = c.ConnectToNSQLookupds(strings.Split(lookupdAddrs, ";"))
			if err != nil {
				log.Printf("ConnectToNSQLookupds error=%v", err)
			}
			log.Printf("init nsq consumer topic=%v, channel=%v", topic, channel)
		}(c, v.LookupdAddrs, v.Topic, v.Channel)
	}
	return nil
}

func InitProducer(nsqdAddr string) error {
	cfg := nsq.NewConfig()
	var err error
	p, err = nsq.NewProducer(nsqdAddr, cfg)
	if err != nil {
		log.Printf("nsq.NewProducer error=%v", err)
		return err
	}
	p.SetLogger(&NsqLogger{}, nsq.LogLevelWarning)
	return nil
}

func Publish(topic string, body []byte) error {
	err := p.Publish(topic, body)
	return err
}

func DeferPublish(topic string, delay time.Duration, body []byte) error {
	return p.DeferredPublish(topic, delay, body)
}

func StopProducer() {
	p.Stop()
}

func JsonMarshalPublish(topic string, req interface{}) error {
	if topic == "" || req == nil {
		return errors.New("input param error")
	}
	message, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = Publish(topic, message)
	return err
}

func JsonMarshalDeferPublish(topic string, delay time.Duration, req interface{}) error {
	if topic == "" || req == nil {
		return errors.New("input param error")
	}

	message, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return DeferPublish(topic, delay, message)
}

type NsqLookupApiResponse struct {
	Channels  []string              `json:"channels"`
	Producers []*NsqLookupProducers `json:"producers"`
}

type NsqLookupProducers struct {
	RemoteAddress    string `json:"remote_address"`
	Hostname         string `json:"hostname"`
	BroadcastAddress string `json:"broadcast_address"`
	TcpPort          int    `json:"tcp_port"`
	HttpPort         int    `json:"http_port"`
	Version          string `json:"version"`
}

// nsqlookupd API (/lookup) : returns a list of producers for a topic
func ApiLookUp(ctx context.Context, nsqlookupdAddrs string, topic string) (*NsqLookupApiResponse, error) {
	addrs := strings.Split(nsqlookupdAddrs, ";")
	if len(addrs) < 1 {
		return nil, errors.New("input param error")
	}

	rurl := fmt.Sprintf("http://%s/lookup?topic=%s", addrs[0], topic)
	req, err := http.NewRequest("GET", rurl, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("call lookup api failed, body=%v", string(respBody))
	}

	var response NsqLookupApiResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, errors.Wrapf(err, "parse nsq api response body failed, body=%v", string(respBody))
	}
	return &response, nil
}

// nsqd API (/channel/delete) : delete channel info on special nsqd server
func ApiDeleteChannel(ctx context.Context, nsqdAddr string, topic string, channel string) (bool, error) {
	rurl := fmt.Sprintf("http://%s/channel/delete?topic=%s&channel=%s", nsqdAddr, topic, channel)
	req, err := http.NewRequest("POST", rurl, nil)
	if err != nil {
		return false, err
	}

	req = req.WithContext(ctx)
	req.Header.Add("Accept", "application/vnd.nsq; version=1.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return false, errors.Errorf("call nsqd delete channel api failed, body=%v", string(respBody))
	}
	return resp.StatusCode == http.StatusOK, nil
}
