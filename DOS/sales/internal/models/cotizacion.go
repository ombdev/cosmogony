package models

type CotizacionItem struct {
	ItemID  int    `json:"item_id" bson:"item_id,omitempty"`
	ProdID  int    `json:"prod_id" bson:"prod_id,omitempty"`
	ProdSku string `json:"prod_sku" bson:"prod_sku,omitempty"`
}

type Cotizacion struct {
	ID            int              `json:"_id" bson:"_id,omitempty"`
	Folio         int              `json:"folio" bson:"folio,omitempty"`
	Tipo          string           `json:"tipo" bson:"tipo,omitempty"`
	Observaciones string           `json:"observaciones" bson:"observaciones,omitempty"`
	Items         []CotizacionItem `json:"items" bson:"items,omitempty"`
}
