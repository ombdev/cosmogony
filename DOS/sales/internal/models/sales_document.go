package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Prospecto struct {
	RazonSocial string `json:"razon_social" bson:"razon_social,omitempty"`
	Domicilio   string `json:"domicilio" bson:"domicilio,omitempty"`
	Tel         string `json:"tel" bson:"tel,omitempty"`
	Email       string `json:"email" bson:"email,omitempty"`
	Contacto    string `json:"contacto" bson:"contacto,omitempty"`
}

type Impuesto struct {
	Tasa  float32 `json:"tasa" bson:"tasa,omitempty"`
	Monto float32 `json:"monto" bson:"monto,omitempty"`
}

type SalesDocumentItem struct {
	ID                 primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	SalesDocumentID    primitive.ObjectID `json:"sales_document_id" bson:"sales_document_id,omitempty"`
	ProdID             int                `json:"prod_id" bson:"prod_id,omitempty"`
	ProdCantidad       float32            `json:"prod_cantidad" bson:"prod_cantidad,omitempty"`
	ProdPrecioUnitario float32            `json:"prod_precio_unitario" bson:"prod_precio_unitario,omitempty"`
	ProdImporte        float32            `json:"prod_importe" bson:"prod_importe,omitempty"`
	ProdTraslados      []Impuesto         `json:"prod_traslados" bson:"prod_traslados,omitempty"`
	ProdRetenciones    []Impuesto         `json:"prod_retenciones" bson:"prod_retenciones,omitempty"`
	ProdStatusAut      string             `json:"prod_status_aut" bson:"prod_status_aut,omitempty"`
	ProdPrecioAut      float32            `json:"prod_precio_aut" bson:"prod_precio_aut,omitempty"`
	ProdGralUsrIDAut   int                `json:"prod_gral_usr_id_aut" bson:"prod_gral_usr_id_aut,omitempty"`
	ProdRequiereAut    bool               `json:"prod_requiere_aut" bson:"prod_requiere_aut"`
}

type SalesDocument struct {
	ID                 primitive.ObjectID  `json:"_id" bson:"_id,omitempty"`
	Tipo               string              `json:"tipo" bson:"tipo,omitempty"`
	ClienteID          int                 `json:"cliente_id" bson:"cliente_id,omitempty"`
	Prospecto          Prospecto           `json:"prospecto" bson:"prospecto,omitempty"`
	Observaciones      string              `json:"observaciones" bson:"observaciones,omitempty"`
	Fecha              string              `json:"fecha" bson:"fecha,omitempty"`
	MonedaDocumento    string              `json:"moneda_documento" bson:"moneda_documento,omitempty"`
	TipoCambio         float32             `json:"tipo_cambio" bson:"tipo_cambio,omitempty"`
	MonedaBase         string              `json:"moneda_base" bson:"moneda_base,omitempty"`
	AgenteVentas       string              `json:"agente_ventas" bson:"agente_ventas,omitempty"`
	DiasVigencia       int                 `json:"dias_vigencia" bson:"dias_vigencia,omitempty"`
	TrasladosEnabled   bool                `json:"traslados_enabled" bson:"traslados_enabled"`
	RetencionesEnabled bool                `json:"retenciones_enabled" bson:"retenciones_enabled"`
	MontoSubtotal      float32             `json:"monto_subtotal" bson:"monto_subtotal,omitempty"`
	MontoTraslados     float32             `json:"monto_traslados" bson:"monto_traslados,omitempty"`
	MontoRetenciones   float32             `json:"monto_retenciones" bson:"monto_retenciones,omitempty"`
	MontoTotal         float32             `json:"monto_total" bson:"monto_total,omitempty"`
	Items              []SalesDocumentItem `json:"items" bson:"items,omitempty"`
}
