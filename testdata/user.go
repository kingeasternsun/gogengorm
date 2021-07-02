package taskdata

type User struct {
	UserName string `json:"userName,omitempty" form:"userName" binding:"required" gorm:"column:user_name;type:varchar(64);primaryKey"`
	UserAge  int    `json:"userAge,omitempty" form:"userAge" binding:"required" gorm:"type:smallint;column:user_age"`
	Email    string `json:"email,omitempty" form:"email" binding:"required" gorm:"column:email;type:varchar(64);unique"`
}

type UserScore struct {
	UserName string `json:"userName,omitempty" form:"userName" binding:"required" gorm:"column:user_name;type:varchar(64);primaryKey"`
	ClassID  string `json:"class_id,omitempty" form:"class_id" binding:"required" gorm:"type:varchar(64);column:primaryKey"`
	Score    int    `json:"score,omitempty" form:"score" binding:"required" gorm:"type:smallint;column:score"`
}

type UserLink struct {
	ID       int    `json:"id,omitempty" form:"id" binding:"required" gorm:"column:id;primaryKey"`
	FromUser string `json:"fromUser,omitempty" form:"fromUser" binding:"required" gorm:"column:from_user;uniqueIndex:idx_from_to"`
	ToUser   string `json:"toUser,omitempty" form:"toUser" binding:"required" gorm:"column:to_user;uniqueIndex:idx_from_to"`
	Score    int    `json:"score,omitempty" form:"score" binding:"required" gorm:"column:score"`
}
