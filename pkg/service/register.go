package service

import "context"

var serviceLists []func(ctx context.Context, args []string) (Service, error)

func RegisterService(f func(ctx context.Context, args []string) (Service, error)) {
	serviceLists = append(serviceLists, f)
}

func ServiceList() []func(ctx context.Context, args []string) (Service, error) {
	return serviceLists
}
