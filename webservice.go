package shopline

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	URL_BOLETO   = "https://shopline.itau.com.br/shopline/Itaubloqueto.asp"
	URL_CONSULTA = "https://shopline.itau.com.br/shopline/consulta.asp"
	URL_SHOPLINE = "https://shopline.itau.com.br/shopline/shopline.asp"
	CPF          = "01"
	CNPJ         = "02"
)

type Webservice struct {
	Codigo    string
	Chave     string
	ChaveItau string
}

type BoletoDef struct {
	Pedido          int
	Valor           float64
	Observacao      string
	NomeCliente     string
	CodigoInscricao string
	NumeroInscricao string
	CEP             string
	Endereco        string
	Bairro          string
	Cidade          string
	Estado          string
	Vencimento      time.Time
	URL_Retorno     string
	Obs1            string
	Obs2            string
	Obs3            string
}

func (b BoletoDef) Clean() {

}

func (b BoletoDef) cleanPedido() string {
	return rjust(strconv.Itoa(b.Pedido), "0", 8)
}

func (b BoletoDef) cleanValor() string {
	vlrs := strings.Replace(strings.Replace(moneyf(b.Valor), ",", "", 0), ".", "", 0)
	return rjust(vlrs, "0", 10)
}

func (b BoletoDef) cleanVencimento() string {
	return fmt.Sprintf("%02d%02d%04d", b.Vencimento.Day(), int(b.Vencimento.Month()), b.Vencimento.Year())
}

func (b BoletoDef) cleanCodigoInscricao() string {
	return rjust(b.CodigoInscricao, "0", 2)
}

func (b BoletoDef) ToToken() string {
	var buffer bytes.Buffer
	buffer.WriteString(b.cleanPedido())
	buffer.WriteString(b.cleanValor())
	buffer.WriteString(b.Observacao)
	buffer.WriteString(b.NomeCliente)
	buffer.WriteString(b.cleanCodigoInscricao())
	buffer.WriteString(b.NumeroInscricao)
	buffer.WriteString(b.Endereco)
	buffer.WriteString(b.Bairro)
	buffer.WriteString(b.CEP)
	buffer.WriteString(b.Cidade)
	buffer.WriteString(b.cleanVencimento())
	buffer.WriteString(b.URL_Retorno)
	buffer.WriteString(b.Obs1)
	buffer.WriteString(b.Obs2)
	buffer.WriteString(b.Obs3)
	return buffer.String()
}

func New(codigo, chave string) *Webservice {
	ws := &Webservice{codigo, chave, "SEGUNDA12345ITAU"}
	return ws
}

func (ws *Webservice) NewBoleto(boleto BoletoDef) {

}
func (ws *Webservice) assert() error {
	if len(ws.Codigo) != 26 {
		return errors.New("Tamanho do codigo da empresa diferente de 26 posições")
	}
	if len(ws.Chave) != 16 {
		return errors.New("Tamanho da chave da chave diferente de 16 posições")
	}
	return nil
}
func (ws *Webservice) process(boleto BoletoDef) (res string, err error) {
	err = ws.assert()
	if err != nil {
		return
	}
	if v := rjust(boleto.CodigoInscricao, "0", 2); v != CPF && v != CNPJ {
		err = errors.New("Código de Inscrição Inválido 01 = CPF, 02 = CNPJ")
		return
	}
	boleto.Clean()
	chave1 := algoritmo(boleto.ToToken(), ws.Chave)
	chave2 := algoritmo(ws.Codigo+chave1, ws.ChaveItau)
	res = converte(chave2)
	return
}

func rjust(in, fillchar string, length int) string {
	for len(in) < length {
		in = fillchar + in
	}
	return in
}

func ljust(in, fillchar string, length int) string {
	for len(in) < length {
		in = in + fillchar
	}
	return in
}

// Itau usa encryption própria
// token =
// 'pedido', 'valor', 'observacao',
//          'nome', 'codigo_inscricao', 'numero_inscricao', 'endereco', 'bairro', 'cep',
//          'cidade', 'estado', 'vencimento', 'url_retorno', 'obs_1', 'obs_2', 'obs_3'
func algoritmo(token, chave string) string {
	// inicializa
	indices := make([]int, 256)
	asc_codes := make([]rune, 256)
	for i := 0; i < 256; i++ {
		asc_codes[i] = rune(chave[i%len(chave)])
		indices[i] = i
	}
	l := 0
	for k := 0; k < 256; k++ {
		l = (l + indices[k] + int(asc_codes[k])) % 256

		i := indices[k]
		indices[k] = indices[i]
		indices[l] = i
	}
	// algoritmo
	var data_chave bytes.Buffer
	l = 0
	for j := 1; j < len(token)+1; j++ {
		k := j % 256
		l = (l + indices[k]) % 256
		i := indices[k]
		indices[k] = indices[l]
		indices[l] = i
		//caracter = int(ord(token[(j-1):j]) ^ int(self.indices[(self.indices[k] + self.indices[l]) % 256]))
		caracter := rune(int(token[j-1]) ^ indices[(indices[k]+indices[l])%256])
		log.Println(caracter)
		data_chave.WriteRune(caracter)
	}
	return data_chave.String()
}
func rr(min, max int) int {
	b := make([]byte, 1)
	rand.Read(b)
	return min + int((float64(b[0])/256)*float64(max-min))
}
func rnd() rune {
	alfa := "ABCDEFGHIJKLMNOPQRSTUVXWYZ"
	return rune(alfa[rr(0, len(alfa))])
}
func converte(chave string) string {
	var data_rand bytes.Buffer
	data_rand.WriteRune(rnd())
	for i := 0; i < len(chave); i++ {
		data_rand.WriteString(strconv.Itoa(int(chave[i])))
		data_rand.WriteRune(rnd())
	}
	return data_rand.String()
}
func make_post_page(dc string) string {
	return fmt.Sprintf("<html><body><form method=\"post\" action=\"%v\" id=\"itaushopline\"><input type=\"hidden\" name=\"DC\" value=\"%v\"></form><script>document.getElementById('itaushopline').submit();</script></body></html>", URL_BOLETO, dc)
}
func (ws *Webservice) sonda(pedido int, formato string) (io.ReadCloser, error) {
	if formato != "0" && formato != "1" {
		if err := ws.assert(); err != nil {
			return nil, err
		}
	}
	chave1 := algoritmo(rjust(strconv.Itoa(pedido), "0", 8)+formato, ws.Chave)
	chave2 := algoritmo(ws.Codigo+chave1, ws.ChaveItau)
	dc := converte(chave2)
	cl := http.DefaultClient
	cl.Timeout = time.Second * 30
	buffer := new(bytes.Buffer)
	v := url.Values{}
	v.Set("DC", dc)
	buffer.WriteString(v.Encode())
	req, err := http.NewRequest("POST", URL_CONSULTA, buffer)
	if err != nil {
		return nil, err
	}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("HTTP STATUS CODE " + strconv.Itoa(resp.StatusCode) + " " + resp.Status)
	}
	return resp.Body, nil
}

// gabs copypasta

func moneyf(inp float64) string {
	str0 := fmt.Sprintf("%0.2f\n", inp)
	na := strings.Split(str0, ".")
	big := na[0]
	var frac string
	if len(na) < 2 {
		frac = "0"
	} else {
		frac = na[1]
	}
	b := new(bytes.Buffer)
	n := 0
	for i := len(big) - 1; i >= 0; i-- {
		b.WriteByte(big[i])
		n++
		if n == 3 {
			if i > 0 {
				if big[i-1] != '-' {
					b.WriteByte('.')
					n = 0
				}
			}
		}
	}
	// invert again
	str0 = b.String()
	b.Truncate(0)
	for i := len(str0) - 1; i >= 0; i-- {
		b.WriteByte(str0[i])
	}
	big = b.String()
	return big + "," + frac
}
