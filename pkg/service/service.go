package service

import (
	"context"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Service interface {
	Ctx() context.Context

	// Tools returns the tools of the service
	Tools() []server.ServerTool

	// Prompts returns the prompts of the service
	Prompts() []PromptEntry

	// Resources returns the resources of the service
	Resources() map[mcp.Resource]server.ResourceHandlerFunc

	// ResourceTemplates returns the resource templates of the service
	ResourceTemplates() map[mcp.ResourceTemplate]server.ResourceTemplateHandlerFunc

	// NotificationHandlers returns the notification handlers of the service
	NotificationHandlers() map[string]server.NotificationHandlerFunc

	// Config returns the config of the service
	Config() string

	// Name returns the name of the service
	Name() string

	// Close closes the service and releases all resources
	Close() error
}

type PromptEntry struct {
	prompt mcp.Prompt
	phf    server.PromptHandlerFunc
}

func (p *PromptEntry) Prompt() mcp.Prompt {
	return p.prompt
}

func (p *PromptEntry) PromptHandlerFunc() server.PromptHandlerFunc {
	return p.phf
}

type ServiceManager struct {
	ctx                  context.Context
	lock                 *sync.Mutex
	tools                []server.ServerTool
	prompts              []PromptEntry
	notificationHandlers map[string]server.NotificationHandlerFunc
	reources             map[mcp.Resource]server.ResourceHandlerFunc
	resourceTemplates    map[mcp.ResourceTemplate]server.ResourceTemplateHandlerFunc
}

func NewServiceManager(ctx context.Context) *ServiceManager {
	return &ServiceManager{
		ctx:                  ctx,
		lock:                 &sync.Mutex{},
		tools:                make([]server.ServerTool, 0),
		prompts:              make([]PromptEntry, 0),
		notificationHandlers: make(map[string]server.NotificationHandlerFunc),
		reources:             make(map[mcp.Resource]server.ResourceHandlerFunc),
		resourceTemplates:    make(map[mcp.ResourceTemplate]server.ResourceTemplateHandlerFunc),
	}
}

func (sm *ServiceManager) AddTool(tool mcp.Tool, handler server.ToolHandlerFunc) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.tools = append(sm.tools, server.ServerTool{Tool: tool, Handler: handler})
}

func (sm *ServiceManager) AddPrompt(prompt mcp.Prompt, phf server.PromptHandlerFunc) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.prompts = append(sm.prompts, PromptEntry{prompt: prompt, phf: phf})
}

func (sm *ServiceManager) AddNotificationHandler(name string, handler server.NotificationHandlerFunc) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.notificationHandlers[name] = handler
}

func (sm *ServiceManager) AddResourceHandler(resource mcp.Resource, handler server.ResourceHandlerFunc) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.reources[resource] = handler
}

func (sm *ServiceManager) AddResourceTemplateHandler(resourceTemplate mcp.ResourceTemplate, handler server.ResourceTemplateHandlerFunc) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.resourceTemplates[resourceTemplate] = handler
}

func (sm *ServiceManager) Ctx() context.Context {
	return sm.ctx
}

func (sm *ServiceManager) Tools() []server.ServerTool {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	return sm.tools
}

func (sm *ServiceManager) Prompts() []PromptEntry {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	return sm.prompts
}

func (sm *ServiceManager) NotificationHandlers() map[string]server.NotificationHandlerFunc {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	return sm.notificationHandlers
}

func (sm *ServiceManager) Resources() map[mcp.Resource]server.ResourceHandlerFunc {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	return sm.reources
}

func (sm *ServiceManager) ResourceTemplates() map[mcp.ResourceTemplate]server.ResourceTemplateHandlerFunc {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	return sm.resourceTemplates
}
