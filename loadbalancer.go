package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	servers []string
	current int
	mutex   sync.Mutex
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
		current: 0,
	}
}

func (lb *LoadBalancer) GetNextServer() string {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	server := lb.servers[lb.current]
	lb.current = (lb.current + 1) % len(lb.servers)
	return server
}

func (lb *LoadBalancer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	server := lb.GetNextServer()
	fmt.Printf("🔄 Routing to: %s\n", server)

	serverURL, err := url.Parse(server)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(serverURL)
	proxy.ServeHTTP(w, r)
}

func main() {
	servers := []string{
		"http://localhost:8001",
		"http://localhost:8002",
	}

	lb := NewLoadBalancer(servers)

	http.HandleFunc("/", lb.HandleRequest)

	fmt.Println("🚀 Load Balancer running on :8080")
	fmt.Println("📡 Balancing between servers:", servers)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

