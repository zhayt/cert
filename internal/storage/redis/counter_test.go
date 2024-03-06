package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestCounterStorage_IncreaseCounter(t *testing.T) {
	// redis testcontainers
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		t.Error(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	type fields struct {
		client *redis.Client
	}

	field := fields{client: client}

	time.Sleep(2 * time.Second)

	type args struct {
		ctx context.Context
		key string
		val int64
	}

	arg1 := args{context.TODO(), "counter", 10}
	arg2 := args{context.TODO(), "counter", 1}
	arg3 := args{context.TODO(), "counter", -5}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantVal string
	}{
		{"success", field, arg1, false, "10"},
		{"success", field, arg2, false, "11"},
		{"success", field, arg3, false, "6"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CounterStorage{
				client: tt.fields.client,
			}
			if err := r.IncreaseCounter(tt.args.ctx, tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("IncreaseCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if val, err := tt.fields.client.Get(context.TODO(), "counter").Result(); err != nil || val != tt.wantVal {
				t.Errorf("IncreaseCounter() val = %v, wantVal %v", val, tt.wantVal)
			}
		})
	}
}

func TestCounterStorage_DecreaseCounter(t *testing.T) {
	// redis testcontainers
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		t.Error(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	type fields struct {
		client *redis.Client
	}

	field := fields{client: client}

	time.Sleep(2 * time.Second)
	field.client.Set(context.TODO(), "counter", 20, 0)

	type args struct {
		ctx context.Context
		key string
		val int64
	}

	arg1 := args{context.TODO(), "counter", 10}
	arg2 := args{context.TODO(), "counter", 1}
	arg3 := args{context.TODO(), "counter", -5}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantVal string
	}{
		{"success", field, arg1, false, "10"},
		{"success", field, arg2, false, "9"},
		{"success", field, arg3, false, "14"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CounterStorage{
				client: tt.fields.client,
			}
			if err := r.DecreaseCounter(tt.args.ctx, tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("DecreaseCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if val, err := tt.fields.client.Get(context.TODO(), "counter").Result(); err != nil || val != tt.wantVal {
				t.Errorf("DecreaseCounter() val = %v, wantVal %v", val, tt.wantVal)
			}
		})
	}
}

func TestCounterStorage_GetCounter(t *testing.T) {
	// redis testcontainers
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		t.Error(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	type fields struct {
		client *redis.Client
	}

	field := fields{client: client}

	time.Sleep(2 * time.Second)
	field.client.Set(context.TODO(), "counter", 20, 0)

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"success", field, args{context.TODO(), "counter"}, "20", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CounterStorage{
				client: tt.fields.client,
			}
			got, err := r.GetCounter(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCounter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
