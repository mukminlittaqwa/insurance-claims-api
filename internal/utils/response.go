package utils

import "github.com/gin-gonic/gin"

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Message string      `json:"message,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
    c.JSON(200, Response{Success: true, Data: data})
}

func ErrorResponse(c *gin.Context, code int, err string) {
    c.JSON(code, Response{Success: false, Error: err})
}

func PaginatedResponse(c *gin.Context, data interface{}, total int64, page, limit int) {
    c.JSON(200, Response{
        Success: true,
        Data: map[string]interface{}{
            "items": data,
            "pagination": map[string]interface{}{
                "page":  page,
                "limit": limit,
                "total": total,
            },
        },
    })
}