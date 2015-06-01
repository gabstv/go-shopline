package shopline

import (
	"bytes"
	"crypto/rand"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	R  = false
	r0 = regexp.MustCompile("<consulta>(.+)\\s*<\\/consulta>")
)

type SondaResult struct {
	XMLName       xml.Name   `xml:"PARAMETER"`
	Z             []XMLParam `xml:"PARAM"`
	CodigoEmpresa string
	Pedido        int
	Valor         float64
	TipoPgto      string
	SitPgto       string
	DtPgto        time.Time
}

func (s *SondaResult) GetRawVal(name string) string {
	for _, v := range s.Z {
		if v.Id == name {
			return v.Value
		}
	}
	return ""
}

func (s *SondaResult) Unwrap() {
	s.CodigoEmpresa = s.GetRawVal("CodEmp")
	s.Pedido, _ = strconv.Atoi(s.GetRawVal("Pedido"))
	s.Valor, _ = strconv.ParseFloat(strings.Replace(s.GetRawVal("Valor"), ",", ".", -1), 64)
	s.TipoPgto = s.GetRawVal("tipPag")
	s.SitPgto = s.GetRawVal("sitPag")
	dmy := s.GetRawVal("dtPag")
	if len(dmy) == 8 {
		d, _ := strconv.Atoi(dmy[:2])
		m, _ := strconv.Atoi(dmy[2:4])
		y, _ := strconv.Atoi(dmy[4:])
		s.DtPgto = time.Date(y, time.Month(m), d, 12, 0, 0, 0, time.UTC)
	}
}

type XMLParam struct {
	Id    string `xml:"ID,attr"`
	Value string `xml:"VALUE,attr"`
}

const (
	URL_BOLETO         = "https://shopline.itau.com.br/shopline/Impressao.aspx"
	URL_BOLETO_HOMOLOG = "https://shopline.itau.com.br/shopline/emissao_teste.asp"
	URL_CONSULTA       = "https://shopline.itau.com.br/shopline/consulta.aspx"
	URL_SHOPLINE       = "https://shopline.itau.com.br/shopline/shopline.aspx"
	CPF                = "01"
	CNPJ               = "02"
	// status pagamento
	STAT_PAGAMENTO_EFETUADO           = "00"
	STAT_SIT_PAGAMENTO_NAO_FINALIZADO = "01"
	STAT_ERR_NA_CONSULTA              = "02"
	STAT_PEDIDO_NAO_LOCALIZADO        = "03"
	STAT_BOLETO_EMITIDO_COM_SUCESSO   = "04"
	STAT_PGTO_EFETUADO_AG_COMPENSACAO = "05"
	STAT_PGTO_NAO_COMPENSADO          = "06"
)

var trmap = map[rune]rune{
	'Š': 'S', 'Œ': 'O', 'Ž': 'Z', 'š': 's', 'œ': 'o', 'ž': 'z', 'Ÿ': 'Y', '¥': 'Y',
	'µ': 'u', 'À': 'A', 'Á': 'A', 'Â': 'A', 'Ã': 'A', 'Ä': 'A', 'Å': 'A', 'Æ': 'A',
	'Ç': 'C', 'È': 'E', 'É': 'E', 'Ê': 'E', 'Ë': 'E', 'Ì': 'I', 'Í': 'I', 'Î': 'I',
	'Ï': 'I', 'Ð': 'D', 'Ñ': 'N', 'Ò': 'O', 'Ó': 'O', 'Ô': 'O', 'Õ': 'O', 'Ö': 'O',
	'Ø': 'O', 'Ù': 'U', 'Ú': 'U', 'Û': 'U', 'Ü': 'U', 'Ý': 'Y', 'ß': 's', 'à': 'a',
	'á': 'a', 'â': 'a', 'ã': 'a', 'ä': 'a', 'å': 'a', 'æ': 'a', 'ç': 'c', 'è': 'e',
	'é': 'e', 'ê': 'e', 'ë': 'e', 'ì': 'i', 'í': 'i', 'î': 'i', 'ï': 'i', 'ð': 'o',
	'ñ': 'n', 'ò': 'o', 'ó': 'o', 'ô': 'o', 'õ': 'o', 'ö': 'o', 'ø': 'o', 'ù': 'u',
	'ú': 'u', 'û': 'u', 'ü': 'u', 'ý': 'y', 'ÿ': 'y',
}

var whitelist = []int{
	32, 33, 34, 37, 39, 40, 41, 42, 43, 44, 45, 45, 46, 48, 49, 50, 51, 52, 53, 54,
	55, 56, 57, 58, 60, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
	77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 95, 97, 98, 99, 100, 101,
	102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118,
	119, 120, 121, 122,
}

func fixrune(r rune) rune {
	if v, ok := trmap[r]; ok {
		r = v
	}
	v2 := int(r)
	i := sort.Search(len(whitelist), func(i int) bool {
		return whitelist[i] >= v2
	})
	if i < len(whitelist) && whitelist[i] == v2 {
		return r
	}
	// not whitelisted
	return ' '
}

func cleanstr(v string) string {
	b := new(bytes.Buffer)
	for _, r := range v {
		b.WriteRune(fixrune(r))
	}
	return strings.ToUpper(b.String())
}

func strim(v string, maxlen int) string {
	if len(v) > maxlen {
		return v[:maxlen]
	}
	return v
}

type Webservice struct {
	Codigo      string
	Chave       string
	ChaveItau   string
	Homologacao bool
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

func (b BoletoDef) Clean() BoletoDef {
	c := BoletoDef{}
	c.Pedido = b.Pedido
	c.Valor = b.Valor
	c.Vencimento = b.Vencimento
	c.CodigoInscricao = b.cleanCodigoInscricao()
	c.Observacao = ljust(strim(cleanstr(b.Observacao), 40), " ", 40)
	c.Obs1 = ljust(strim(cleanstr(b.Obs1), 60), " ", 60)
	c.Obs2 = ljust(strim(cleanstr(b.Obs2), 60), " ", 60)
	c.Obs3 = ljust(strim(cleanstr(b.Obs3), 60), " ", 60)
	c.NomeCliente = ljust(strim(cleanstr(b.NomeCliente), 30), " ", 30)
	c.NumeroInscricao = ljust(strim(cleanstr(b.NumeroInscricao), 14), " ", 14)
	c.Endereco = ljust(strim(cleanstr(b.Endereco), 40), " ", 40)
	c.Bairro = ljust(strim(cleanstr(b.Bairro), 15), " ", 15)
	c.CEP = ljust(strim(cleanstr(b.CEP), 8), " ", 8)
	c.Cidade = ljust(strim(cleanstr(b.Cidade), 15), " ", 15)
	c.Estado = ljust(strim(cleanstr(b.Estado), 2), " ", 2)
	c.URL_Retorno = ljust(strim(cleanstr(b.URL_Retorno), 60), " ", 60)
	return c
}

func (b BoletoDef) cleanPedido() string {
	return rjust(strconv.Itoa(b.Pedido), "0", 8)
}

func (b BoletoDef) cleanValor() string {
	vlrs := strings.Replace(strings.Replace(moneyf(b.Valor), ",", "", -1), ".", "", -1)
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
	c := b.Clean()
	buffer.WriteString(c.cleanPedido())
	buffer.WriteString(c.cleanValor())
	buffer.WriteString(c.Observacao)
	buffer.WriteString(c.NomeCliente)
	buffer.WriteString(c.CodigoInscricao)
	buffer.WriteString(c.NumeroInscricao)
	buffer.WriteString(c.Endereco)
	buffer.WriteString(c.Bairro)
	buffer.WriteString(c.CEP)
	buffer.WriteString(c.Cidade)
	buffer.WriteString(c.Estado)
	buffer.WriteString(c.cleanVencimento())
	buffer.WriteString(c.URL_Retorno)
	buffer.WriteString(c.Obs1)
	buffer.WriteString(c.Obs2)
	buffer.WriteString(c.Obs3)
	return buffer.String()
}

func New(codigo, chave string) *Webservice {
	ws := &Webservice{codigo, chave, "SEGUNDA12345ITAU", false}
	return ws
}

func (ws *Webservice) GetDC(boleto BoletoDef) (string, error) {
	dc, err := ws.process(boleto)
	return dc, err
}

// encapsula o boleto e retorna o PDF
// BETA
func (ws *Webservice) GetBoletoPDF(boleto BoletoDef) ([]byte, error) {
	dc, err := ws.process(boleto)
	if err != nil {
		return nil, err
	}

	//
	// 1 - Visit Shopline landing page
	cjar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.New("HTTP COOKIEJAR ERR: " + err.Error())
	}
	httpc := &http.Client{
		Jar:     cjar,
		Timeout: time.Second * 30,
	}
	//
	postv := url.Values{}
	postv.Set("DC", dc)
	buffer := new(bytes.Buffer)
	buffer.WriteString(postv.Encode())
	//
	req, err := http.NewRequest("POST", URL_SHOPLINE, buffer)
	if err != nil {
		return nil, err
	}
	// (͡° ͜ʖ ͡°)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	// Do request
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	//
	if resp.StatusCode != 200 {
		return nil, errors.New("HTTP STATUS CODE " + strconv.Itoa(resp.StatusCode))
	}
	buffer.Reset()
	io.Copy(buffer, resp.Body) // read all, so the server is happy
	resp.Body.Close()
	//
	// some cookies should have been grabbed by now
	//
	postv = url.Values{}
	postv.Set("cliente", "N")
	postv.Set("CodEmp", ws.Codigo)
	postv.Set("IdSite", "29772")
	postv.Set("npedido", boleto.cleanPedido())
	postv.Set("flag", "1")
	postv.Set("DC", dc)
	postv.Set("emissao", "1")
	postv.Set("soagenda", "")
	buffer.Reset()
	buffer.WriteString(postv.Encode())
	//
	req, err = http.NewRequest("POST", URL_BOLETO, buffer)
	if err != nil {
		return nil, err
	}
	// (͡° ͜ʖ ͡°)
	req.Header.Set("Referer", URL_SHOPLINE)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	// Finally do the request we want to
	resp, err = httpc.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("HTTP STATUS CODE " + strconv.Itoa(resp.StatusCode))
	}
	buffer.Reset()
	io.Copy(buffer, resp.Body)
	resp.Body.Close()
	return buffer.Bytes(), nil
}

func (ws *Webservice) GetBoletoRedirectHTML(boleto BoletoDef) (string, error) {
	dc, err := ws.process(boleto)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.WriteString("<html><body><form method=\"post\" action=\"")
	buf.WriteString(URL_SHOPLINE)
	buf.WriteString("\" id=\"itaushopline\"><input type=\"hidden\" name=\"DC\" value=\"")
	buf.WriteString(dc)
	buf.WriteString("\"></form><script>document.getElementById('itaushopline').submit();</script></body></html>")
	return buf.String(), nil
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
	chave1 := algoritmo([]byte(boleto.ToToken()), ws.Chave)
	chave2 := algoritmo(append([]byte(ws.Codigo), chave1...), ws.ChaveItau)
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

func inicializa(chave string) (indices []int, asc_codes []rune) {
	chave = strings.ToUpper(chave)
	// inicializa
	indices = make([]int, 256)
	asc_codes = make([]rune, 256)
	for i := 0; i < 256; i++ {
		asc_codes[i] = rune(chave[i%len(chave)])
		indices[i] = i
	}
	l := 0
	for k := 0; k < 256; k++ {
		l = (l + indices[k] + int(asc_codes[k])) % 256

		i := indices[k]
		indices[k] = indices[l]
		indices[l] = i
	}
	return
}

// Itau usa encryption própria
// token =
// 'pedido', 'valor', 'observacao',
//          'nome', 'codigo_inscricao', 'numero_inscricao', 'endereco', 'bairro', 'cep',
//          'cidade', 'estado', 'vencimento', 'url_retorno', 'obs_1', 'obs_2', 'obs_3'
func algoritmo(token []byte, chave string) []byte {
	chave = strings.ToUpper(chave)
	// inicializa
	indices, _ := inicializa(chave)
	// algoritmo
	var data_chave bytes.Buffer
	m := 0
	k := 0
	for j := 1; j <= len(token); j++ {
		k = (k + 1) % 256
		m = (m + indices[k]) % 256
		i := indices[k]
		indices[k] = indices[m]
		indices[m] = i
		n := indices[((indices[k] + indices[m]) % 256)]
		caracter := byte(int(token[j-1]) ^ n)
		data_chave.WriteByte(caracter)
	}
	return data_chave.Bytes()
}
func rr(min, max int) int {
	if R {
		return min
	}
	b := make([]byte, 1)
	rand.Read(b)
	return min + int((float64(b[0])/256)*float64(max-min))
}
func rnd() rune {
	alfa := "ABCDEFGHIJKLMNOPQRSTUVXWYZ"
	return rune(alfa[rr(0, len(alfa))])
}
func converte(chave []byte) string {
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
	chave1 := algoritmo([]byte(rjust(strconv.Itoa(pedido), "0", 8)+formato), ws.Chave)
	chave2 := algoritmo(append([]byte(ws.Codigo), chave1...), ws.ChaveItau)
	dc := converte(chave2)
	resp, err := http.Get(URL_CONSULTA + "?DC=" + dc)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("HTTP STATUS CODE " + resp.Status)
	}
	return resp.Body, nil
}

func (ws *Webservice) Sonda(pedido int) (*SondaResult, error) {
	rc, err := ws.sonda(pedido, "1")
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, rc)
	rc.Close()
	if !r0.MatchString(buf.String()) {
		return nil, errors.New("XML MISMATCH: " + buf.String())
	}
	subm := r0.FindStringSubmatch(buf.String())
	v := &SondaResult{}
	err = xml.Unmarshal([]byte(subm[1]), v)
	if err != nil {
		return nil, err
	}
	v.Unwrap()
	return v, nil
}

// gabs copypasta

func moneyf(inp float64) string {
	str0 := fmt.Sprintf("%0.2f", inp)
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
