package utils

import "testing"

func TestCNPJValido(t *testing.T) {
	if err := ValidarCnpj("17.060.704/0001-52"); err != nil {
		t.Fatal(err)
	}
}

func TestCNPJValido2(t *testing.T) {
	if err := ValidarCnpj("24.871.503/0001-09"); err != nil {
		t.Fatal(err)
	}
}

func TestCNPJInvalido(t *testing.T) {
	if err := ValidarCnpj("24.871.503/0001-03"); err == nil {
		t.Fatal("Falha ao validar o cnpj")
	}
}

func TestCPF(t *testing.T) {
	if err := ValidarCpf("529.982.247-25"); err != nil {
		t.Fatal(err)
	}
}

func TestCPFInvalido(t *testing.T) {
	if err := ValidarCpf("529.982.247-35"); err == nil {
		t.Fatal(err)
	}
}
