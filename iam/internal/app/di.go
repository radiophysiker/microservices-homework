package app

import (
	"context"
	"fmt"
	"net"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	v1 "github.com/radiophysiker/microservices-homework/iam/internal/api/auth/v1"
	userapiv1 "github.com/radiophysiker/microservices-homework/iam/internal/api/user/v1"
	"github.com/radiophysiker/microservices-homework/iam/internal/config"
	"github.com/radiophysiker/microservices-homework/iam/internal/repository"
	"github.com/radiophysiker/microservices-homework/iam/internal/repository/session"
	userRepo "github.com/radiophysiker/microservices-homework/iam/internal/repository/user"
	"github.com/radiophysiker/microservices-homework/iam/internal/service"
	authSvc "github.com/radiophysiker/microservices-homework/iam/internal/service/auth"
	userSvc "github.com/radiophysiker/microservices-homework/iam/internal/service/user"
	"github.com/radiophysiker/microservices-homework/platform/pkg/cache"
	redisclient "github.com/radiophysiker/microservices-homework/platform/pkg/cache/redis"
	"github.com/radiophysiker/microservices-homework/platform/pkg/closer"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

// diContainer содержит все зависимости приложения с lazy initialization.
type diContainer struct {
	pool              *pgxpool.Pool
	redisPool         *redigo.Pool
	redisClient       cache.RedisClient
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	authService       service.AuthService
	userService       service.UserService
	authAPI           *v1.API
	userAPI           *userapiv1.API
}

// newDiContainer создает новый DI контейнер.
func newDiContainer() *diContainer {
	return &diContainer{}
}

// Pool возвращает пул соединений PostgreSQL с lazy initialization.
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

// RedisPool возвращает пул соединений Redis с lazy initialization.
func (d *diContainer) RedisPool(ctx context.Context) (*redigo.Pool, error) {
	if d.redisPool == nil {
		cfg := config.AppConfig()
		redisCfg := cfg.Redis

		address := net.JoinHostPort(redisCfg.Host(), redisCfg.Port())

		pool := &redigo.Pool{
			MaxIdle:     redisCfg.MaxIdle(),
			IdleTimeout: redisCfg.IdleTimeout(),
			Dial: func() (redigo.Conn, error) {
				conn, err := redigo.Dial("tcp", address)
				if err != nil {
					return nil, fmt.Errorf("dial redis: %w", err)
				}
				return conn, nil
			},
			TestOnBorrow: func(c redigo.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		}

		conn := pool.Get()
		defer func() {
			if err := conn.Close(); err != nil {
				logger.Error(ctx, "failed to close redis connection", zap.Error(err))
			}
		}()

		if _, err := conn.Do("PING"); err != nil {
			if err := pool.Close(); err != nil {
				logger.Error(ctx, "failed to close redis pool", zap.Error(err))
			}

			return nil, fmt.Errorf("ping redis: %w", err)
		}

		closer.AddNamed("Redis pool", func(ctx context.Context) error {
			return pool.Close()
		})

		d.redisPool = pool
	}

	return d.redisPool, nil
}

// RedisClient возвращает Redis клиент с lazy initialization.
func (d *diContainer) RedisClient(ctx context.Context) (cache.RedisClient, error) {
	if d.redisClient == nil {
		redisPool, err := d.RedisPool(ctx)
		if err != nil {
			return nil, err
		}

		d.redisClient = redisclient.NewClient(
			redisPool,
			logger.Logger(),
			config.AppConfig().Redis.ConnectionTimeout(),
		)
	}

	return d.redisClient, nil
}

// UserRepository возвращает репозиторий пользователей с lazy initialization.
func (d *diContainer) UserRepository(ctx context.Context) (repository.UserRepository, error) {
	if d.userRepository == nil {
		pool, err := d.Pool(ctx)
		if err != nil {
			return nil, err
		}

		d.userRepository = userRepo.NewRepository(pool)
	}

	return d.userRepository, nil
}

// SessionRepository возвращает репозиторий сессий с lazy initialization.
func (d *diContainer) SessionRepository(ctx context.Context) (repository.SessionRepository, error) {
	if d.sessionRepository == nil {
		redisClient, err := d.RedisClient(ctx)
		if err != nil {
			return nil, err
		}

		d.sessionRepository = session.NewRepository(redisClient, config.AppConfig().Session.TTL())
	}

	return d.sessionRepository, nil
}

// AuthService возвращает сервис аутентификации с lazy initialization.
func (d *diContainer) AuthService(ctx context.Context) (service.AuthService, error) {
	if d.authService == nil {
		userRepo, err := d.UserRepository(ctx)
		if err != nil {
			return nil, err
		}

		sessionRepo, err := d.SessionRepository(ctx)
		if err != nil {
			return nil, err
		}

		userSvc, err := d.UserService(ctx)
		if err != nil {
			return nil, err
		}

		d.authService = authSvc.NewService(
			userRepo,
			sessionRepo,
			userSvc,
			config.AppConfig().Session.TTL(),
		)
	}

	return d.authService, nil
}

// UserService возвращает сервис пользователей с lazy initialization.
func (d *diContainer) UserService(ctx context.Context) (service.UserService, error) {
	if d.userService == nil {
		userRepo, err := d.UserRepository(ctx)
		if err != nil {
			return nil, err
		}

		d.userService = userSvc.NewService(userRepo)
	}

	return d.userService, nil
}

// AuthAPI возвращает API слой для аутентификации с lazy initialization.
func (d *diContainer) AuthAPI(ctx context.Context) (*v1.API, error) {
	if d.authAPI == nil {
		authService, err := d.AuthService(ctx)
		if err != nil {
			return nil, err
		}

		d.authAPI = v1.NewAPI(authService)
	}

	return d.authAPI, nil
}

// UserAPI возвращает API слой для пользователей с lazy initialization.
func (d *diContainer) UserAPI(ctx context.Context) (*userapiv1.API, error) {
	if d.userAPI == nil {
		userService, err := d.UserService(ctx)
		if err != nil {
			return nil, err
		}

		d.userAPI = userapiv1.NewAPI(userService)
	}

	return d.userAPI, nil
}
