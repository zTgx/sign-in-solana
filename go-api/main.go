package main

import (
	"crypto/ed25519"
	"fmt"
	"net/http"
	"time"

	"github.com/akamensky/base58"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	// JWT 密钥，用于签名和验证
	jwtSecretKey = "your_secret_key" // 请确保使用更安全的密钥
)

// 验证 Solana 签名的函数
func verifySolanaSignature(publicAddress, message, signature string) (bool, error) {
	publicKeyBytes, err := base58.Decode(publicAddress)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %v", err)
	}

	signatureBytes, err := base58.Decode(signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature: %v", err)
	}

	// 添加 Solana 消息头
	// messageWithHeader := append([]byte{0xff}, []byte("solana offchain")...)
	// messageWithHeader = append(messageWithHeader, []byte(message)...)

	// 验证签名
	return ed25519.Verify(publicKeyBytes, []byte(message), signatureBytes), nil
}

// 生成 JWT 的函数
func generateJWT(publicAddress string) (string, error) {
	claims := jwt.MapClaims{
		"address": publicAddress,
		"exp":     jwt.TimeFunc().Add(time.Hour * 72).Unix(), // 过期时间设置为72小时
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecretKey))
}

func main() {
	r := gin.Default()

	// 将 CORS 中间件应用到路由
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 或者指定具体的来源
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	r.POST("/api/verify", func(c *gin.Context) {
		var req struct {
			PublicAddress string `json:"publicAddress"`
			Signature     string `json:"signature"`
		}

		// 绑定请求参数
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// 验证签名
		message := "Sign me" // 您可以根据需要更改此消息
		isValid, err := verifySolanaSignature(req.PublicAddress, message, req.Signature)

		if err != nil || !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}

		// 生成 JWT
		token, err := generateJWT(req.PublicAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	r.Run(":8080") // 启动服务，监听8080端口
}
