package model

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

type Cluster struct {
	nodes   map[string]*Node
	lock    sync.Mutex
	redis   *redis.Client
	ipPool  []string
	usedIPs map[string]bool
}

func NewClusterManager(redisHost string, ipPool []string) *Cluster {
	client := redis.NewClient(&redis.Options{
		Addr: redisHost,
	})
	return &Cluster{
		nodes:   make(map[string]*Node),
		redis:   client,
		ipPool:  ipPool,
		usedIPs: make(map[string]bool),
	}
}

// AllocateIP assigns an IP address to a new node
func (cm *Cluster) AllocateIP() (string, error) {
	for _, ip := range cm.ipPool {
		if !cm.usedIPs[ip] {
			cm.usedIPs[ip] = true
			return ip, nil
		}
	}
	return "", fmt.Errorf("no available IPs")
}

// ReleaseIP releases an IP address back to the pool
func (cm *Cluster) ReleaseIP(ip string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.usedIPs, ip)
}
