# Json file prometheus exporter

**This project is a proof of concept, should not be used in any production scenarios!**

json-file-prometheus-exporter is a command line application,
that is able to read files with [JSON content](https://www.json.org/),
and export them to [prometheus](https://prometheus.io/).

The most common use-case is [ndjson](http://ndjson.org/),
where the json values are separated by a new line (`\n`).

The application reads the whole file,
and waits for it to be written again.

## Example

Input `/tmp/analytics.json`

```json
{"analytics":{"total": 150,"mode": "PER_NODE"},"from":"192.168.1.2"}
{"analytics":{"total": 120,"mode": "PER_NODE"},"from":"192.168.1.3"}
{"analytics":{"total": 130,"mode": "PER_NODE"},"from":"192.168.1.2"}
{"analytics":{"total": 110,"mode": "PER_NODE"},"from":"192.168.1.3"}
```

Command

```sh
./json-file-prometheus-exporter --metric-name test_metrics --labels "analytics.total,analytics.mode,from" /tmp/analytics.json
```

Metrics on `http://localhost:8080/metrics`

```text
test_metrics{analytics_total="150",analytics_mode="PER_NODE",from="192.168.1.2"} 1
test_metrics{analytics_total="120",analytics_mode="PER_NODE",from="192.168.1.3"} 1
test_metrics{analytics_total="130",analytics_mode="PER_NODE",from="192.168.1.2"} 1
test_metrics{analytics_total="110",analytics_mode="PER_NODE",from="192.168.1.3"} 1
```

You can use the following command to explore usage and required parameters:

```sh
./json-file-prometheus-exporter -h
```

## Filebeat example usage

[Filebeat](https://www.elastic.co/guide/en/beats/filebeat/current/filebeat-overview.html) has a number of
[modules](https://www.elastic.co/guide/en/beats/filebeat/current/configuration-filebeat-modules.html) that can parse log files,
and export the parsed logs in JSON format to another file.

Once the new file is written by filebeat,
this application can read it,
and export it to prometheus.

Partial filebeat config `/etc/filebeat/filebeat.yml`:

```yml
# all configs above

output.file:
  path: "/tmp/filebeat"
  filename: filebeat
```

You can start the application by:

```sh
./json-file-prometheus-exporter --metric-name test_metrics --labels "analytics.total,analytics.mode,from" /tmp/filebeat/filebeat-20221031.ndjson
```
