package storage

import (
	"context"
	"fmt"
	"log"

	"cosmogony.com/sales/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateSalesDocument(salesDoc *models.SalesDocument) (string, error) {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("cosmogony")

	items := salesDoc.Items
	// Para omitir el atributo Items en docto de MongoDB, se le asigna el zero value de un slice: nil
	salesDoc.Items = nil

	salesDocumentColl := db.Collection("sales_document")

	result, err := salesDocumentColl.InsertOne(ctx, salesDoc)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("Error al intentar insercion de un documento de ventas: %v", err)
	}

	log.Println("Successfully inserted", *result, *salesDoc)
	fmt.Println("Successfully inserted", *result, *salesDoc)
	fmt.Printf("*****************(decoded from JSON: %+v\n", *salesDoc)

	salesDocumentItemColl := db.Collection("sales_document_item")

	// result2, err := salesDocumentItemColl.InsertMany(ctx, salesDoc.Items)
	// if err != nil {
	// 	log.Println(err)
	// 	return "", fmt.Errorf("Error al intentar insercion de los items de un documento de ventas: %v", err)
	// }

	var docs []interface{}
	for _, v := range items {
		docs = append(docs, bson.D{
			{"_id", primitive.NewObjectID()},
			{"sales_document_id", result.InsertedID},
			{"prod_id", v.ProdID},
			{"prod_cantidad", v.ProdCantidad},
			{"prod_precio_unitario", v.ProdPrecioUnitario},
			{"prod_importe", v.ProdImporte},
			{"prod_traslados", v.ProdTraslados},
			{"prod_retenciones", v.ProdRetenciones},
			{"prod_status_aut", v.ProdStatusAut},
			{"prod_precio_aut", v.ProdPrecioAut},
			{"prod_gral_usr_id_aut", v.ProdGralUsrIDAut},
			{"prod_requiere_aut", v.ProdRequiereAut},
		})
	}
	result2, err := salesDocumentItemColl.InsertMany(ctx, docs)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("Error al intentar insercion de los items de un documento de ventas: %v", err)
	}
	log.Println("Successfully inserted", *result2)
	fmt.Println("Successfully inserted", *result2)
	fmt.Printf("*****************(decoded from JSON: %+v\n", result2.InsertedIDs)

	return fmt.Sprintf("%+v", *result), nil
}

func ReadSalesDocument(objID string) (*models.SalesDocument, error) {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("cosmogony")
	objectID, err := primitive.ObjectIDFromHex(objID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	salesDocumentColl := db.Collection("sales_document")
	filter := bson.D{{"_id", objectID}}
	salesDoc := models.SalesDocument{}

	if err = salesDocumentColl.FindOne(ctx, filter).Decode(&salesDoc); err != nil {
		fmt.Println(err)
		return nil, err
	}

	salesDocumentItemColl := db.Collection("sales_document_item")
	filter = bson.D{{"sales_document_id", objectID}}

	cur, err := salesDocumentItemColl.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var item models.SalesDocumentItem
		if err = cur.Decode(&item); err != nil {
			fmt.Println(err)
			return nil, err
		}
		salesDoc.Items = append(salesDoc.Items, item)
	}

	fmt.Printf("*****************(decoded from MongoDB): %+v\n", salesDoc)
	return &salesDoc, nil
}

func UpdateSalesDocument(id int, salesDoc *models.SalesDocument) {

}

func DeleteSalesDocument(id int) {

}

func ReadSalesDocumentList(filtros []string) {

}
