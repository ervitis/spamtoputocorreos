package models

import (
	"fmt"
	"github.com/ervitis/spamtoputocorreos/regtools"
	"regexp"
	"strings"
	"time"
)

const (
	InquiryCategoryEnvios             InquiryCategoryType = "ENVIOS" // at this moment, use this
	InquiryCategoryRegistroWeb        InquiryCategoryType = "REGISTRO_WEB"
	InquiryCategoryDevoluciones       InquiryCategoryType = "DEVOLUCIONES"
	InquiryCategoryInformacionGeneral InquiryCategoryType = "INFORMACION_GENERAL"
	InquiryCategoryPresupuestos       InquiryCategoryType = "PRESUPUESTOS"
	InquiryCategoryDocumentacion      InquiryCategoryType = "DOCUMENTACION"
	InquiryCategoryOtros              InquiryCategoryType = "OTROS"

	InquiryDescriptionNoLlegadaEnvio           InquiryDescriptionType = "NO_LLEGADA_ENVIO"
	InquiryDescriptionEnvioExtraviado          InquiryDescriptionType = "ENVIO_EXTRAVIADO"
	InquiryDescriptionCuandoReciboEnvio        InquiryDescriptionType = "CUANDO_RECIBO_ENVIO"
	InquiryDescriptionPagado                   InquiryDescriptionType = "PAGADO"
	InquiryDescriptionEnvioSigueInspeccion     InquiryDescriptionType = "ENVIO_SIGUE_INSPECCION"
	InquiryDescriptionEnvioRechazadoInspeccion InquiryDescriptionType = "ENVIO_RECHAZADO_INSPECCION"
	InquiryDescriptionEnvioLiberado            InquiryDescriptionType = "ENVIO_LIBERADO"
)

var (
	InquiryCategoryDescriptionData = map[InquiryCategoryType]map[InquiryDescriptionType]string{
		InquiryCategoryEnvios: {
			InquiryDescriptionNoLlegadaEnvio:           "NO HE RECIBIDO LA NOTIFICACIÓN DE LLEGADA DE MI ENVÍO",
			InquiryDescriptionEnvioExtraviado:          "¿MI ENVIO ESTÁ EXTRAVIADO?",
			InquiryDescriptionCuandoReciboEnvio:        "YA HE PAGADO, ¿CUANDO RECIBO MI ENVÍO?",
			InquiryDescriptionPagado:                   "YA HE PAGADO, ¿QUÉ SUCEDE AHORA?",
			InquiryDescriptionEnvioSigueInspeccion:     "MI ENVÍO SIGUE EN INSPECCIÓN",
			InquiryDescriptionEnvioRechazadoInspeccion: "MI ENVÍO HA SIDO RECHAZADO TRAS LA INSPECCIÓN",
			InquiryDescriptionEnvioLiberado:            "MI ENVÍO YA ESTÁ LIBERADO Y NO ME HA LLEGADO",
		},
		InquiryCategoryRegistroWeb:        nil,
		InquiryCategoryDevoluciones:       nil,
		InquiryCategoryInformacionGeneral: nil,
		InquiryCategoryPresupuestos:       nil,
		InquiryCategoryDocumentacion:      nil,
		InquiryCategoryOtros:              nil,
	}

	inquiryRegex = regexp.MustCompile(`(?P<Cat>\w+)-(?P<Desc>\w+)-(?P<Content>.+)`)
)

type (
	InquiryCategoryType    string
	InquiryDescriptionType string

	ContactData struct {
		Name               string
		Phone              string
		Email              string
		Category           InquiryCategoryType
		InquiryDescription string
		InquiryCategory    InquiryDescriptionType
		RefCode            string
		Query              string
		AcceptPrivacy      bool
	}

	StatusTrace struct {
		RefCode  string
		Statuses []*StatusData
	}

	StatusData struct {
		Date   time.Time
		Status string
		Detail string
	}

	Tokens struct {
		Captcha string
		Csrf    string
		Session string
	}

	ReqStatusTraceBody struct {
		Code           string `json:"envio.numEnvio"`
		TokenRecaptcha string `json:"tokenRecaptcha"`
	}

	InquiryBodyData struct {
		Category    InquiryCategoryType
		Description string                 // text of InquiryCategoryDescriptionData
		Type        InquiryDescriptionType // InquiryDescriptionType
		Query       string
	}
)

func (bd *InquiryBodyData) Marshal(payload string) (*InquiryBodyData, error) {
	if strings.TrimSpace(payload) == "" {
		return nil, fmt.Errorf("payload empty")
	}
	data := regtools.GetParams(inquiryRegex, payload)
	bd.Query = strings.TrimSpace(data["Content"])
	bd.Category = InquiryCategoryType(strings.TrimSpace(data["Cat"]))
	bd.Type = InquiryDescriptionType(strings.TrimSpace(data["Desc"]))
	bd.Description = InquiryCategoryDescriptionData[bd.Category][bd.Type]

	if err := bd.Validate(); err != nil {
		return nil, err
	}
	return bd, nil
}

func (bd *InquiryBodyData) Validate() error {
	if bd.Type == "" || bd.Description == "" || bd.Category == "" || bd.Query == "" {
		return fmt.Errorf("some data are empty")
	}
	return nil
}

func (bd *InquiryBodyData) String() string {
	return fmt.Sprintf("%s, %s -> %s", bd.Category, bd.Type, bd.Query)
}

func (c InquiryCategoryType) Value() string {
	return string(c)
}

func (c InquiryCategoryType) All() []string {
	return []string{
		InquiryCategoryEnvios.Value(),
	}
}

func (d InquiryDescriptionType) All() []string {
	return []string{
		InquiryDescriptionNoLlegadaEnvio.Value(),
		InquiryDescriptionEnvioExtraviado.Value(),
		InquiryDescriptionCuandoReciboEnvio.Value(),
		InquiryDescriptionPagado.Value(),
		InquiryDescriptionEnvioSigueInspeccion.Value(),
		InquiryDescriptionEnvioRechazadoInspeccion.Value(),
		InquiryDescriptionEnvioLiberado.Value(),
	}
}

func (d InquiryDescriptionType) Value() string {
	return string(d)
}
