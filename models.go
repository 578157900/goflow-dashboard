package main

// Flow object to retrieve and response flows details
type Flow struct {
	Name string `json:"name"`
}

type DashboardSpec struct {
	TotalFlows     int
	ReadyFlows     int
	TotalRequests  int
	ActiveRequests int
}

type Location struct {
	Name string
	Link string
}

type FlowDesc struct {
	Name            string `json:"name"`
	Dot             string `json:"dot,omitempty"`
	InvocationCount int
}

type FlowRequests struct {
	Flow             string
	TracingEnabled   bool
	Requests         map[string]*RequestTrace
	CurrentRequestID string
}

// NodeTrace traces of each nodes in a dag
type NodeTrace struct {
	Node      string `json:"node"`
	StartTime int    `json:"start-time"`
	Duration  int    `json:"duration"`
	// Other can be added based on the needs
}

// RequestTrace object to retrieve and response traces details
type RequestTrace struct {
	RequestID       string       `json:"request-id"`
	FlowName        string       `json:"flow-name"`
	TraceId         string       `json:"trace-id"`
	NodeTraces      []*NodeTrace `json:"traces"`
	StartTime       int          `json:"start-time"`
	StartTimeFormat string       `json:"start-time-format"`
	Duration        int          `json:"duration"`
	Status          string       `json:"status"`
}
