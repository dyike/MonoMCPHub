package server

import (
	"context"

	"github.com/dyike/MonoMCPHub/pkg/service"
	mcp_server "github.com/mark3labs/mcp-go/server"
)

type HubServer struct {
	ctx    context.Context
	server *mcp_server.MCPServer

	services []service.Service
}

func NewHubServer(ctx context.Context, serverName string, srvs []service.Service) (*HubServer, error) {
	mcpServer := mcp_server.NewMCPServer(
		serverName,
		"0.0.1",
		mcp_server.WithResourceCapabilities(true, true),
		mcp_server.WithPromptCapabilities(true),
		mcp_server.WithToolCapabilities(true),
	)

	hs := &HubServer{
		ctx:      ctx,
		server:   mcpServer,
		services: srvs,
	}

	err := hs.init()

	return hs, err
}

func (hs *HubServer) init() error {
	var err error

	for _, srv := range hs.services {
		hs.loadService(srv)
	}

	return err
}

func (hs *HubServer) loadService(srv service.Service) error {
	for r, rhf := range srv.Resources() {
		hs.server.AddResource(r, rhf)
	}

	for rt, rtf := range srv.ResourceTemplates() {
		hs.server.AddResourceTemplate(rt, rtf)
	}

	hs.server.AddTools(srv.Tools()...)

	for n, nhf := range srv.NotificationHandlers() {
		hs.server.AddNotificationHandler(n, nhf)
	}

	for _, pe := range srv.Prompts() {
		hs.server.AddPrompt(pe.Prompt(), pe.PromptHandlerFunc())
	}
	return nil
}

func (hs *HubServer) Serve() error {
	return mcp_server.ServeStdio(hs.server)
}
