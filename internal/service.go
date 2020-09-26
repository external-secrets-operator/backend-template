package internal

import (
	"backend-template/generated/api"
	"context"
)

func NewBackendService() api.BackendService {
	return api.BackendService{GetSecrets: getSecrets}
}

func getSecrets(ctx context.Context, req *api.SecretsRequest) (*api.SecretsResponse, error) {
	secrets := []*api.SecretsResponse_Secret{{
		Key:    req.Keys[0],
		Values: [][]byte{[]byte("test")},
	}}
	return &api.SecretsResponse{Secrets: secrets}, nil
}
