package controllers

import (
	"context"
	"time"

	"gin-app/config"
	"gin-app/models"

	"gin-app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func InitUserController(db *mongo.Client) {
	userCollection = config.GetCollection(db, "users")
}

func CreateUser(c *gin.Context) {

	// context.WithTimeout → gives MongoDB 10 seconds to finish the request.
	// If DB hangs or is slow, Go kills the request cleanly.
	// defer cancel() → frees resources after function returns.

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// You're preparing a struct to store the incoming JSON data.
	var user models.User

	// 	Gin tries to read the incoming JSON body.
	// Fills the user struct fields.
	// If request body is invalid → return 400 Bad Request.

	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 	Search MongoDB for a user with the same email.
	// bson.M{"email": user.Email} → search filter like { email: "ck@gmail.com" }
	// If a user is found → Decode() succeeds → no error.

	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)

	if err == nil {
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}

	// 	bcrypt hashes the plaintext password.
	// Makes it safe to store in DB.
	// user.Password is replaced with the hashed version.
	// Never store plain passwords. This is the correct standard.

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	res, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"insertedId": res.InsertedID})
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var body LoginBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
	}

	email := body.Email
	password := body.Password

	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)

	if err != nil {
		c.JSON(404, gin.H{"error": "user does not exits."})
		return
	}

	result := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))

	if result != nil {
		c.JSON(404, gin.H{"error": "invalid creedentials"})
		return
	}

	token, err := utils.GenerateToken(existingUser.ID.Hex(), existingUser.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.SetCookie(
		"token",
		token,
		3600*24,
		"/",
		"",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"message": "login success",
		"token":   token,
	})
}

type PublicUser struct {
	ID    primitive.ObjectID `json:"_id"`
	Email string             `json:"email"`
	Name  string             `json:"username"`
}

func UserProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get userId from middleware
	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	// Convert string ID → ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	// Fetch user from DB
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)

	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	publicUser := PublicUser{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	c.JSON(200, gin.H{
		"message": "profile fetched",
		"user":    publicUser,
	})
}
