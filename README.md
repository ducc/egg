# egg - the simple error aggregator
egg ingests errors and aggregates them

egg has 1st class support for sentry SDKs so you dont have to change any code to start using egg

## goals
1. simplicity - egg should only ingest errors and aggegate them
1. easy to install - unlike sentry, egg only requires one dependency (clickhouse)
1. easy to extend - egg should have a good api so additional features can be built on top of it
1. scalability - if clickhouse can scale, so can egg
1. compatible with sentry sdks

## usage
1. deploy egg's ingress and egress services
1. choose the extension services that you'd like (documented below)
1. start ingesting errors - see `examples` to learn how
1. use the [egg cli](https://github.com/ducc/egg/blob/master/cli/README.md)

## extensions
1. rest: forwards http rest requests to the grpc ingest endpoint and supports the sentry sdk `POST /api/{project_id}/store` api call
1. web: you can make it and pr it

## deploying on kubernetes
you can find instructions and manifests [here](https://github.com/ducc/egg/tree/master/.deploy)

## development
run `docker-compose up` to get egg running locally. when you save a file it will be restarted (think `npm start` hot reloading in react)

## tech used 
go, grpc, rest, clickhouse, sentry sdk, docker, kubernetes

## forking, cloning, stealing etc
it would make me happy if you give credit but do what you want, MIT license :)
