package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rohitchauraisa1997/calorie-tracker/dbconn"
	"github.com/rohitchauraisa1997/calorie-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate = validator.New()
var entryCollection *mongo.Collection = dbconn.OpenCollection(dbconn.Client, "calories")

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "ping test for calorie tracker!!")
}

func AddEntry(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry models.Entry
	defer cancel()

	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now()
	entry.CreatedAt = &currentTime

	validationErr := validate.Struct(entry)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
	}

	entry.ID = primitive.NewObjectID()
	_, insertErr := entryCollection.InsertOne(ctx, entry)
	if insertErr != nil {
		msg := fmt.Sprintf("entry was not added due to %s", insertErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func GetEntries(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{"softDeletedAt": primitive.Null{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func GetEntryById(c *gin.Context) {
	entryId := c.Query("id")
	docId, _ := primitive.ObjectIDFromHex(entryId)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var entry bson.M

	if err := entryCollection.FindOne(ctx, bson.M{"_id": docId, "softDeletedAt": primitive.Null{}}).Decode(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, entry)
}

func GetEntriesByIngredient(c *gin.Context) {
	ingredient := c.Query("ingredient")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	fmt.Println("Getting entries by ingredient...", ingredient)

	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{"ingredients": ingredient, "softDeletedAt": primitive.Null{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, entries)
}

func UpdateEntry(c *gin.Context) {
	entryId := c.Query("id")
	docID, _ := primitive.ObjectIDFromHex(entryId)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var entry models.Entry

	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(entry)
	if validationErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
		return
	}

	after := options.After
	updateTime := time.Now()
	update := bson.M{
		"dish":        entry.Dish,
		"fat":         entry.Fat,
		"ingredients": entry.Ingredients,
		"calories":    entry.Calories,
		"updatedAt":   &updateTime,
	}

	err := entryCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": docID, "softDeletedAt": primitive.Null{}},
		bson.M{"$set": update},
		&options.FindOneAndUpdateOptions{ReturnDocument: &after},
	).Decode(&entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, entry)
}

func UpdateIngredients(c *gin.Context) {
	entryId := c.Query("id")
	docID, _ := primitive.ObjectIDFromHex(entryId)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var entryWithUpdatedIngredients models.Entry
	type Ingredient struct {
		Ingredients *string `json:"ingredients"`
	}
	var ingredient Ingredient

	if err := c.BindJSON(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	after := options.After
	updateTime := time.Now()
	update := bson.A{
		bson.M{"$set": bson.M{
			"ingredients": ingredient.Ingredients,
			"updatedAt":   &updateTime,
		}}}

	err := entryCollection.FindOneAndUpdate(ctx,
		bson.M{"_id": docID, "softDeletedAt": primitive.Null{}},
		update,
		&options.FindOneAndUpdateOptions{ReturnDocument: &after},
	).Decode(&entryWithUpdatedIngredients)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entryWithUpdatedIngredients)
}

func SoftDeleteEntry(c *gin.Context) {
	entryId := c.Query("id")
	docId, _ := primitive.ObjectIDFromHex(entryId)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var entry models.Entry
	after := options.After
	deleteTime := time.Now()
	update := bson.A{
		bson.M{"$set": bson.M{
			"softDeletedAt": &deleteTime,
		}}}

	err := entryCollection.FindOneAndUpdate(ctx,
		bson.M{"_id": docId, "softDeletedAt": primitive.Null{}},
		update,
		&options.FindOneAndUpdateOptions{ReturnDocument: &after}).Decode(&entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func DeleteEntry(c *gin.Context) {
	entryId := c.Query("id")
	docId, _ := primitive.ObjectIDFromHex(entryId)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var entry models.Entry

	err := entryCollection.FindOneAndDelete(ctx, bson.M{"_id": docId}, &options.FindOneAndDeleteOptions{}).Decode(&entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}
