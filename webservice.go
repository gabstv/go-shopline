package shopline

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

func New(codigo, chave string) *Webservice {
	ws := &Webservice{codigo, chave, "SEGUNDA12345ITAU"}
	return ws
}
