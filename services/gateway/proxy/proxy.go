// Package proxy 提供 Gateway 到后端各服务的连接代理，负责服务地址发现与自动重连。
package proxy

import (
	"IM/api/gen/group"
	"IM/api/gen/message"
	"IM/api/gen/user"
	"IM/pkg/config"
	"IM/pkg/discovery"
	"IM/pkg/logger"
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceProxy 管理与 User/Group/Message 等后端服务的 gRPC 连接与客户端实例。
type ServiceProxy struct {
	resolver      *discovery.Resolver
	cfg           *config.GatewayConfig
	userClient    user.UserServiceClient
	groupClient   group.GroupServiceClient
	messageClient message.MessageServiceClient
	userConn      *grpc.ClientConn
	groupConn     *grpc.ClientConn
	messageConn   *grpc.ClientConn
	mu            sync.RWMutex
}

// NewServiceProxy 创建 ServiceProxy，并在后台启动服务观察与自动重连逻辑。
func NewServiceProxy(resolver *discovery.Resolver, cfg *config.GatewayConfig) *ServiceProxy {
	p := &ServiceProxy{
		resolver: resolver,
		cfg:      cfg,
	}

	go p.watchServices()

	return p
}

func (p *ServiceProxy) watchServices() {
	ctx := context.Background()

	go func() {
		if err := p.resolver.Watch(ctx, p.cfg.Services.User.Name, func(addrs []string) {
			p.reconnectUserService(addrs)
		}); err != nil {
			logger.Errorw("Failed to watch user service", "component", "gateway_proxy", "err", err)
		}
	}()

	go func() {
		if err := p.resolver.Watch(ctx, p.cfg.Services.Group.Name, func(addrs []string) {
			p.reconnectGroupService(addrs)
		}); err != nil {
			logger.Errorw("Failed to watch group service", "component", "gateway_proxy", "err", err)
		}
	}()

	go func() {
		if err := p.resolver.Watch(ctx, p.cfg.Services.Message.Name, func(addrs []string) {
			p.reconnectMessageService(addrs)
		}); err != nil {
			logger.Errorw("Failed to watch message service", "component", "gateway_proxy", "err", err)
		}
	}()
}

func (p *ServiceProxy) reconnectUserService(addrs []string) {
	if len(addrs) == 0 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.userConn != nil {
		p.userConn.Close()
	}

	addr := addrs[0]
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorw("Failed to connect to user service", "component", "gateway_proxy", "err", err, "addr", addr)
		return
	}

	p.userConn = conn
	p.userClient = user.NewUserServiceClient(conn)
	logger.Infow("Connected to user service", "component", "gateway_proxy", "addr", addr)
}

func (p *ServiceProxy) reconnectGroupService(addrs []string) {
	if len(addrs) == 0 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.groupConn != nil {
		p.groupConn.Close()
	}

	addr := addrs[0]
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorw("Failed to connect to group service", "component", "gateway_proxy", "err", err, "addr", addr)
		return
	}

	p.groupConn = conn
	p.groupClient = group.NewGroupServiceClient(conn)
	logger.Infow("Connected to group service", "component", "gateway_proxy", "addr", addr)
}

func (p *ServiceProxy) reconnectMessageService(addrs []string) {
	if len(addrs) == 0 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.messageConn != nil {
		p.messageConn.Close()
	}

	addr := addrs[0]
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorw("Failed to connect to message service", "component", "gateway_proxy", "err", err, "addr", addr)
		return
	}

	p.messageConn = conn
	p.messageClient = message.NewMessageServiceClient(conn)
	logger.Infow("Connected to message service", "component", "gateway_proxy", "addr", addr)
}

func (p *ServiceProxy) UserClient() user.UserServiceClient {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.userClient
}

func (p *ServiceProxy) GroupClient() group.GroupServiceClient {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.groupClient
}

func (p *ServiceProxy) MessageClient() message.MessageServiceClient {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.messageClient
}

func (p *ServiceProxy) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.userConn != nil {
		p.userConn.Close()
	}
	if p.groupConn != nil {
		p.groupConn.Close()
	}
	if p.messageConn != nil {
		p.messageConn.Close()
	}
}

func (p *ServiceProxy) DiscoverService(ctx context.Context, serviceName string) ([]string, error) {
	return p.resolver.Discover(ctx, serviceName)
}

// GetServiceAddr 返回指定服务的一个可用实例地址（超时 5s）。
func (p *ServiceProxy) GetServiceAddr(serviceName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addrs, err := p.resolver.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}

	if len(addrs) == 0 {
		return "", fmt.Errorf("no instances found for service %s", serviceName)
	}

	return addrs[0], nil
}
