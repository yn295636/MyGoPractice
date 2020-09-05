package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/resolver"
	"log"
)

const schema = "etcd"

type Resolver struct {
	endpoints []string
	service   string
	cli       *clientv3.Client
	cc        resolver.ClientConn
}

func NewResolver(endpoints []string, service string) resolver.Builder {
	return &Resolver{endpoints: endpoints, service: service}
}

func (r *Resolver) Scheme() string {
	return schema + "_" + r.service
}

func (r *Resolver) ResolveNow(rn resolver.ResolveNowOption) {
}

func (r *Resolver) Close() {
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var err error

	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints: r.endpoints,
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
	}

	r.cc = cc

	go r.watch(r.service)

	return r, nil
}

func (r *Resolver) watch(prefix string) {
	addrDict := make(map[string]resolver.Address)

	update := func() {
		addrList := make([]resolver.Address, 0, len(addrDict))
		for _, v := range addrDict {
			addrList = append(addrList, v)
		}
		r.cc.NewAddress(addrList)
	}

	resp, err := r.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err == nil {
		for i, kv := range resp.Kvs {
			info := &ServiceInfo{}
			if err := json.Unmarshal(kv.Value, info); err != nil {
				log.Println(err)
				continue
			}

			addrDict[string(resp.Kvs[i].Value)] = resolver.Address{Addr: info.IP}
		}
	}

	update()

	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for n := range rch {
		for _, ev := range n.Events {
			switch ev.Type {
			case mvccpb.PUT:
				info := &ServiceInfo{}
				if err := json.Unmarshal(ev.Kv.Value, info); err != nil {
					log.Println(err)
					continue
				}
				addrDict[string(ev.Kv.Key)] = resolver.Address{Addr: info.IP}
			case mvccpb.DELETE:
				delete(addrDict, string(ev.PrevKv.Key))
			}
		}
		update()
	}
}
