package app

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/radiophysiker/microservices-homework/order/internal/api/order/v1"
	clientGrpc "github.com/radiophysiker/microservices-homework/order/internal/client/grpc"
	inventoryClient "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/payment/v1"
	"github.com/radiophysiker/microservices-homework/order/internal/config"
	"github.com/radiophysiker/microservices-homework/order/internal/repository"
	orderRepo "github.com/radiophysiker/microservices-homework/order/internal/repository/order"
	"github.com/radiophysiker/microservices-homework/order/internal/service"
	orderSvc "github.com/radiophysiker/microservices-homework/order/internal/service/order"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	inventorypb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	pool            *pgxpool.Pool
	inventoryConn   *grpc.ClientConn
	paymentConn     *grpc.ClientConn
	orderRepository repository.OrderRepository
	inventoryClient clientGrpc.InventoryClient
	paymentClient   clientGrpc.PaymentClient
	orderService    service.OrderService
	api             *apiv1.API
}

func newDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Pool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.pool == nil {
		pc, err := pgxpool.ParseConfig(config.AppConfig().Postgres.DSN())
		if err != nil {
			return nil, fmt.Errorf("parse postgres config: %w", err)
		}

		pc.MaxConns = config.AppConfig().Postgres.PoolMaxConns()
		pc.MinConns = config.AppConfig().Postgres.PoolMinConns()
		pc.MaxConnLifetime = config.AppConfig().Postgres.PoolMaxConnLifetime()
		pc.MaxConnIdleTime = config.AppConfig().Postgres.PoolMaxConnIdleTime()

		ctxConnect, cancelConnect := context.WithTimeout(ctx, 10*time.Second)
		defer cancelConnect()

		pool, err := pgxpool.NewWithConfig(ctxConnect, pc)
		if err != nil {
			return nil, fmt.Errorf("create postgres pool: %w", err)
		}

		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			return nil, fmt.Errorf("ping postgres: %w", err)
		}

		closer.AddNamed("PostgreSQL pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pool = pool
	}

	return d.pool, nil
}

func (d *diContainer) InventoryConn(ctx context.Context) (*grpc.ClientConn, error) {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.InventoryAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("connect inventory grpc: %w", err)
		}

		closer.AddNamed("inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn, nil
}

func (d *diContainer) PaymentConn(ctx context.Context) (*grpc.ClientConn, error) {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.PaymentAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("connect payment grpc: %w", err)
		}

		closer.AddNamed("payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn, nil
}

func (d *diContainer) OrderRepository(ctx context.Context) (repository.OrderRepository, error) {
	if d.orderRepository == nil {
		pool, err := d.Pool(ctx)
		if err != nil {
			return nil, err
		}

		d.orderRepository = orderRepo.NewRepository(pool)
	}

	return d.orderRepository, nil
}

func (d *diContainer) InventoryClient(ctx context.Context) (clientGrpc.InventoryClient, error) {
	if d.inventoryClient == nil {
		conn, err := d.InventoryConn(ctx)
		if err != nil {
			return nil, err
		}

		d.inventoryClient = inventoryClient.NewClient(
			inventorypb.NewInventoryServiceClient(conn),
		)
	}

	return d.inventoryClient, nil
}

func (d *diContainer) PaymentClient(ctx context.Context) (clientGrpc.PaymentClient, error) {
	if d.paymentClient == nil {
		conn, err := d.PaymentConn(ctx)
		if err != nil {
			return nil, err
		}

		d.paymentClient = paymentClient.NewClient(
			paymentpb.NewPaymentServiceClient(conn),
		)
	}

	return d.paymentClient, nil
}

func (d *diContainer) OrderService(ctx context.Context) (service.OrderService, error) {
	if d.orderService == nil {
		orderRepo, err := d.OrderRepository(ctx)
		if err != nil {
			return nil, err
		}

		inventoryClient, err := d.InventoryClient(ctx)
		if err != nil {
			return nil, err
		}

		paymentClient, err := d.PaymentClient(ctx)
		if err != nil {
			return nil, err
		}

		d.orderService = orderSvc.NewService(
			orderRepo,
			inventoryClient,
			paymentClient,
		)
	}

	return d.orderService, nil
}

func (d *diContainer) API(ctx context.Context) (*apiv1.API, error) {
	if d.api == nil {
		orderService, err := d.OrderService(ctx)
		if err != nil {
			return nil, err
		}

		d.api = apiv1.NewAPI(orderService)
	}

	return d.api, nil
}
