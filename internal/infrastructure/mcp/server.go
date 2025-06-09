package mcp

import (
	"encoding/json"
	"fmt"
	"log"
)

// =============================================================================
// MCP Server - Model Context Protocol サーバー実装
// =============================================================================

// MCPServer はMCPプロトコルサーバーを表します
type MCPServer struct {
	name        string
	version     string
	description string
	tools       map[string]Tool
}

// Tool はMCPツールを表します
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
	Handler     ToolHandler
}

// ToolHandler はツールのハンドラー関数です
type ToolHandler func(arguments map[string]interface{}) (interface{}, error)

// NewMCPServer は新しいMCPサーバーを作成します
func NewMCPServer(name, version, description string) *MCPServer {
	return &MCPServer{
		name:        name,
		version:     version,
		description: description,
		tools:       make(map[string]Tool),
	}
}

// RegisterTool はツールを登録します
func (s *MCPServer) RegisterTool(tool Tool) {
	s.tools[tool.Name] = tool
	log.Printf("Registered MCP tool: %s", tool.Name)
}

// GetTools は登録されているツール一覧を返します
func (s *MCPServer) GetTools() []Tool {
	tools := make([]Tool, 0, len(s.tools))
	for _, tool := range s.tools {
		// ハンドラーを除外して返す（JSONシリアライズのため）
		tools = append(tools, Tool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}
	return tools
}

// ExecuteTool はツールを実行します
func (s *MCPServer) ExecuteTool(name string, arguments map[string]interface{}) (interface{}, error) {
	tool, exists := s.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	log.Printf("Executing tool: %s with arguments: %v", name, arguments)
	return tool.Handler(arguments)
}

// ServerInfo はサーバー情報を返します
func (s *MCPServer) ServerInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":    s.name,
		"version": s.version,
	}
}

// =============================================================================
// JSON-RPC 2.0 Message Types
// =============================================================================

// JSONRPCRequest はJSON-RPC 2.0リクエストです
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCResponse はJSON-RPC 2.0レスポンスです
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError はJSON-RPCエラーです
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HandleRequest はJSON-RPCリクエストを処理します
func (s *MCPServer) HandleRequest(requestData []byte) []byte {
	log.Printf("Received request: %s", string(requestData))
	
	var request JSONRPCRequest
	if err := json.Unmarshal(requestData, &request); err != nil {
		response := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      nil,
			Error: &JSONRPCError{
				Code:    -32700,
				Message: "Parse error",
				Data:    err.Error(),
			},
		}
		responseData, _ := json.Marshal(response)
		log.Printf("Sending error response: %s", string(responseData))
		return responseData
	}

	log.Printf("Processing method: %s", request.Method)

	// 通知メッセージの場合はレスポンスを返さない
	if request.ID == nil {
		// 通知メッセージの処理
		switch request.Method {
		case "notifications/initialized":
			log.Printf("Client initialized")
		default:
			log.Printf("Unhandled notification: %s", request.Method)
		}
		return nil // 通知にはレスポンスを返さない
	}

	// メソッドの処理
	var result interface{}
	var rpcError *JSONRPCError

	switch request.Method {
	case "initialize":
		result = s.handleInitialize(request.Params)
	case "tools/list":
		result = s.handleToolsList()
	case "tools/call":
		result, rpcError = s.handleToolsCall(request.Params)
	default:
		rpcError = &JSONRPCError{
			Code:    -32601,
			Message: "Method not found",
			Data:    request.Method,
		}
	}

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
		Error:   rpcError,
	}

	responseData, _ := json.Marshal(response)
	log.Printf("Sending response: %s", string(responseData))
	return responseData
}

// handleInitialize は初期化リクエストを処理します
func (s *MCPServer) handleInitialize(params interface{}) interface{} {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": s.ServerInfo(),
	}
}

// handleToolsList はツール一覧リクエストを処理します
func (s *MCPServer) handleToolsList() interface{} {
	tools := s.GetTools()
	log.Printf("Returning %d tools: %v", len(tools), tools)
	return map[string]interface{}{
		"tools": tools,
	}
}

// handleToolsCall はツール実行リクエストを処理します
func (s *MCPServer) handleToolsCall(params interface{}) (interface{}, *JSONRPCError) {
	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, &JSONRPCError{
			Code:    -32602,
			Message: "Invalid params",
		}
	}

	name, ok := paramsMap["name"].(string)
	if !ok {
		return nil, &JSONRPCError{
			Code:    -32602,
			Message: "Missing or invalid tool name",
		}
	}

	arguments, ok := paramsMap["arguments"].(map[string]interface{})
	if !ok {
		arguments = make(map[string]interface{})
	}

	result, err := s.ExecuteTool(name, arguments)
	if err != nil {
		return nil, &JSONRPCError{
			Code:    -32603,
			Message: "Tool execution error",
			Data:    err.Error(),
		}
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": result,
			},
		},
	}, nil
}
