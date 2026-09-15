package backend

func evaluateMetric(v MetricSample) *Alert { if v.Name=="packet_loss" && v.Value>=5 { return &Alert{Severity:"warning",Source:"telemetry",Message:"packet loss threshold exceeded",Active:true,ObservedAt:v.ObservedAt} }; if v.Name=="latency_ms" && v.Value>=100 { return &Alert{Severity:"warning",Source:"telemetry",Message:"latency threshold exceeded",Active:true,ObservedAt:v.ObservedAt} }; return nil }
