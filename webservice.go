package shopline

import (
	"time"
)

const (
	URL_BOLETO   = "https://shopline.itau.com.br/shopline/Itaubloqueto.asp"
	URL_CONSULTA = "https://shopline.itau.com.br/shopline/consulta.asp"
	URL_SHOPLINE = "https://shopline.itau.com.br/shopline/shopline.asp"
)

type Webservice struct {
	Codigo    string
	Chave     string
	ChaveItau string
}

type BoletoDef struct {
	Valor           float64
	Obs             string
	NomeCliente     string
	CodigoInscricao string
	NumeroInscricao string
	CEP             string
	Endereco        string
	Bairro          string
	Cidade          string
	Estado          string
	Vencimento      time.Time
}

func New(codigo, chave string) *Webservice {
	ws := &Webservice{codigo, chave, "SEGUNDA12345ITAU"}
	return ws
}

func (ws *Webservice) NewBoleto(boleto BoletoDef) {

}

func rjust(in, fill string, length int) string {
	for len(in) < length {
		in = fill + in
	}
	return in
}

func ljust(in, fill string, length int) string {
	for len(in) < length {
		in = in + fill
	}
	return in
}
