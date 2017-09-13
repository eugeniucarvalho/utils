package utils

import(
    "fmt"
    "regexp"
    "strconv"
)

func cnpjDigitoVerificado(cnpj []int, b int) (dv int) {
    i, j, soma, limit := 0, 4 + b, 0, 11 + b

    for ; i < limit; i++ {
        soma +=  cnpj[i] * j;
        if j == 2 {  j = 9 } else { j -= 1 }
    }
    resto := soma % 11
    if resto < 2 { dv = 0  } else { dv = 11 - resto }
    return
}

func somenteNumeros(c string, vsize int) ([]int, error) {
    tmp   := regexp.MustCompile("[^0-9]").ReplaceAll([]byte(c), []byte("") )
    size  := len(tmp)
    cnpji := make([]int, size)
    // Valida tamanho
    if len(cnpji) != vsize {
        return cnpji, fmt.Errorf("Tamanho deve conter %d digitos sem caracteres especiais", vsize);
    }

    for i, val := range tmp {
        cnpji[i], _ = strconv.Atoi(string(val))
    }

    return cnpji, nil
}

func ValidarCnpj(cnpj string) error {

    cnpji, err := somenteNumeros(cnpj, 14)
    if err != nil {
        return err
    }

    // Valida primeiro dígito verificador
    if cnpji[12] == cnpjDigitoVerificado(cnpji, 1) {
        // Valida segundo dígito verificador
        if cnpji[13] == cnpjDigitoVerificado(cnpji, 2) {
            return nil
        }
    }
    return fmt.Errorf("CNPJ inválido");
}


func cpfDigitoVerificado(cpf []int, b int) (dv int) {
    soma, ini := 0, 9 + b
    for i := ini; i > 1; i -- {
        soma += i * cpf[ini - i]
    }
    resto := soma * 10%11
    if resto == 10 { dv = 0 } else { dv = resto }
    return
}

func ValidarCpf(cpf string) error {
    cpfi, err := somenteNumeros(cpf, 11)
    if err != nil {
        return err
    }

    // Valida primeiro dígito verificador
    if cpfi[9] == cpfDigitoVerificado(cpfi, 1) {
        // Valida segundo dígito verificador
        if cpfi[10] == cpfDigitoVerificado(cpfi, 2) {
            return nil
        }
    }
    return fmt.Errorf("CPF inválido");
}
