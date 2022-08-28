package admin

type AccountService struct {
	AccountName     string `json:"accountName" binding:"required"`
	AccountPassword string `json:"accountPassword" binding:"required"`
}
