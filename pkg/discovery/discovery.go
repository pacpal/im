package discovery

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Registry 管理服务注册
type Registry struct {
	client     *clientv3.Client
	leaseID    clientv3.LeaseID
	cancelFunc context.CancelFunc
}

// NewRegistry 创建注册器
func NewRegistry(endpoints []string, dialTimeout time.Duration) (*Registry, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}
	return &Registry{
		client: client,
	}, nil
}

// Register 注册服务，带 TTL 和自动续约
func (r *Registry) Register(ctx context.Context, serviceName, addr string, ttl int64) error {
	// 创建租约
	leaseResp, err := r.client.Grant(ctx, ttl)
	if err != nil {
		return fmt.Errorf("failed to create lease: %w", err)
	}
	r.leaseID = leaseResp.ID

	// 注册服务
	key := fmt.Sprintf("/services/%s/%s", serviceName, addr)
	_, err = r.client.Put(ctx, key, addr, clientv3.WithLease(r.leaseID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	// 自动续约
	keepAliveCtx, cancel := context.WithCancel(ctx)
	r.cancelFunc = cancel

	keepAliveCh, err := r.client.KeepAlive(keepAliveCtx, r.leaseID)
	if err != nil {
		cancel()
		return fmt.Errorf("failed to keep alive: %w", err)
	}

	// 后台消费 keepAlive 响应
	go func() {
		for {
			select {
			case <-keepAliveCtx.Done():
				return
			case _, ok := <-keepAliveCh:
				if !ok {
					return
				}
			}
		}
	}()

	return nil
}

// Deregister 注销服务
func (r *Registry) Deregister(ctx context.Context, serviceName, addr string) error {
	key := fmt.Sprintf("/services/%s/%s", serviceName, addr)
	_, err := r.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	// 取消续约
	if r.cancelFunc != nil {
		r.cancelFunc()
	}

	// 撤销租约
	if r.leaseID != 0 {
		_, err = r.client.Revoke(ctx, r.leaseID)
		if err != nil {
			return fmt.Errorf("failed to revoke lease: %w", err)
		}
	}

	return nil
}

// Close 关闭连接
func (r *Registry) Close() error {
	if r.cancelFunc != nil {
		r.cancelFunc()
	}
	return r.client.Close()
}

// Resolver 服务发现
type Resolver struct {
	client *clientv3.Client
}

// NewResolver 创建服务发现解析器
func NewResolver(endpoints []string, dialTimeout time.Duration) (*Resolver, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}
	return &Resolver{client: client}, nil
}

// Discover 获取服务地址列表
func (r *Resolver) Discover(ctx context.Context, serviceName string) ([]string, error) {
	prefix := fmt.Sprintf("/services/%s/", serviceName)
	resp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to discover services: %w", err)
	}

	addrs := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		addrs = append(addrs, string(kv.Value))
	}
	return addrs, nil
}

// Watch 监听服务变化
func (r *Resolver) Watch(ctx context.Context, serviceName string, callback func([]string)) error {
	prefix := fmt.Sprintf("/services/%s/", serviceName)

	// 先获取当前列表
	addrs, err := r.Discover(ctx, serviceName)
	if err != nil {
		return err
	}
	callback(addrs)

	// 开始监听变化
	watchCh := r.client.Watch(ctx, prefix, clientv3.WithPrefix())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-watchCh:
				if !ok {
					return
				}
				// 有变化时重新获取完整列表
				currentAddrs, err := r.Discover(ctx, serviceName)
				if err == nil {
					callback(currentAddrs)
				}
			}
		}
	}()

	return nil
}

// Close 关闭连接
func (r *Resolver) Close() error {
	return r.client.Close()
}
