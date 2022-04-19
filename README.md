# eth-rpc-reverse-proxy

## What I've done
1. Used fasthtp, zap - libraries which focused on performance, less allocations
2. HTTP and WS proxying
3. Removing backends if they are down
4. Round robin backend resolver
5. Retries for non-200 return code
   1. For http retries to one backend
   2. For WS reviving WS channel if it's down
6. Duplicating calls to another list of backends for `eth_call` method request
7. Prometheus + Grafana metrics
   1. Chainstack console-like dashboard

## What I didn't do, but can
1. Measure performance well
   1. I tried, but my old laptop shows RPS for clean FastHTTP like 2-3k RPS (but it can be 60k-200k rps easy) 
   2. And my reverse-proxy shows 2-2.5k RPS on my laptop, but if i did everything fine - proxy have to show RPS close to clean FastHTTP
2. Retries with changing backend if non-200 return code
3. ETCD with realtime updating list of backends
4. Even cleaner code (code would be better if I got at least one review)
5. More tests - controllers, clients


## How to run

1. Rename example config to `config.json`
2. `docker-compose up`

## Grafana

localhost:8080 -> log: `admin` pass: `admin`

![Metrics](https://i.imgur.com/ak88R4u.jpg)
