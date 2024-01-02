package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ssr0016/ecommerceCart/database"
)

type SingedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SERCRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, first_name string, last_name string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SingedDetails{
		Email:      email,
		First_Name: first_name,
		Last_Name:  last_name,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SingedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("SECRET_KEY"))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshClaims).SignedString([]byte("SECRET_KEY"))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SingedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SingedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("SECRET_KEY"), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SingedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is already expired"
		return
	}

	return claims, msg

}

func UpdateAllTokens(signedToken string, signedFreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updatedObj primitive.D

	updatedObj = append(updatedObj, bson.E{Key: "token", Value: signedToken})
	updatedObj = append(updatedObj, bson.E{Key: "refresh_token", Value: signedFreshToken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: updated_at})

	upsert := true

	filter := bson.M{"user_id": userId}

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserData.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updatedObj}}, &opt)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

}
