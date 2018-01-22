# elasticappproxy


The server listens on port 9090.

It expects elasticsearch to be listening on localhost:9200, so start that first
You can run elasticsearch with docker.

```bash
docker run -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:6.1.2
```

Start the server

```bash
go get
go build && ./elasticappproxy
```

You can hook up kibana to elasticsearch too. It listens on localhost:5601

```bash
docker run -e ELASTICSEARCH_URL=http://172.17.0.1:9200 -p 5601:5601 docker.elastic.co/kibana/kibana:6.1.2
```

When you run the App from https://github.com/david-martin/elasticsearch-mobile-poc (configured to talk to the running elastic app proxy), you should then be able to create an index `metric*`, using the `date` field as a time filter, in Kibana and see at least 1 entry for the app metric on the Discover tab.
You can also Visualize the data using charts.

![Kibana](/kibana.png)

You can push some dummy data to the server too using curl e.g.

```bash
curl -H "ClientID:aaa" -H "X-User-Agent:Dalvik/2.1.0 (Linux; U; Android 4.4; TEST)" 10.201.82.209:9090
curl -H "ClientID:bbb" -H "X-User-Agent:Dalvik/2.1.0 (Linux; U; Android 4.4; TEST)" 10.201.82.209:9090
curl -H "ClientID:ccc" -H "X-User-Agent:Dalvik/2.1.0 (Linux; U; Android 5.0; TEST)" 10.201.82.209:9090
curl -H "ClientID:ddd" -H "X-User-Agent:Dalvik/2.1.0 (Linux; U; Android 5.0; TEST)" 10.201.82.209:9090
# Deliberate duplicate calls to simulate multiple inits
curl -H "ClientID:eee" -H "X-User-Agent:Dalvik/2.1.0 (Linux; U; Android 5.0; TEST)" 10.201.82.209:9090
curl -H "ClientID:eee" -H "X-User-Agent:Dalvik/2.1.0 (Linux; U; Android 5.0; TEST)" 10.201.82.209:9090
curl -H "ClientID:eee" -H "X-User-Agent:Dalvik/2.1.0 (Linux; U; Android 5.0; TEST)" 10.201.82.209:9090
curl -H "ClientID:eee" -H "X-User-Agent:Dalvik/2.1.0 (Linux; U; Android 5.0; TEST)" 10.201.82.209:9090
curl -H "ClientID:eee" -H "X-User-Agent:Dalvik/2.1.0 (Linux; U; Android 5.0; TEST)" 10.201.82.209:9090
```

With this extra data you can visualize the distribution of OS versions across unique clients with a pie chart.
The below chart shows 2 users with Android 5.0 (ccc, ddd & eee), 2 with 4.4 (aaa & bbb) 1 with 5.1 (actual device)

![Kibana Pie Chart](/kibana_pie_chart.png)
