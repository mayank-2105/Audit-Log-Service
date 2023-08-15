package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic/v7"
)

// Event represents an audit log event. There are some common fields like ID, Timestamp
type Event struct {
	ID        string                 `json:"id"` // this will serve as a unique identifier for the event
	Timestamp time.Time              `json:"timestamp"` //timestamp of the event, I set this to time.Now by default whenever a request is received
	Type      string                 `json:"type"` //type will basically correspond to log level like ERROR, DEBUG or INFO 
	Action    string                 `json:"action"` //this will correspond to actions like account_deleted, account_created etc
	Identity  string                 `json:"identity"` // this specifies the identity of the user which performed the action
	Data      map[string]interface{} `json:"data"`  /* you can send event specific fields as json and they will be stored in data field, this is map where keys are strings and values are denoted by interface type which can correspond to any data type*/ 
}

/* There were multiple options available for data storage here like ElasticSearch, Cassandra, MongoDB etc but I proceeded with ElasticSearch for the following reasons:-

1. logs are write-intensive. ElasticSearch offers good horizontal scaling support

2. there can be some data which is specific to the event and not defined in advance. 
We may also need to perform queries on such data. Thus ElasticSearch with its mechanism
to store data in JSON documents seems to be the best choice here.

3. possible future integration with Logstash and Kibana for easy data ingestion and Analytics
*/

var esClient *elastic.Client
var indexName = "audit_log" // this variable specifies name of index where the document(event is stored)


/* This func will run a new elasticsearch client on http://localhost:9200*/
func initElasticsearch() {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %v", err)
	}

	esClient = client
}

/*the storeEvent function takes an Event object, 
indexes it in the specified Elasticsearch index with the provided event ID
and returns any errors that occurred during the indexing process. 
This function is responsible for persisting event data in Elasticsearch, 
which will allow us to later search, retrieve, and analyze the stored events.

Below is a documentation of what is happening overall via each method

ctx := context.Background(): This line creates a new background context, which is used for executing Elasticsearch operations.

esClient.Index(): This is a method call on the Elasticsearch client that prepares an index operation.

.Index(indexName): This sets the name of the index where the document (event) will be stored. The indexName is likely a variable that holds the name of the index.

.Id(event.ID): This sets the document ID for the indexed event. The event.ID is assumed to be a unique identifier for the event.

.BodyJson(event): This sets the JSON representation of the Event object as the body of the index operation. Elasticsearch will store this JSON data as the document in the specified index.

.Do(ctx): This executes the index operation with the provided context. It sends the index request to the Elasticsearch cluster.
*/
func storeEvent(event Event) error {
	ctx := context.Background()
	_, err := esClient.Index().
		Index(indexName). 
		Id(event.ID). 
		BodyJson(event).
		Do(ctx)

	return err
}

/* Note- For simplicity, I have not placed a strict validation check on the request based on field values
, however we may need this in real world.
 
This func is mainly used as an event handler, It is 
mainly called whenever a user makes a POST request to /event
in order to basically log an event
It creates a new event ID based on current timestamp 
and calls the storeevent method to store the event in elasticsearch*/
func handleEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Bad Request, Empty event body", http.StatusBadRequest)
		return
	}

	event.ID = fmt.Sprintf("%d", time.Now().UnixNano())  // we are setting the ID of the log based on the current time, this can be changed to a UUID based implementation as well
	event.Timestamp = time.Now()

	err = storeEvent(event)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

/*Note- For simplicity, I have not placed a strict validation check on the request based on field values
, however we may need this in real world.

This function is used when we want to query both common and event specific data. 
It extracts the query params from request URL and then queries on the esClient based on them" 
Finally it results all the events which match the query params*/
func handleQuery(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	termQuery := elastic.NewBoolQuery() // create an empty query
	for key, value := range queryParams {
		if key == "data" {
			// Construct a nested query for data field
			dataQuery := elastic.NewNestedQuery("data",
				elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("data."+value[0], true),
				),
			)
			termQuery.Must(dataQuery)
		} else {
			termQuery.Must(elastic.NewMatchQuery(key, value[0]))
		}
	}

	searchResult, err := esClient.Search().
		Index(indexName).
		Query(termQuery).
		Do(context.Background())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var events []Event
	for _, hit := range searchResult.Hits.Hits {
		var event Event
		err := json.Unmarshal(hit.Source, &event)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		events = append(events, event)
	}

	jsonEvents, err := json.Marshal(events)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonEvents)
}

func main() {
	initElasticsearch()

	r := mux.NewRouter()
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/event", TokenAuthMiddleware(handleEvent)).Methods("POST")
	r.HandleFunc("/query", TokenAuthMiddleware(handleQuery)).Methods("GET")

	log.Println("Audit Log Service is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
