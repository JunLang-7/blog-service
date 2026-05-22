package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JunLang-7/blog-service/global"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func setupTestTracer(t *testing.T) func() {
	t.Helper()
	tracer, closer := jaeger.NewTracer(
		"blog-service-test",
		jaeger.NewConstSampler(true),
		jaeger.NewNullReporter(),
	)
	global.Tracer = tracer
	opentracing.SetGlobalTracer(tracer)
	return func() {
		closer.Close()
		global.Tracer = nil
	}
}

func TestTracing_CreatesRootSpan(t *testing.T) {
	cleanup := setupTestTracer(t)
	defer cleanup()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	var traceID, spanID any
	r.Use(Tracing())
	r.GET("/api/v1/tags", func(c *gin.Context) {
		traceID, _ = c.Get("X-Trace-ID")
		spanID, _ = c.Get("X-Span-ID")
		c.String(http.StatusOK, "ok")
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/tags", nil)
	r.ServeHTTP(w, c.Request)

	if traceID == nil || traceID == "" {
		t.Error("expected X-Trace-ID to be set, got empty")
	}
	if spanID == nil || spanID == "" {
		t.Error("expected X-Span-ID to be set, got empty")
	}
	if traceID == "0" {
		t.Error("trace ID should not be zero")
	}
	if spanID == "0" {
		t.Error("span ID should not be zero")
	}
}

func TestTracing_CreatesChildSpanFromUpstreamHeaders(t *testing.T) {
	cleanup := setupTestTracer(t)
	defer cleanup()

	parentSpan := global.Tracer.StartSpan("parent-operation")
	defer parentSpan.Finish()

	headers := http.Header{}
	carrier := opentracing.HTTPHeadersCarrier(headers)
	if err := global.Tracer.Inject(parentSpan.Context(), opentracing.HTTPHeaders, carrier); err != nil {
		t.Fatalf("failed to inject span: %v", err)
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	var traceID, spanID any
	r.Use(Tracing())
	r.GET("/api/v1/articles/1", func(c *gin.Context) {
		traceID, _ = c.Get("X-Trace-ID")
		spanID, _ = c.Get("X-Span-ID")
		c.String(http.StatusOK, "ok")
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/articles/1", nil)
	c.Request.Header = headers
	r.ServeHTTP(w, c.Request)

	if traceID == nil || traceID == "" {
		t.Error("expected X-Trace-ID to be set, got empty")
	}
	if spanID == nil || spanID == "" {
		t.Error("expected X-Span-ID to be set, got empty")
	}
}

func TestTracing_DifferentRequestsHaveDifferentSpanIDs(t *testing.T) {
	cleanup := setupTestTracer(t)
	defer cleanup()

	gin.SetMode(gin.TestMode)

	var spanIDs []string
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.Use(Tracing())
		r.GET("/test", func(c *gin.Context) {
			sid, _ := c.Get("X-Span-ID")
			spanIDs = append(spanIDs, sid.(string))
		})
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, c.Request)
	}

	if spanIDs[0] == spanIDs[1] || spanIDs[1] == spanIDs[2] || spanIDs[0] == spanIDs[2] {
		t.Error("each request should have a unique span ID")
	}
}
