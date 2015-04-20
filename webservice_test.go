package shopline

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

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

/*func Test2(t *testing.T) {
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
}*/

func TestBoletoReal(t *testing.T) {
	ws := New("J0012345678901234567890123", "A3G8E4C19N6W7BPS")
	ws.Homologacao = true
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
	dc, err := ws.process(bl)
	t.Log("DC", dc, "\n")
	if err != nil {
		t.Fatalf("Err ftt %v\n", err.Error())
	}
	resp, err := http.PostForm(URL_BOLETO, url.Values{"DC": {dc}})
	if err != nil {
		t.Fatalf("[ITAU ERR] %v\n", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("[HTTP STATUS ERR] %v - %v\n", resp.StatusCode, resp.Status)
	}
	b := new(bytes.Buffer)
	io.Copy(b, resp.Body)
	resp.Body.Close()
	t.Log("[OK]", resp.StatusCode, resp.Status, "{{{ ", b.String(), " }}}\n")
}

func TestBoletoReal2(t *testing.T) {
	ws := New("J0012345678901234567890123", "A3G8E4C19N6W7BPS")
	ws.Homologacao = true
	bl := BoletoDef{}
	bl.Pedido = 101
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
	dc, err := ws.process(bl)
	t.Log("DC", dc, "\n")
	if err != nil {
		t.Fatalf("Err ftt %v\n", err.Error())
	}
	t.Log(make_post_page(dc) + "\n")
}
