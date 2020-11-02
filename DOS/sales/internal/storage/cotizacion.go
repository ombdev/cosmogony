package storage

import (
	"context"
	"fmt"
	"log"

	"cosmogony.com/sales/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateCotizacion(cotizacion *models.Cotizacion) (string, error) {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	cotizCollection := client.Database("cosmogony").Collection("cotizaciones")

	result, err := cotizCollection.InsertOne(ctx, cotizacion)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("Error al intentar insercion de una cotizacion: %v", err)
	}

	log.Println("Successfully inserted", *result, *cotizacion)
	fmt.Println("Successfully inserted", *result, *cotizacion)
	fmt.Printf("*****************(decoded from JSON: %+v\n", *cotizacion)
	return fmt.Sprintf("%+v", *result), nil
}

func ReadCotizacion(id int) (*models.Cotizacion, error) {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	cotizCollection := client.Database("cosmogony").Collection("cotizaciones")
	filter := bson.D{{"_id", id}}
	cotizacion := models.Cotizacion{}

	if err = cotizCollection.FindOne(ctx, filter).Decode(&cotizacion); err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("*****************(decoded from MongoDB): %+v\n", cotizacion)
	return &cotizacion, nil
}

func UpdateCotizacion(id int, cot *models.Cotizacion) {

}

func DeleteCotizacion(id int) {

}

func ReadCotizaciones(filtros []string) {

}
