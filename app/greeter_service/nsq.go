package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/yn295636/MyGoPractice/nsqwrap"
	"github.com/yn295636/MyGoPractice/proto/greeter_service"
)

const (
	NsqTopicTest            = "test_topic"
	NsqTopicUpdateUserPhone = "update_user_phone_topic"
	NsqConsGreeterSrv       = "greeter_service"
)

func InitNsq(nsqAddr, nsqLookupAddrs string) {
	if err := nsqwrap.InitProducer(nsqAddr); err != nil {
		log.Printf("Init nsq got error %v", err)
		return
	}
	//startPublishTestTopic()
	nsqConfigs := []*nsqwrap.NsqConfig{
		{
			Topic:        NsqTopicTest,
			Channel:      NsqConsGreeterSrv,
			LookupdAddrs: nsqLookupAddrs,
			Handler: nsq.HandlerFunc(func(message *nsq.Message) error {
				log.Printf("Nsq consuming msg \"%v\"", string(message.Body))
				return nil
			}),
			MsgTimeout:         0,
			ConcurrentHandlers: 0,
		},
		{
			Topic:              NsqTopicUpdateUserPhone,
			Channel:            NsqConsGreeterSrv,
			LookupdAddrs:       nsqLookupAddrs,
			Handler:            nsq.HandlerFunc(updateUserPhoneMsgHandler),
			MsgTimeout:         0,
			ConcurrentHandlers: 0,
		},
	}
	if err := nsqwrap.InitConsumer(nsqConfigs); err != nil {
		log.Printf("Init nsq got error %v", err)
		return
	}
}

func startPublishTestTopic() {
	go func() {
		for {
			if err := nsqwrap.Publish(NsqTopicTest, []byte(fmt.Sprintf("test message %v", time.Now()))); err != nil {
				log.Printf("Nsq producer publishing messages got error %v", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func updateUserPhoneMsgHandler(msg *nsq.Message) error {
	storePhoneInDbRequest := &greeter_service.StorePhoneInDbRequest{}
	if err := proto.Unmarshal(msg.Body, storePhoneInDbRequest); err != nil {
		log.Printf("updateUserPhoneMsgHandler unmarshal nsq msg got error %v", err)
		return err
	}
	log.Printf("updateUserPhoneMsgHandler receives %v", storePhoneInDbRequest)
	if _, err := StorePhone(storePhoneInDbRequest); err != nil {
		log.Printf("updateUserPhoneMsgHandler storing phone in DB got error %v", err)
		return err
	}
	time.Sleep(3 * time.Second)
	return nil
}
