package internal

import (
	"context"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestApp_should_start_http_server(t *testing.T) {
	app, err := prepareApp()
	assert.NoError(t, err)
	assert.NoError(t, app.Start())

	r, err := http.Get("http://:" + os.Getenv(appName+"_HTTP_PORT") + "/health")
	assert.NoError(t, err)
	//noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{"status": "SERVING"}`, string(body))

	assert.NoError(t, app.Stop())
}

func TestApp_should_start_grpc_server(t *testing.T) {
	app, err := prepareApp()
	assert.NoError(t, err)
	assert.NoError(t, app.Start())

	conn, err := grpc.Dial(":"+os.Getenv(appName+"_GRPC_PORT"), grpc.WithInsecure(), grpc.WithBlock())
	assert.NoError(t, err)

	c := grpc_health_v1.NewHealthClient(conn)

	r, err := c.Check(context.TODO(), &grpc_health_v1.HealthCheckRequest{})
	assert.NoError(t, err)

	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, r.Status)

	assert.NoError(t, app.Stop())
}

//noinspection GoUnhandledErrorResult
func prepareApp() (App, error) {
	httpPort, _ := freeport.GetFreePort()
	os.Setenv(appName+"_HTTP_PORT", strconv.Itoa(httpPort))
	grpcPort, _ := freeport.GetFreePort()
	os.Setenv(appName+"_GRPC_PORT", strconv.Itoa(grpcPort))
	return NewApp()
}
