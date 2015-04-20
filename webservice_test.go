package shopline

import (
	"testing"
	"time"
)

func Test0(t *testing.T) {
	t.Log("0\n")
	r := algoritmo("Now in the street, there is violence And-and a lots of work to be done No place to hang out our washing And-and I can't blame all on the sun", "1234567890123456")
	t.Log(r + "\n")
}

func Test1(t *testing.T) {
	r0 := rr(5, 7)
	r1 := rr(5, 7)
	r2 := rr(5, 7)
	if r0 < 5 || r0 > 6 {
		t.Fatalf("r0 is %v\n", r0)
	}
	if r1 < 5 || r1 > 6 {
		t.Fatalf("r1 is %v\n", r1)
	}
	if r2 < 5 || r2 > 6 {
		t.Fatalf("r2 is %v\n", r2)
	}
}

func Test2(t *testing.T) {
	ws := New("10101010101010101010101010", "2020202020222222")
	bl := BoletoDef{}
	bl.Pedido = 100
	bl.Valor = 100.25
	bl.Observacao = "Teste."
	bl.NomeCliente = "GABRIEL OCHSENHOFER"
	bl.CodigoInscricao = CNPJ
	bl.NumeroInscricao = "00000000000"
	bl.CEP = "04041002"
	bl.Endereco = "AV ONZE DE JUNHO 600 APT 82"
	bl.Bairro = "VILA CLEMENTINO"
	bl.Cidade = "SAO PAULO"
	bl.Estado = "SP"
	bl.Vencimento = time.Now().AddDate(0, 0, 2)
	proc, err := ws.process(bl)
	t.Log("PROC", proc, "\n")
	if err != nil {
		t.Fatalf("Err ftt %v\n", err.Error())
	}
}
