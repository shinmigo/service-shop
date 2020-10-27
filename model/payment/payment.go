package payment

type Payment struct {
	Id     uint64 `json:"id" gorm:"PRIMARY_KEY"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Params string `json:"params"`
	Status int8   `json:"status"`
	Sort   uint64 `json:"sort"`
}

func GetTableName() string {
	return "payment"
}
