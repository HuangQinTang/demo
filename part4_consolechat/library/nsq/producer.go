package nsq

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nsqio/go-nsq"
	"chat/library/logicerror"
	"io/ioutil"
	"net/http"
	"time"
)

//生产者对象
type NsqProducer struct {
	producer *nsq.Producer
	addr     string
}

func NewNsqProducer(addr string) (*NsqProducer, error) {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	err = producer.Ping()
	if nil != err {
		producer.Stop()
		return nil, err
	}

	return &NsqProducer{
		producer: producer,
		addr:     addr,
	}, nil
}

//检测当前连接状态（消费者掉线会自动重连，这个包的生产者不会）
func (p *NsqProducer) CheckConn() {
	ticker := time.NewTicker(60 * time.Second) //1分钟检查一次连接状态
	for range ticker.C {
		err := p.producer.Ping()
		//连接状态正常
		if err == nil {
			continue
		}

		//连接异常重新连接
		p.producer.Stop()
		newConn, err := NewNsqProducer(p.addr)
		if logicerror.PrintError(err) != nil {
			continue
		}
		p.producer = newConn.producer
	}
}

//生产者延迟发布消息
func (p *NsqProducer) DeferredPublish(topic string, delay time.Duration, message string) (err error) {
	if p.producer != nil {
		if message == "" { //不能发布空串，否则会导致error
			return nil
		}
		err = p.producer.DeferredPublish(topic, delay, []byte(message)) // 发布消息
		return err
	}
	return fmt.Errorf("producer is nil", err)
}

//生产者发布消息
func (p *NsqProducer) Publish(topic string, message string) (err error) {
	if p.producer != nil {
		if message == "" { //不能发布空串，否则会导致error
			return nil
		}
		err = p.producer.Publish(topic, []byte(message)) // 发布消息
		return err
	}
	return fmt.Errorf("producer is nil", err)
}

//生产者停止连接
func (p *NsqProducer) Close() {
	p.producer.Stop()
}

//http发送消息
func HttpPush(nsqdAddr, topic, mes string) bool {
	myteMes := []byte(mes)
	buffer := bytes.NewBuffer(myteMes)
	//request, err := http.NewRequest("POST", nsqdAddr + "/pub?topic=" + topic,  strings.NewReader(mes))
	request, err := http.NewRequest("POST", nsqdAddr + "/pub?topic=" + topic, buffer)
	if logicerror.PrintError(err) != nil {
		return false
	}
	client := http.Client{}
	resp, err := client.Do(request.WithContext(context.TODO()))
	if logicerror.PrintError(err) != nil {
		return false
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if logicerror.PrintError(err) != nil {
		return false
	}
	if string(respBytes) == "OK" {
		return true
	}
	return false
}
