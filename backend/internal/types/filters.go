package types

type UserFilter struct {
    Status string `form:"status"`
    City   string `form:"city"`
    Role   string `form:"role"`
}