package shopline

import (
	"testing"
)

func TestStrip(t *testing.T) {
	a := `<html><body><form method="post" action="https://shopline.itau.com.br/shopline/shopline.aspx" id="itaushopline"><input type="hidden" name="DC" value="TENDIES"></form><script>document.getElementById('itaushopline').submit();</script></body></html>`
	result, err := StripDCFromRedirectHTML(a)
	if err != nil {
		t.Fatal(err)
	}
	if result != "TENDIES" {
		t.Fatalf("Result should be TENDIES but it is %v\n", result)
	}
}

func TestStrim(t *testing.T) {
	a := "MARCIO DE MORAIS CAVEíRA DA ROZA"
	b := ljust(strim(cleanstr(a), 30), " ", 30)
	if b != "MARCIO DE MORAIS CAVEIRA DA RO" {
		t.Fatalf("b is %v\n", b)
	}
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

func TestConverte(t *testing.T) {
	a := "A49A50A51A52A53A54A"
	R = true
	b := converte([]byte("123456"))
	if a != b {
		t.Fatalf("a != b (%v != %v)\n", a, b)
	}
}

func TestAlgo(t *testing.T) {
	R = true
	vvb := algoritmo([]byte("123456789"), "BOLOVO")
	vv := converte(algoritmo(append([]byte("LOL"), vvb...), "ZUZUZU"))
	if vv != "A208A60A194A122A171A13A242A39A77A155A165A248A" {
		t.Fatalf("Algoritmo error: '%v' != '%v'\n%v\n", vv, "A208A60A194A122A171A13A242A39A77A155A165A248A", append([]byte("LOL"), vvb...))
	}
}

func TestInicializa(t *testing.T) {
	indices, asc_codes := inicializa("TESTE123")
	//
	ta := []rune{
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
		'T', 'E', 'S', 'T', 'E', '1', '2', '3', 'T', 'E', 'S', 'T', 'E', '1', '2', '3',
	}
	ti := []int{
		84, 154, 239, 70, 143, 30, 190, 55, 32, 225, 40, 6, 133, 44, 171, 114, 36, 104, 234, 50, 38, 140, 152, 99, 207, 42, 151, 253, 7, 74, 176, 87, 113, 93, 19,
		181, 131, 220, 57, 142, 62, 120, 9, 96, 39, 180, 112, 46, 86, 61, 75, 162, 58, 23, 226, 132, 230, 0, 179, 254, 195, 188, 161, 223, 148, 77, 135, 136, 20,
		85, 48, 144, 17, 246, 51, 166, 60, 192, 10, 219, 231, 235, 211, 121, 191, 34, 22, 185, 160, 118, 240, 232, 115, 206, 94, 90, 111, 209, 106, 182, 186,
		116, 238, 59, 247, 88, 83, 49, 24, 175, 158, 173, 210, 229, 170, 187, 145, 95, 63, 137, 215, 196, 194, 199, 81, 250, 228, 177, 13, 78, 236, 134, 150,
		163, 31, 244, 89, 5, 37, 167, 97, 3, 126, 105, 208, 1, 205, 69, 71, 76, 237, 200, 212, 155, 252, 233, 15, 41, 213, 222, 165, 64, 169, 102, 216, 201,
		128, 248, 47, 8, 174, 241, 119, 28, 98, 217, 127, 141, 227, 139, 18, 52, 4, 130, 79, 65, 91, 80, 183, 25, 11, 123, 149, 224, 26, 33, 53, 255, 27, 29,
		21, 122, 73, 147, 156, 100, 107, 101, 204, 2, 202, 146, 14, 153, 197, 129, 245, 243, 16, 110, 67, 172, 221, 124, 251, 68, 184, 35, 178, 218, 168,
		109, 72, 214, 203, 125, 249, 117, 103, 164, 242, 108, 12, 193, 159, 56, 92, 43, 82, 45, 54, 198, 189, 157, 66, 138,
	}
	//
	if len(ta) != len(asc_codes) {
		t.Fatal("len mismatch", len(ta), len(asc_codes))
	}
	if len(ti) != len(asc_codes) {
		t.Fatal("len mismatch", len(ti), len(indices))
	}
	for k := range ta {
		if ta[k] != asc_codes[k] {
			t.Fatalf("ASC CODE MISMATCH: %v != %v", string(asc_codes[k]), string(ta[k]))
		}
	}
	for k := range ti {
		if ti[k] != indices[k] {
			t.Fatalf("INDICES MISMATCH: %v != %v\n%v", indices[k], ti[k], indices)
		}
	}
}
