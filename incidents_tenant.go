package main

import (
	"flag"
	"github.com/leizhu/incidents_tenant/controller"
	"log"
)

var (
	operation = flag.String(
		"operation", "",
		"create or remove tenant",
	)
	tenant_name = flag.String(
		"tenant_name", "",
		"tenant name",
	)
	tenant_type = flag.String(
		"tenant_type", "vip",
		"tenant type",
	)
	number_of_shards = flag.Int(
		"number_of_shards", 3,
		"The number of index shards",
	)
	number_of_replicas = flag.Int(
		"number_of_replicas", 1,
		"The number of index replica shards",
	)
	es_url = flag.String(
		"elasticsearch.url", "http://elasticsearch:9200",
		"URL of elasticsearch",
	)
	log_level = flag.String(
		"loglevel", "INFO",
		"log level",
	)
)

func main() {
	flag.Parse()

	if *operation == "" || (*operation != "create" && *operation != "remove") {
		log.Println("Please specify operation type, create or remove")
		return
	}
	if *tenant_name == "" {
		log.Println("Please specify tenant name")
		return
	}
	log.Println("elasticsearch.url: ", *es_url)
	log.Println("operation: ", *operation)
	log.Println("tenant_name: ", *tenant_name)
	log.Println("tenant_type: ", *tenant_type)
	log.Println("number_of_shards: ", *number_of_shards)
	log.Println("number_of_replicas: ", *number_of_replicas)
	log.Println("log level: ", *log_level)
	controller.InitLog(*log_level)
	tc := controller.NewTenantController(*es_url, *operation, *tenant_name, *number_of_shards, *number_of_replicas)
	tc.ConfigIndex()
	tc.Operate()
}
