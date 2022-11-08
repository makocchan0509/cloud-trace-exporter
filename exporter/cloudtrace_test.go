package exporter

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"os"
	"testing"
	"time"
)

func TestInit(t *testing.T) {

	type args struct {
		project        string
		serviceNameKey string
	}

	t.Run("Regular pattern", func(t *testing.T) {
		a := args{
			project:        os.Getenv("PROJECT_ID"),
			serviceNameKey: "test code service",
		}
		ctx := context.Background()
		if tp, err := Init(ctx, a.project, a.serviceNameKey); err != nil {
			t.Errorf("Init() = %v \n", err)
		} else {
			defer tp.Flush(ctx)

			ctx := context.Background()
			tracer := otel.GetTracerProvider().Tracer("exporter/testing")
			ctx, span := tracer.Start(ctx, "testing", trace.WithAttributes(attribute.String("Id", "testing")))

			time.Sleep(time.Second * 1)

			defer span.End()
		}
	})

	t.Run("Not found project pattern", func(t *testing.T) {
		a := args{
			project:        "",
			serviceNameKey: "test code service",
		}
		ctx := context.Background()
		if _, err := Init(ctx, a.project, a.serviceNameKey); err == nil {
			t.Errorf("Init() want error,but not found error")
		}
	})

	t.Run("Not found service pattern(regular)", func(t *testing.T) {
		a := args{
			project:        os.Getenv("PROJECT_ID"),
			serviceNameKey: "",
		}
		ctx := context.Background()
		if _, err := Init(ctx, a.project, a.serviceNameKey); err != nil {
			t.Errorf("Init() = %v \n", err)
		}
	})
}
