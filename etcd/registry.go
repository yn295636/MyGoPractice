package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

type ServiceInfo struct {
	Name string
	IP   string
}

type Service struct {
	ServiceInfo ServiceInfo
	stop        chan error
	leaseId     clientv3.LeaseID
	client      *clientv3.Client
}

func NewService(info ServiceInfo, endpoints []string) (service *Service, err error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 200,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	service = &Service{
		ServiceInfo: info,
		client:      client,
	}
	return
}

func (info *ServiceInfo) getKey() string {
	return info.Name + "/" + info.IP
}

func (service *Service) Start() (err error) {

	ch, err := service.keepAlive()
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		select {
		case err := <-service.stop:
			return err
		case <-service.client.Ctx().Done():
			return errors.New(fmt.Sprintf("service %s closed", service.ServiceInfo.getKey()))
		case _, ok := <-ch:
			if !ok {
				log.Printf("keep alive channel of service %s closed", service.ServiceInfo.getKey())
				return service.revoke()
			}
		}
	}

	return
}

func (service *Service) Stop() {
	service.stop <- nil
}

func (service *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	info := &service.ServiceInfo
	key := service.ServiceInfo.getKey()
	val, _ := json.Marshal(info)

	resp, err := service.client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = service.client.Put(context.TODO(), key, string(val), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	service.leaseId = resp.ID
	return service.client.KeepAlive(context.TODO(), resp.ID)
}

func (service *Service) revoke() error {
	_, err := service.client.Revoke(context.TODO(), service.leaseId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("servide:%s stop\n", service.ServiceInfo.getKey())
	return err
}
