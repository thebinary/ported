package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var red *redis.Client

func keysForServiceName(serviceName string) (routerKey, serviceKey, loadBalancerKey string) {
	routerKey = fmt.Sprintf("%s/http/routers/%s", rootKey, serviceName)
	serviceKey = fmt.Sprintf("%s/http/routers/%s/service", rootKey, serviceName)
	loadBalancerKey = fmt.Sprintf("%s/http/services/%s/loadbalancer", rootKey, serviceName)
	return routerKey, serviceKey, loadBalancerKey
}

func createPortedService(ctx *context.Context, red *redis.Client, serviceName, addr string) (serviceURL string, err error) {
	//TODO: handle errors
	serviceURL = "http://" + serviceName + "." + rootDomain
	routerKey, serviceKey, loadBalancerKey := keysForServiceName(serviceName)
	log.Println("[KEY]", red.Set(*ctx, loadBalancerKey+"/servers/0/url", addr, keepAliveTimeout+time.Second*2))
	log.Println("[KEY]", red.Set(*ctx, serviceKey, serviceName, keepAliveTimeout+time.Second*4))
	log.Println("[KEY]", red.Set(*ctx, routerKey+"/rule", fmt.Sprintf("Host(`%s.%s`)", serviceName, rootDomain), keepAliveTimeout))
	return serviceURL, nil
}

func updateKeepAlive(ctx *context.Context, red *redis.Client, serviceName string) (err error) {
	//TODO: handle errors
	routerKey, serviceKey, loadBalancerKey := keysForServiceName(serviceName)

	log.Println("[KEY]", red.Expire(*ctx, loadBalancerKey, keepAliveTimeout+time.Second*2))
	log.Println("[KEY]", red.Expire(*ctx, serviceKey, keepAliveTimeout+time.Second*4))
	log.Println("[KEY]", red.Expire(*ctx, routerKey+"/rule", keepAliveTimeout))
	return nil
}
