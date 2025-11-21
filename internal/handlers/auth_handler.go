package handlers

import (
    "insurance-claims-api/internal/models"
    "insurance-claims-api/internal/services"
    "insurance-claims-api/internal/utils"
    "net/http"

    "github.com/gin-gonic/gin"
)

func Login(authService services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req models.LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
            return
        }

        resp, err := authService.Login(req)
        if err != nil {
            utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
            return
        }

        utils.SuccessResponse(c, resp)
    }
}