package transport_factory

import (
	"course-go-avito-SitnikovArtem06/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDeadline_Success(t *testing.T) {

	t.Parallel()

	df := NewTransportFactory()

	tests := []struct {
		name      string
		transport model.TransportType
		wantMin   time.Duration
	}{
		{"on_foot", model.OnFoot, 30 * time.Minute},
		{"scooter", model.Scooter, 15 * time.Minute},
		{"car", model.Car, 5 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			now := time.Now().UTC()

			tr := df.Get(tt.transport)
			got := tr.Deadline()

			diff := got.Sub(now)

			if diff < tt.wantMin-1*time.Second || diff > tt.wantMin+1*time.Second {
				t.Fatalf("deadline diff = %v, want around %v", diff, tt.wantMin)
			}
		})
	}

}

func TestTransportFactory_DefaultOnUnknown(t *testing.T) {
	f := NewTransportFactory()

	tr := f.Get(model.TransportType(rune(100)))
	require.NotNil(t, tr)
	
	deadline := tr.Deadline()
	require.WithinDuration(t,
		time.Now().Add(30*time.Minute),
		deadline,
		2*time.Second,
	)
}
