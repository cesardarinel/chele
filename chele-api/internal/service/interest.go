package service

import (
	"math"
	"time"
)

func CalcularInteresDiario(saldo, tasaAnual float64, diasAtraso int) float64 {
	if saldo <= 0 || tasaAnual <= 0 || diasAtraso <= 0 {
		return 0
	}
	return math.Round(saldo*(tasaAnual/12/30)*float64(diasAtraso)*100) / 100
}

func GenerarIntereses(saldo, tasaAnual float64, desde, hasta time.Time) float64 {
	if hasta.IsZero() {
		hasta = time.Now()
	}
	dias := int(hasta.Sub(desde).Hours() / 24)
	if dias < 0 {
		dias = 0
	}
	return CalcularInteresDiario(saldo, tasaAnual, dias)
}
