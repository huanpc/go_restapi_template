package storage

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/olivere/elastic.v5"
)

// Structures used for serializing/deserializing data in Elasticsearch.
type OnConnectEvent struct {
	Call		string                `json:"call"`
	TcUrl     	string                `json:"tc_url,omitempty"`
	Addr		string                `json:"addr"`
	App    		string                `json:"app,omitempty"`
	FlashVer    string                `json:"flash_ver,omitempty"`
	SwfUrl		string                `json:"swf_url,omitempty"`
	PageUrl		string                `json:"page_url,omitempty"`
	Created  	time.Time             `json:"created,omitempty"`
}

type Event struct {
	Data *OnConnectEvent
	IndexName string
	Type string
	Id string
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"event":{
			"properties":{
				"call":{
					"type":"keyword"
				},
				"tc_url":{
					"type":"keyword"
				},
				"addr":{
					"type":"ip"
				},
				"app":{
					"type":"keyword"
				},
				"flash_ver":{
					"type":"keyword"
				},
				"swf_url":{
					"type":"ip"
				},
				"extra":{
					"type":"keyword"
				},
				"page_url":{
					"type":"keyword"
				},
				"stream_name":{
					"type":"keyword"
				},
				"created":{
					"type":"date"
				}
			}
		}
	}
}`

type EsClient struct {
	Client *elastic.Client
	Context context.Context
}

// init new elastic cli instance
func NewEsClient(esAddress string, context context.Context) EsClient {
	client, err := elastic.NewSimpleClient(elastic.SetURL(esAddress))
	if err != nil {
		// Handle error
		panic(err)
	}
	return EsClient{
		Client: client,
		Context: context,
	}
}

/* CreateIndexIfNotExists */
func (esClient *EsClient) CreateIndexIfNotExists(item *Event) {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := esClient.Client.IndexExists(item.IndexName).Do(esClient.Context)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := esClient.Client.CreateIndex(item.IndexName).BodyString(mapping).Do(esClient.Context)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

/* Index item using json serialization */
func (esClient *EsClient) IndexItem(item *Event) {
	esClient.CreateIndexIfNotExists(item)
	put1, err := esClient.Client.Index().
		Index(item.IndexName).
		Type(item.Type).
		Id(item.Id).
		BodyJson(item.Data).
		Do(esClient.Context)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}

/* TODO Update */
func (esClient *EsClient) UpdateItem(item *Event){
	update, err := esClient.Client.Update().Index(item.IndexName).Type(item.Type).Id(item.Id).
		Script(elastic.NewScriptInline("ctx._source.retweets += params.num").Lang("painless").Param("num", 1)).
		Upsert(map[string]interface{}{"retweets": 0}).
		Do(esClient.Context)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("New version of tweet %q is now %d\n", update.Id, update.Version)

}

func (esClient *EsClient) DeleteItem(esclient *elastic.Client, indexName string, context context.Context){}

func (esClient *EsClient) DeleteIndex(esclient *elastic.Client, indexName string, context context.Context){
	// Delete an index.
	deleteIndex, err := esclient.DeleteIndex(indexName).Do(context)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}
}

func main() {

	ctx := context.Background()

	client := NewEsClient("http://localhost:9200", ctx)

	item := &Event {
		Id: "2",
		IndexName: "event",
		Type: "streaming_event",
		Data: &OnConnectEvent{},
	}

	client.IndexItem(item)
}