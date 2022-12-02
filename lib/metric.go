package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Objects to retrieve specific trace details

type SpanItem struct {
	TraceID       string `json:"traceID"`
	SpanID        string `json:"spanID"`
	OperationName string `json:"operationName"`
	StartTime     int    `json:"startTime"`
	Duration      int    `json:"duration"`
	// Other can be added based on the needs
}

type TraceItem struct {
	TraceID string      `json:"traceID"`
	Spans   []*SpanItem `json:"spans"`
}

type Traces struct {
	Data []*TraceItem `json:"data"`
}

// Objects to retrieve requests lists

type SpanOps struct {
	TraceID       string `json:"traceID"`
	SpanID        string `json:"spanID"`
	OperationName string `json:"operationName"`
	Tags          []Tag  `json:"tags"`
}

func (span *SpanOps) FindRequestID() string {
	for _, tag := range span.Tags {
		if tag.Key == "request" && tag.Type == "string" {
			return tag.Value.(string)
		}
	}
	return ""
}

type Tag struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type RequestItem struct {
	TraceID string     `json:"traceID"`
	Spans   []*SpanOps `json:"spans"`
}

type Requests struct {
	Data []*RequestItem `json:"data"`
}

// traces of each nodes in a dag
type NodeTrace struct {
	StartTime int    `json:"start-time"`
	Duration  int    `json:"duration"`
	Node      string `json:"node"`
	// Other can be added based on the needs
}

// RequestTrace object to response traces details
type RequestTrace struct {
	RequestID  string       `json:"request-id"`
	NodeTraces []*NodeTrace `json:"traces"`
	StartTime  int          `json:"start-time"`
	Duration   int          `json:"duration"`
}

var (
	trace_url = "http://localhost:16686/"
)

func ListRequests(function string) (map[string]string, error) {
	resp, err := http.Get(getTraceUrl() + "api/traces?service=goflow&operation=" + function)
	if err != nil {
		return nil, fmt.Errorf("failed to request trace service, error %v ", err)
	}
	defer resp.Body.Close()
	if resp.Body == nil {
		return nil, fmt.Errorf("failed to request trace service, status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read trace result, read error %v", err)
	}

	if len(bodyBytes) == 0 {
		return nil, fmt.Errorf("failed to get request traces, empty result")
	}
	requests := &Requests{}
	err = json.Unmarshal(bodyBytes, requests)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal requests lists, error %v", err)
	}

	requestMap := make(map[string]string)
	for _, request := range requests.Data {
		if request.Spans == nil {
			continue
		}
		for _, span := range request.Spans {
			if span.TraceID == request.TraceID && span.TraceID == span.SpanID {
				if requestID := span.FindRequestID(); requestID != "" {
					requestMap[requestID] = request.TraceID
				}
				break
			}
		}
	}
	return requestMap, nil
}

func GetTraceByTag(flowName string, tag map[string]string) (*RequestTrace, error) {
	tagData, _ := json.Marshal(tag)
	resp, err := http.Get(getTraceUrl() + fmt.Sprintf("api/traces?service=goflow&tags=%s", string(tagData)))
	if err != nil {
		return nil, fmt.Errorf("failed to request trace service, error %v ", err)
	}
	defer resp.Body.Close()
	if resp.Body == nil {
		return nil, fmt.Errorf("failed to request trace service, status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read trace result, read error %v", err)
	}

	if len(bodyBytes) == 0 {
		return nil, fmt.Errorf("failed to get request traces, empty result")
	}

	traces := &Traces{}
	err = json.Unmarshal(bodyBytes, traces)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal requests lists, error %v", err)
	}

	if traces.Data == nil || len(traces.Data) == 0 {
		return nil, fmt.Errorf("failed to get request traces, empty data")
	}

	requestTrace := traces.Data[0]
	requestTraces := &RequestTrace{}
	requestTraces.NodeTraces = []*NodeTrace{}

	var lastSpanEnd int

	for _, span := range requestTrace.Spans {
		if span.TraceID == span.SpanID {
			// Set RequestID, StartTime and lastestSpan start time
			requestTraces.RequestID = span.OperationName
			requestTraces.StartTime = span.StartTime
			requestTraces.Duration = span.Duration
			lastSpanEnd = span.StartTime
		} else {
			spanEndTime := span.StartTime + span.Duration
			if spanEndTime > lastSpanEnd {
				lastSpanEnd = spanEndTime
			}
			requestTraces.NodeTraces = append(requestTraces.NodeTraces, &NodeTrace{
				span.StartTime,
				span.Duration,
				span.OperationName + "-" + span.SpanID,
			})
		}
	}
	if lastSpanEnd > requestTraces.StartTime {
		requestTraces.Duration = lastSpanEnd - requestTraces.StartTime
	}

	return requestTraces, nil
}

func getTraceUrl() string {
	url := os.Getenv("TRACE_URL")
	if url != "" {
		trace_url = url
	}
	return trace_url
}
