package controller

import (
	"context"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/leizhu/incidents_tenant/logutil"
	elastic "gopkg.in/olivere/elastic.v5"
	"io/ioutil"
	"os"
)

func InitLog(loglevel string) {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	switch loglevel {
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.AddHook(logutil.ContextHook{})
}

type (
	TenantController struct {
		ElasticsearchURL string
		Operation        string
		Tenant           string
		NumberOfShards   int
		NumberOfReplicas int
		IndexConfig      string
	}
)

func NewTenantController(es_url string, operation string, tenant string, number_of_shards int, number_of_replicas int) *TenantController {
	return &TenantController{ElasticsearchURL: es_url, Operation: operation, Tenant: tenant, NumberOfShards: number_of_shards, NumberOfReplicas: number_of_replicas}
}

func (tc *TenantController) ConfigIndex() {
	bytes, err := ioutil.ReadFile("/opt/elasticsearch/index.json")
	if err != nil {
		log.Error("Read index config /opt/incidentstenant/index.json error: ", err.Error())
		return
	}
	tc.IndexConfig = string(bytes)
}

func (tc TenantController) Test() {}

func (tc TenantController) es_client(index string) (*elastic.Client, context.Context, error) {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(tc.ElasticsearchURL))
	if err != nil {
		log.Error("Can not create es client: " + err.Error())
		return nil, context.TODO(), errors.New(fmt.Sprintln("Can not create es client: ", err))
	}
	info, code, err := client.Ping(tc.ElasticsearchURL).Do(ctx)
	if err != nil {
		log.Error("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
		return nil, context.TODO(), errors.New(fmt.Sprintln("Elasticsearch returned with code %d and version %s", code, info.Version.Number))
	}
	return client, ctx, nil
}

func (tc TenantController) create_tenant(ctx context.Context, index string, client *elastic.Client) {
	log.Info("Creating index " + index)
	createIndex, err := client.CreateIndex(index).Body(tc.IndexConfig).Do(ctx)
	if err != nil {
		log.Error("Create index error: " + err.Error())
		return
	}
	if !createIndex.Acknowledged {
		log.Error("IndicesCreateResult.Acknowledged %v", createIndex.Acknowledged)
		return
	}
	log.Info("Index is created")
}

func (tc TenantController) remove_tenant(ctx context.Context, index string, client *elastic.Client) {
	log.Info("Deleting index " + index)
	indexExists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Error("Check index exists error: " + err.Error())
		return
	}
	if indexExists {
		deleteIndex, err := client.DeleteIndex(index).Do(ctx)
		if err != nil {
			log.Error("Delete index error: " + err.Error())
			return
		}
		if !deleteIndex.Acknowledged {
			log.Error("Delete index error: " + err.Error())
			return
		}
	}
	log.Info("Index is deleted")
}

func (tc TenantController) Operate() {
	es_index := "incidents-" + tc.Tenant
	client, ctx, err := tc.es_client(es_index)
	if err != nil {
		log.Error("Can not connect to ES" + err.Error())
		return
	}
	if tc.Operation == "create" {
		tc.create_tenant(ctx, es_index, client)
	} else {
		tc.remove_tenant(ctx, es_index, client)
	}
}
