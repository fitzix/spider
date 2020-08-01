package services

import (
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fitzix/spider/model"
	"github.com/fitzix/spider/pkg/logger"
	"github.com/fitzix/spider/pkg/weibo"
	"gopkg.in/yaml.v3"
)

func Weibo() {
	bytes, err := ioutil.ReadFile("weibo.yml")
	if err != nil {
		logger.Fatalf("read conf file err: %s", err)
	}

	var conf model.WeiboConf

	if err := yaml.Unmarshal(bytes, &conf); err != nil {
		logger.Fatalf("unmarshal conf err: %s", err)
	}

	weibo.Init(conf.Token)

	for k, v := range conf.Comments {
		go watch(k, conf.Tick, append(v, conf.Common.Comments...))
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func watch(uid string, tick int, messages []string) {
	index := 0
	commentLen := len(messages)
	register := make(map[string]bool)
	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	for {
		<-ticker.C
		mid, err := weibo.GetLatestWeibo(uid)
		if err != nil {
			logger.Errorf("获取最新微博失败-->%s", err)
			continue
		}

		if register[mid] {
			continue
		}

		if index == commentLen {
			index = 0
		}

		msg := messages[index]
		logger.Infof("%s === 发现新微博, 开始评论", uid)

		if err := weibo.CreateComment(mid, msg); err != nil {
			logger.Infof("%s === 发送评论失败--->%s", uid, err)
			continue
		}
		index++
		register[mid] = true
		logger.Infof("%s === 评论成功 === %s", uid, msg)
	}
}
