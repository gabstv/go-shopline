# go-shopline
Itau shopline for Google Go.  

## Install
```bash
go get github.com/gabstv/go-shopline
```

## Examples

### Get Bank Redirect Page  

```Go
import(
	"github.com/gabstv/go-shopline"
)

func main(){
	ws := shopline.New("COD_EMP", "KEY")

	boleto := shopline.BoletoDef{}
	boleto.Observacao = "3"
	boleto.Obs1 = "NAO ACEITAR APOS O VENCIMENTO"
	boleto.Obs2 = "NUMERO DO ORCAMENTO: 12"
	boleto.NomeCliente = "Teste"
	boleto.CodigoInscricao = shopline.CPF
	boleto.NumeroInscricao = "45918713638"
	boleto.Endereco = "Praca da Se"
	boleto.Bairro = "Centro"
	boleto.CEP = "00001000"
	boleto.Cidade = "Sao Paulo"
	boleto.Estado = "SP"
	boleto.Vencimento = time.Now().AddDate(0, 0, 4)
	boleto.URL_Retorno = ""
	boleto.Pedido = 123456
	boleto.Valor = 3.5
	htmlstring, err := ws.GetBoletoRedirectHTML(boleto)
	if err != nil {
		print(err.Error()+"\n")
		return
	}
	print(htmlstring+"\n")
}
```

### Get Boleto PDF  

```Go
import(
	"github.com/gabstv/go-shopline"
	"fmt"
)

func main(){
	ws := shopline.New("COD_EMP", "KEY")

	boleto := shopline.BoletoDef{}
	boleto.Observacao = "3"
	boleto.Obs1 = "NAO ACEITAR APOS O VENCIMENTO"
	boleto.Obs2 = "NUMERO DO ORCAMENTO: 12"
	boleto.NomeCliente = "Teste"
	boleto.CodigoInscricao = shopline.CPF
	boleto.NumeroInscricao = "45918713638"
	boleto.Endereco = "Praca da Se"
	boleto.Bairro = "Centro"
	boleto.CEP = "00001000"
	boleto.Cidade = "Sao Paulo"
	boleto.Estado = "SP"
	boleto.Vencimento = time.Now().AddDate(0, 0, 4)
	boleto.URL_Retorno = ""
	boleto.Pedido = 123456
	boleto.Valor = 3.5
	boletopdf, err := ws.GetBoletoPDF(boleto, "https://www.example.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Boleto bytes:", boletopdf)
}
```