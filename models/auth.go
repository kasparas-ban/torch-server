package models

// import (
// 	"fmt"
// 	"net/http"
// 	"os"

// 	jwt "github.com/golang-jwt/jwt/v4"
// )

// type TokenDetails struct {
// 	AccessToken  string
// 	RefreshToken string
// 	AccessUUID   string
// 	RefreshUUID  string
// 	AtExpires    int64
// 	RtExpires    int64
// }

// type AccessDetails struct {
// 	AccessUUID string
// 	UserID     int64
// }

// type Token struct {
// 	AccessToken  string `json:"access_token"`
// 	RefreshToken string `json:"refresh_token"`
// }

// //AuthModel ...x
// type AuthModel struct{}

// //VerifyToken ...
// func (m AuthModel) VerifyToken(r *http.Request) (*jwt.Token, error) {
// 	tokenString := m.ExtractToken(r)
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		//Make sure that the token method conform to "SigningMethodHMAC"
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("ACCESS_SECRET")), nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return token, nil
// }
