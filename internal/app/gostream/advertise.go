package gostream

import (
	"context"
	"runtime"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// Advertise ...
func (i *Implementation) Advertise(ctx context.Context, req *gostreamv1.AdvertiseRequest) (*gostreamv1.AdvertiseResponse, error) {
	result := make([]*gostreamv1.AdvertiseResponse_ServerMetricResponse, 0)
	for _, metric := range req.GetMetrics() {
		switch metric {
		case gostreamv1.ServerMetric_SERVER_METRIC_GOROUTINES:
			result = append(result, &gostreamv1.AdvertiseResponse_ServerMetricResponse{
				Metric: metric,
				Value:  float64(runtime.NumGoroutine()),
			})
		case gostreamv1.ServerMetric_SERVER_METRIC_STREAMS:
			result = append(result, &gostreamv1.AdvertiseResponse_ServerMetricResponse{
				Metric: metric,
				Value:  float64(i.watchers.Count()),
			})
		}
	}

	return &gostreamv1.AdvertiseResponse{
		Response: result,
	}, nil
}
