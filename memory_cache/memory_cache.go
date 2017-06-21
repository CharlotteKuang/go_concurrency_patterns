package main

import ()

type CacheGetter struct {
	Key      string
	Receiver chan *WrapGetterData
}

type CacheWriter struct {
	Key   string
	Value interface{}
}

type WrapGetterData struct {
	Value interface{}
	Exist bool
}

type CacheData struct {
	Value interface{}
}

type MCache struct {
	Cache    map[string]*CacheData
	ReaderCh chan *CacheGetter
	WriterCh chan *CacheWriter
}

func (this *MCache) Init() {
	this.Cache = make(map[string]*CacheData)
	this.ReaderCh = make(chan *CacheGetter)
	this.WriterCh = make(chan *CacheWriter)
	go this.run()
}

func (this *MCache) run() {
	for {
		select {
		case getter := <-this.ReaderCh:
			cache, exist := this.Cache[getter.Key]
			data := new(WrapGetterData)
			if exist {
				data.Value = cache.Value
				data.Exist = true
			} else {
				data.Value = nil
				data.Exist = false
			}
			getter.Receiver <- data
		case writer := <-this.WriterCh:
			data := new(CacheData)
			data.Value = writer.Value
			this.Cache[writer.Key] = data
		}
	}
}

func (this *MCache) Get(key string) (interface{}, bool) {
	getter := new(CacheGetter)
	getter.Key = key
	getter.Receiver = make(chan *WrapGetterData)
	this.ReaderCh <- getter
	rst := <-getter.Receiver
	return rst.Value, rst.Exist
}

func (this *MCache) Set(key string, value interface{}) {
	writer := new(CacheWriter)
	writer.Key = key
	writer.Value = value
	this.WriterCh <- writer
}

func NewMCache() *MCache {
	cache := new(MCache)
	cache.Init()
	return cache
}
