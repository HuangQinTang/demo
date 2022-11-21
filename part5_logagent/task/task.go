package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	etcd "go.etcd.io/etcd/clientv3"
	"log"
	"logagent/common"
	"logagent/tailservice"
	"strings"
)

const EtcdCollectKey = "collect_log_%s_conf"

type CollectEntry struct {
	Path  string `json:"path"`  //要读取的日志路径
	Topic string `json:"topic"` //推送到kafka的哪个主题下
}

type Task struct {
	EtcdClient  *etcd.Client        //etcd客户端
	Tails       map[string]tailTask //文件读取对象集合
	KafkaClient sarama.SyncProducer //kafka生产者客户端
}

func (t *Task) Run() {
	//1.从ectd中读取要收集的日志信息
	ip := common.GetOutboundIP()
	etcdCollectKey := fmt.Sprintf(EtcdCollectKey, ip)
	collectInfo, err := t.EtcdClient.Get(context.Background(), etcdCollectKey)
	if err != nil {
		log.Fatalf("read etcd key failed err: %s", err.Error())
	}
	if len(collectInfo.Kvs) == 0 {
		log.Fatal("etcd 没有配置要读取的日志")
	}

	var collects []CollectEntry
	collectCfg := collectInfo.Kvs[0]
	err = json.Unmarshal(collectCfg.Value, &collects)
	if err != nil {
		log.Fatalf("tailsCfg unmarshal failed err: 配置信息格式不正确 %s", err.Error())
	}

	//2.创建kafka发送消息对象
	kafkaSendMes := t.sendMes()

	//3.根据ectd中获取的信息创建日志收集对象,并收集日志
	t.Tails = make(map[string]tailTask, len(collects))
	for _, collect := range collects {
		tailTaskObj, err := t.createTailTask(collect.Path, collect.Topic, kafkaSendMes)
		if err != nil {
			log.Printf("read %s failed err: %s", collect.Path, err.Error())
			continue
		}
		t.Tails[collect.Path] = tailTaskObj
	}

	//4.监听etcd中设置的日志配置，并根据变更重新维护读取日志对象
	watchCh := t.EtcdClient.Watch(context.Background(), etcdCollectKey)
	for wresp := range watchCh {
		log.Printf("etcd key %s update", etcdCollectKey)
		for _, evt := range wresp.Events {
			switch evt.Type {
			case etcd.EventTypeDelete: //停止当前所有日志读取任务
				for path, task := range t.Tails {
					task.cancel()
					delete(t.Tails, path)
				}
			case etcd.EventTypePut: //更新日志读取任务
				err = json.Unmarshal(evt.Kv.Value, &collects)
				if err != nil {
					log.Printf("etcd collect format error, update failed: 更新配置信息失败 %s", err.Error())
					continue
				}
				t.updateTails(collects, kafkaSendMes)
			}
		}
	}
}

// sendMes 返回一个管道，往管道中发送的消息将投递到kafka
func (t *Task) sendMes() chan<- *sarama.ProducerMessage {
	msgChat := make(chan *sarama.ProducerMessage)
	go func() {
		for {
			msg := <-msgChat
			pid, offset, err := t.KafkaClient.SendMessage(msg)
			if err != nil {
				log.Printf("message send failed err: %s", err.Error())
			}
			log.Printf("message send success pid: %v offset: %v\n", pid, offset)
		}
	}()
	return msgChat
}

// updateTails 重新维护文件读取对象集合
func (t *Task) updateTails(newData []CollectEntry, kafkaSendMes chan<- *sarama.ProducerMessage) {
	//遍历获取当前正在进行的日志读取任务
	nowTails := make([]string, 0, len(t.Tails))
	for path, _ := range t.Tails {
		nowTails = append(nowTails, path)
	}

	//遍历新配置信息
	newTails := make([]string, 0, len(newData))
	for _, item := range newData {
		newTails = append(newTails, item.Path)
		//如果新配置信息不在当前任务中，创建新任务
		_, ok := t.Tails[item.Path]
		if !ok {
			tailTaskObj, err := t.createTailTask(item.Path, item.Topic, kafkaSendMes)
			if err != nil {
				log.Printf("read %s failed err: %s", item.Path, err.Error())
				continue
			}
			t.Tails[item.Path] = tailTaskObj
		}
	}

	//如果当前正在进行的任务不在新的配置信息中，停止执行
	deleteTails := common.Difference(nowTails, newTails)
	if len(deleteTails) > 0 {
		for _, path := range deleteTails {
			t.Tails[path].cancel()
			delete(t.Tails, path)
		}
	}
}

// tailTask 文件读取对象
type tailTask struct {
	path   string //文件路径
	tail   *tail.Tail
	topic  string //要发送到kafka哪个主题下
	ctx    context.Context
	cancel context.CancelFunc
}

// createTailTask 返回文件读取对象
func (t *Task) createTailTask(path string, topic string, kafkaSendMes chan<- *sarama.ProducerMessage) (tailTask, error) {
	tailObj, err := tailservice.InitTail(path)
	if err != nil {
		return tailTask{}, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	tailTaskObj := tailTask{
		path:   path,
		tail:   tailObj,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}

	//起协程读取文件并发送到kafka
	go tailTaskObj.readLog(kafkaSendMes)
	return tailTaskObj, nil
}

// readLog 读取文件
func (tk *tailTask) readLog(msgChat chan<- *sarama.ProducerMessage) {
	for {
		select {
		case line, ok := <-tk.tail.Lines:
			if !ok { //读取错误
				log.Printf("tail file close reopen, filename:%s\n", tk.tail.Filename)
				continue
			}
			if len(strings.Trim(line.Text, "\r")) == 0 { //空行不处理
				continue
			}

			msg := &sarama.ProducerMessage{
				Topic: tk.topic,
				Value: sarama.StringEncoder(line.Text),
			}
			msgChat <- msg
		case <-tk.ctx.Done():
			log.Printf("path:%s is stopping...", tk.path)
			return
		}
	}
}
