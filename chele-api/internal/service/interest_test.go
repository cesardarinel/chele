package service

import (
	"testing"
	"time"
)

func TestCalcularInteresDiario(t *testing.T) {
	tests := []struct {
		name      string
		saldo     float64
		tasa      float64
		dias      int
		esperado  float64
	}{
		{"sin saldo", 0, 0.96, 30, 0},
		{"sin tasa", 1000, 0, 30, 0},
		{"sin dias", 1000, 0.96, 0, 0},
		{"1000 a 96% anual 30 dias", 1000, 0.96, 30, 80},
		{"5000 a 36% anual 15 dias", 5000, 0.36, 15, 75},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcularInteresDiario(tt.saldo, tt.tasa, tt.dias)
			if got != tt.esperado {
				t.Errorf("got %.2f, want %.2f", got, tt.esperado)
			}
		})
	}
}

func TestCalcularInteresDiario_Negativo(t *testing.T) {
	got := CalcularInteresDiario(-500, 0.96, 30)
	if got != 0 {
		t.Errorf("expected 0 for negative balance, got %.2f", got)
	}
}

func TestGenerarIntereses(t *testing.T) {
	desde := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	hasta := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	got := GenerarIntereses(1000, 0.96, desde, hasta)
	if got <= 0 {
		t.Errorf("expected >0 interest, got %.2f", got)
	}
}
