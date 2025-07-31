package services

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AppMestra/mestra-golang/models"
)

type ProcessadorMensagem struct{}

func NewProcessadorMensagem() *ProcessadorMensagem {
	return &ProcessadorMensagem{}
}

func (p *ProcessadorMensagem) ProcessarMensagem(mensagem, autor, dataStr string) ([]models.EntradaJSON, error) {
	dataMensagem, err := time.Parse("02/01/2006", dataStr)
	if err != nil {
		dataMensagem = time.Now()
	}

	linhas := strings.Split(mensagem, "\n")
	var entradas []models.EntradaJSON
	dataReferencia := dataMensagem

	for _, linha := range linhas {
		linha = strings.TrimSpace(linha)
		if linha == "" {
			continue
		}

		if novaData := p.extrairData(linha, dataMensagem); !novaData.IsZero() {
			dataReferencia = novaData
			continue
		}

		if strings.Contains(strings.ToLower(linha), "hoje") {
			dataReferencia = dataMensagem
			linha = p.removerMarcadorTemporal(linha, "hoje")
		} else if strings.Contains(strings.ToLower(linha), "ontem") {
			dataReferencia = dataMensagem.AddDate(0, 0, -1)
			linha = p.removerMarcadorTemporal(linha, "ontem")
		} else if strings.Contains(strings.ToLower(linha), "anteontem") {
			dataReferencia = dataMensagem.AddDate(0, 0, -2)
			linha = p.removerMarcadorTemporal(linha, "anteontem")
		}

		if valor, descricao := p.extrairValorDescricao(linha); valor > 0 {
			entrada := models.EntradaJSON{
				Valor:     valor,
				Descricao: descricao,
				Data:      dataReferencia.Format("2006-01-02"),
			}
			entradas = append(entradas, entrada)
		}
	}

	return entradas, nil
}

func (p *ProcessadorMensagem) extrairData(linha string, anoReferencia time.Time) time.Time {
	regexData := regexp.MustCompile(`(\d{1,2})\s*(de\s*)?(\w+|\d{1,2})`)
	matches := regexData.FindStringSubmatch(strings.ToLower(linha))

	if len(matches) >= 3 {
		dia, _ := strconv.Atoi(matches[1])
		mesStr := matches[3]

		meses := map[string]int{
			"janeiro": 1, "jan": 1,
			"fevereiro": 2, "fev": 2,
			"março": 3, "mar": 3,
			"abril": 4, "abr": 4,
			"maio": 5, "mai": 5,
			"junho": 6, "jun": 6,
			"julho": 7, "jul": 7,
			"agosto": 8, "ago": 8,
			"setembro": 9, "set": 9,
			"outubro": 10, "out": 10,
			"novembro": 11, "nov": 11,
			"dezembro": 12, "dez": 12,
		}

		var mes int
		if mesNum, ok := meses[mesStr]; ok {
			mes = mesNum
		} else if mesNum, err := strconv.Atoi(mesStr); err == nil {
			mes = mesNum
		} else {
			return time.Time{}
		}

		if dia > 0 && dia <= 31 && mes > 0 && mes <= 12 {
			return time.Date(anoReferencia.Year(), time.Month(mes), dia, 0, 0, 0, 0, time.Local)
		}
	}

	return time.Time{}
}

func (p *ProcessadorMensagem) removerMarcadorTemporal(linha, marcador string) string {
	regex := regexp.MustCompile(`(?i)\b` + marcador + `\b`)
	return strings.TrimSpace(regex.ReplaceAllString(linha, ""))
}

func (p *ProcessadorMensagem) extrairValorDescricao(linha string) (float64, string) {
	linha = strings.TrimSpace(linha)
	if linha == "" {
		return 0, ""
	}

	padroes := []string{
		`(\d+)\s+(\d+)\s+(.+)`,
		`(\d+),(\d+)\s+(.+)`,
		`(.+?)\s+(\d+)\s+(\d+)$`,
		`(.+?)\s+(\d+),(\d+)$`,
		`(\d+)\s+(.+)`,
		`(.+?)\s+(\d+)$`,
		`^(\d+)\s+(\d+)$`,
		`^(\d+),(\d+)$`,
		`^(\d+)$`,
	}

	for i, pattern := range padroes {
		regex := regexp.MustCompile(pattern)
		matches := regex.FindStringSubmatch(linha)

		if len(matches) >= 2 {
			switch i {
			case 0, 1:
				if len(matches) >= 4 {
					reais, _ := strconv.Atoi(matches[1])
					centavos, _ := strconv.Atoi(matches[2])
					descricao := strings.TrimSpace(matches[3])
					valor := float64(reais) + float64(centavos)/100
					return valor, descricao
				}
			case 2, 3:
				if len(matches) >= 4 {
					descricao := strings.TrimSpace(matches[1])
					reais, _ := strconv.Atoi(matches[2])
					centavos, _ := strconv.Atoi(matches[3])
					valor := float64(reais) + float64(centavos)/100
					return valor, descricao
				}
			case 4:
				reais, _ := strconv.Atoi(matches[1])
				descricao := strings.TrimSpace(matches[2])
				valor := float64(reais)
				return valor, descricao
			case 5:
				descricao := strings.TrimSpace(matches[1])
				reais, _ := strconv.Atoi(matches[2])
				valor := float64(reais)
				return valor, descricao
			case 6, 7:
				if len(matches) >= 3 {
					reais, _ := strconv.Atoi(matches[1])
					centavos, _ := strconv.Atoi(matches[2])
					valor := float64(reais) + float64(centavos)/100
					return valor, "sem descrição"
				}
			case 8:
				reais, _ := strconv.Atoi(matches[1])
				valor := float64(reais)
				return valor, "sem descrição"
			}
		}
	}

	return 0, ""
}
