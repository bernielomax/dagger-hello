package main

import (
	"context"
	"dagger/dagger-hello/internal/dagger"
	"fmt"
	"math"
	"math/rand/v2"
)

type DaggerHello struct{}

// Build a ready-to-use development environment
func (m *DaggerHello) BuildEnv(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:1.22.5-alpine3.20").
		WithDirectory("/src", source).
		WithDirectory("/dist", dagger.Connect().Directory()).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache")
}

// HttpService to run during E2E tests.
func (m *DaggerHello) HttpService() *dagger.Service {
	return dag.Container().
		From("python").
		WithWorkdir("/srv").
		WithNewFile("index.html", "Hello, world!").
		WithExec([]string{"python", "-m", "http.server", "8080"}).
		WithExposedPort(8080).
		AsService()
}

// Return the result of running unit tests
func (m *DaggerHello) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BuildEnv(source).
		WithServiceBinding("www", m.HttpService()).
		WithExec([]string{"go", "test", "-v", "./..."}).
		Stdout(ctx)
}

// Build the Hello Dagger application.
func (m *DaggerHello) Build(source *dagger.Directory) *dagger.Container {
	build := m.BuildEnv(source).
		WithExec([]string{"go", "build", "-o", "/dist/hello", "."}).
		Directory("/dist")

	return dag.Container().
		From("alpine:3.20").
		WithDirectory("/app", build).
		WithExec([]string{"adduser", "-D", "hello"}).
		WithEntrypoint([]string{"/app/hello"}).
		WithExec([]string{"chown", "-R", "hello:", "/app"}).
		WithUser("hello")
}

// Publish the application container after building and testing it on-the-fly.
func (m *DaggerHello) Publish(ctx context.Context, source *dagger.Directory) (string, error) {
	_, err := m.Test(ctx, source)
	if err != nil {
		return "", err
	}

	address, err := m.
		Build(source).
		Publish(ctx, fmt.Sprintf("bernielomax/dagger-hello:%.0f", math.Floor(rand.Float64()*10000000)))
	if err != nil {
		return "", err
	}

	return address, nil
}
