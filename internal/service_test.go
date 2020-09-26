package internal

import (
	"backend-template/generated/api"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_should_get_secrets(t *testing.T) {
	s := NewBackendService()

	r, err := s.GetSecrets(context.TODO(), &api.SecretsRequest{Keys: []string{"test"}})

	assert.NoError(t, err)
	assert.Equal(t, &api.SecretsResponse{
		Secrets: []*api.SecretsResponse_Secret{{
			Key:    "test",
			Values: [][]byte{[]byte("test")},
		}},
	}, r)
}
