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
			logger.Errorf("Failed to watch user service: %v", err)
		}
	}()

	go func() {
		if err := p.resolver.Watch(ctx, p.cfg.Services.Group.Name, func(addrs []string) {
			p.reconnectGroupService(addrs)
		}); err != nil {
			logger.Errorf("Failed to watch group service: %v", err)
		}
	}()

	go func() {
		if err := p.resolver.Watch(ctx, p.cfg.Services.Message.Name, func(addrs []string) {
			p.reconnectMessageService(addrs)
		}); err != nil {
			logger.Errorf("Failed to watch message service: %v", err)
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
		logger.Errorf("Failed to connect to user service: %v", err)
		return
	}

	p.userConn = conn
	p.userClient = user.NewUserServiceClient(conn)
	logger.Infof("Connected to user service at %s", addr)
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
		logger.Errorf("Failed to connect to group service: %v", err)
		return
	}

	p.groupConn = conn
	p.groupClient = group.NewGroupServiceClient(conn)
	logger.Infof("Connected to group service at %s", addr)
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
		logger.Errorf("Failed to connect to message service: %v", err)
		return
	}

	p.messageConn = conn
	p.messageClient = message.NewMessageServiceClient(conn)
	logger.Infof("Connected to message service at %s", addr)
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
