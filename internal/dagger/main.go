package main

import (
	"context"
	"dagger/dagger-hello/internal/dagger"
)

type DaggerHello struct{}

// Build a ready-to-use development environment
func (m *DaggerHello) BuildEnv(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:1.22").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache")
}

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
