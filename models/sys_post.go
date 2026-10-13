package models

// SysPost is the RuoYi `sys_post` table (job positions).
type SysPost struct {
	PostID   uint   `gorm:"column:post_id;primaryKey;autoIncrement" json:"postId"`
	PostCode string `gorm:"column:post_code" json:"postCode"`
	PostName string `gorm:"column:post_name" json:"postName"`
	PostSort int    `gorm:"column:post_sort" json:"postSort"`
	Status   string `gorm:"column:status" json:"status"`
	BaseModel
}

// TableName sets the table name of the model.
func (r *SysPost) TableName() string {
	return "sys_post"
}
