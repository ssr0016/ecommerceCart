package controllers

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/ssr0016/ecommerceCart/models"
)

func AddAddress() gin.HandlerFunc {
	//YT 3:27:56
}

func EditHomeAddress() gin.HandlerFunc {
}

func EditWorkAddress() gin.HandlerFunc {
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(404, "Wrong command")
			return
		}

		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully deleted")

	}
}
