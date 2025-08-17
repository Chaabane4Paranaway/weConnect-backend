// auth.controller.go - stub
// Package controllers contient les handlers de routes HTTP pour l'API Auth
//
// Tous les handlers de ce package sont liés à l'authentification et la gestion des utilisateurs.
//
// @tag.name Auth
// @tag.description Opérations liées à l'authentification des utilisateurs.
package controllers

import (
	"go-backend/database"
	"go-backend/dto"
	"go-backend/models"
	"go-backend/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// var users = make(map[string]*model.User) // In-memory DB

// Register godoc
// @Summary Enregistrer un nouvel utilisateur
// @Description Envoie un OTP à l’email fourni
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param register_req body model.RegisterRequest true "User credentials"
// @Success 201 {object} map[string]string
// @Failure 400 string string
// @Router /auth/register [post]
func Register(c *gin.Context) {

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, dto.GlobalError{
			CodeStatus:       http.StatusBadRequest,
			TechnicalMessage: err.Error(),
			Message:          "Invalid JSON",
		})
		log.Println(err)
		return
	}

	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		if existing.Otp == "" {
			c.JSON(http.StatusBadRequest, dto.GlobalError{
				CodeStatus:       http.StatusBadRequest,
				TechnicalMessage: "",
				Message:          "Utilisateur déjà existant",
			})
			log.Println(err)
		}
		{
			otp := utils.GenerateOTP()
			existing.Otp = otp
			existing.OtpExpiresAt = time.Now().Add(5 * time.Minute)
			if err := database.DB.Save(&existing).Error; err != nil {
				c.JSON(http.StatusInternalServerError, dto.GlobalError{
					CodeStatus:       http.StatusInternalServerError,
					TechnicalMessage: err.Error(),
					Message:          "Failed to save user",
				})
				log.Println(err)
				return
			}
			c.JSON(http.StatusCreated, dto.GlobalError{
				CodeStatus:       http.StatusCreated,
				TechnicalMessage: "",
				Message:          "OTP renvoyé avec succès!",
			})
			log.Printf("OTP envoyé à %s", req.Email)
			return
		}
	}

	// fmt.Printf("Password is %s", req.Password)
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			TechnicalMessage: err.Error(),
			Message:          "Failed to hash password",
		})
		log.Println(err)
		return
	}
	log.Printf("ERROR hashed: %s", hashed)

	otp := utils.GenerateOTP()

	user := &models.User{
		Email:        req.Email,
		Password:     hashed,
		Otp:          otp,
		Verified:     false,
		OtpExpiresAt: time.Now().Add(5 * time.Minute),
	}

	// Création de l'utilisateur en base de données
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			TechnicalMessage: err.Error(),
			Message:          "Failed to create user",
		})
		log.Println(err)
		return
	}

	utils.SendMockOTP(req.Email, otp)

	c.JSON(http.StatusCreated, dto.GlobalSuccess{
		CodeStatus: http.StatusCreated,
		Message:    "Utilisateur créé avec succès. OTP envoyé à l'email",
		Data: dto.RegisterResponse{
			Verified: user.Verified,
			Email:    user.Email,
		},
	})
}

// VerifyOTP godoc
// @Summary Verify OTP sent to the user's email.
// @Description Verify the OTP sent to the user's email.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param otp body model.VerifyOTPRequest true "OTP"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Router /auth/verify/{email} [post]
func VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	var user models.User
	if err := database.DB.First(&user, "email = ?", req.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.GlobalError{
			CodeStatus:       http.StatusUnauthorized,
			TechnicalMessage: err.Error(),
			Message:          "User not found",
		})
		log.Println(err)
		return
	}

	if user.Otp != req.OTP {
		c.JSON(http.StatusUnauthorized, dto.GlobalError{
			CodeStatus:       http.StatusUnauthorized,
			TechnicalMessage: "Invalid OTP",
			Message:          "Invalid OTP",
		})
		return
	}

	if time.Now().After(user.OtpExpiresAt) {
		c.JSON(http.StatusBadRequest, dto.GlobalError{
			CodeStatus:       http.StatusBadRequest,
			TechnicalMessage: "OTP expired",
			Message:          "OTP expired",
		})
		return
	}

	user.Verified = true
	user.Otp = ""
	user.OtpExpiresAt = time.Time{}
	database.DB.Save(&user)
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			TechnicalMessage: err.Error(),
			Message:          "Failed to save user",
		})
		log.Println(err)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			TechnicalMessage: err.Error(),
			Message:          "Failed to generate token",
		})
		log.Println(err)
		return
	}
	database.DB.Save(&user)
	c.JSON(http.StatusOK, dto.GlobalSuccess{
		CodeStatus: http.StatusOK,
		Message:    "Utilisateur vérifié avec succès",
		Data: gin.H{
			"verified": user.Verified,
			"email":    user.Email,
			"token":    token,
		},
	})
}

// Login godoc
// @Summary Login user.
// @Description Login user.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param password body model.LoginRequest true "Password"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.GlobalError{
			CodeStatus:       http.StatusBadRequest,
			TechnicalMessage: err.Error(),
			Message:          "Invalid JSON",
		})
		log.Println(err)
		return
	}

	var user models.User
	if err := database.DB.First(&user, "email = ?", req.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.GlobalError{
			CodeStatus:       http.StatusUnauthorized,
			TechnicalMessage: err.Error(),
			Message:          "Invalid credentials",
		})
		log.Println(err)
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, dto.GlobalError{
			CodeStatus:       http.StatusUnauthorized,
			TechnicalMessage: "Password incorrect",
			Message:          "Password incorrect",
		})
		return
	}
	if !user.Verified {
		c.JSON(http.StatusUnauthorized, dto.GlobalError{
			CodeStatus:       http.StatusUnauthorized,
			TechnicalMessage: "Email not verified",
			Message:          "Email not verified",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			TechnicalMessage: err.Error(),
			Message:          "Erreur génération JWT",
		})
		return
	}
	c.JSON(http.StatusOK, dto.GlobalSuccess{
		CodeStatus: http.StatusOK,
		Message:    "Utilisateur connecté avec succès",
		Data: gin.H{
			"id":    user.ID,
			"token": token,
		},
	})
}

// Me godoc
// @Summary Get user information.
// @Description Get user information.
// @Tags Authentication
// @Accept json
// @Produce json
// @Router /auth/me [get]
func Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id missing in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}
