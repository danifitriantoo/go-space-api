package storage

import (
	"context"
	"errors"
	"go-space-api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var planetCollection *mongo.Collection

// Init initializes the storage with the given MongoDB client
func Init(client *mongo.Client) {
	planetCollection = client.Database("go-space-db").Collection("planets")
}

// GetPlanetCollection returns the planet collection for external use (e.g., tests)
func GetPlanetCollection() *mongo.Collection {
	return planetCollection
}
func GetAllPlanets() ([]models.Planet, error) {
	var planets []models.Planet
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := planetCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var planet models.Planet
		if err := cursor.Decode(&planet); err != nil {
			return nil, err
		}
		planets = append(planets, planet)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return planets, nil
}

func GetPlanetById(id string) (models.Planet, bool) {
	var planet models.Planet
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return planet, false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = planetCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&planet)
	if err != nil {
		return planet, false
	}
	return planet, true
}

func AddPlanet(planet models.Planet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	planet.ID = primitive.NewObjectID()
	_, err := planetCollection.InsertOne(ctx, planet)
	return err
}

func UpdatePlanet(updatedPlanet models.Planet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": updatedPlanet.ID}
	update := bson.M{
		"$set": bson.M{
			"name":     updatedPlanet.Name,
			"diameter": updatedPlanet.Diameter,
			"moons":    updatedPlanet.Moons,
			"distance": updatedPlanet.Distance,
		},
	}

	res, err := planetCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("planet not found")
	}
	return nil
}

func DeletePlanet(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := planetCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("planet not found")
	}
	return nil
}
