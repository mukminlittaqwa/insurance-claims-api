package handlers

import (
    "insurance-claims-api/internal/models"
    "insurance-claims-api/internal/services"
    "insurance-claims-api/internal/utils"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateClaim(svc services.ClaimService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("user_id").(primitive.ObjectID)
        var req models.CreateClaimRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
            return
        }
        claim, err := svc.CreateClaim(userID, req)
        if err != nil {
            utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
            return
        }
        utils.SuccessResponse(c, claim)
    }
}

func GetMyClaims(svc services.ClaimService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("user_id").(primitive.ObjectID)
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
        claims, total, _ := svc.GetMyClaims(userID, page, limit)
        utils.PaginatedResponse(c, claims, total, page, limit)
    }
}

func GetAllClaims(svc services.ClaimService) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Role-Required", "verifier or approver")
        role := c.GetString("role")
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
        claims, total, err := svc.GetAllClaims(role, page, limit)
        if err != nil {
            utils.ErrorResponse(c, http.StatusForbidden, err.Error())
            return
        }
        utils.PaginatedResponse(c, claims, total, page, limit)
    }
}

func GetClaimByID(svc services.ClaimService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("user_id").(primitive.ObjectID)
        role := c.GetString("role")
        id, _ := primitive.ObjectIDFromHex(c.Param("id"))
        claim, err := svc.GetClaimByID(userID, role, id)
        if err != nil {
            utils.ErrorResponse(c, http.StatusNotFound, "claim not found")
            return
        }
        utils.SuccessResponse(c, claim)
    }
}

func UpdateClaim(svc services.ClaimService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("user_id").(primitive.ObjectID)
        id, _ := primitive.ObjectIDFromHex(c.Param("id"))
        var req models.UpdateClaimRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
            return
        }
        if err := svc.UpdateClaim(userID, id, req); err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
            return
        }
        utils.SuccessResponse(c, map[string]string{"message": "claim updated"})
    }
}

func DeleteClaim(svc services.ClaimService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("user_id").(primitive.ObjectID)
        id, _ := primitive.ObjectIDFromHex(c.Param("id"))
        if err := svc.DeleteClaim(userID, id); err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
            return
        }
        utils.SuccessResponse(c, map[string]string{"message": "claim deleted"})
    }
}

func SubmitClaim(svc services.ClaimService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("user_id").(primitive.ObjectID)
        id, _ := primitive.ObjectIDFromHex(c.Param("id"))
        if err := svc.SubmitClaim(userID, id); err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
            return
        }
        utils.SuccessResponse(c, map[string]string{"message": "claim submitted"})
    }
}

func ReviewClaim(svc services.ClaimService) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Role-Required", "verifier")
        if c.GetString("role") != "verifier" {
            utils.ErrorResponse(c, http.StatusForbidden, "only verifier")
            return
        }
        verifierID := c.MustGet("user_id").(primitive.ObjectID)
        id, _ := primitive.ObjectIDFromHex(c.Param("id"))
        note := c.PostForm("note")
        if err := svc.ReviewClaim(verifierID, id, note); err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
            return
        }
        utils.SuccessResponse(c, map[string]string{"message": "claim reviewed"})
    }
}

func ApproveClaim(svc services.ClaimService) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Role-Required", "approver")
        if c.GetString("role") != "approver" {
            utils.ErrorResponse(c, http.StatusForbidden, "only approver")
            return
        }
        approverID := c.MustGet("user_id").(primitive.ObjectID)
        id, _ := primitive.ObjectIDFromHex(c.Param("id"))
        if err := svc.ApproveClaim(approverID, id); err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
            return
        }
        utils.SuccessResponse(c, map[string]string{"message": "claim approved"})
    }
}

func RejectClaim(svc services.ClaimService) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Role-Required", "approver")
        if c.GetString("role") != "approver" {
            utils.ErrorResponse(c, http.StatusForbidden, "only approver")
            return
        }
        approverID := c.MustGet("user_id").(primitive.ObjectID)
        id, _ := primitive.ObjectIDFromHex(c.Param("id"))
        reason := c.PostForm("reason")
        if err := svc.RejectClaim(approverID, id, reason); err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
            return
        }
        utils.SuccessResponse(c, map[string]string{"message": "claim rejected"})
    }
}