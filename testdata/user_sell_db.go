package taskdata

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func AddUserSell(db *gorm.DB, m *UserSell) error {
	return db.Save(m).Error
}

// https://gorm.cn/zh_CN/docs/query.html
// 注意 当使用结构作为条件查询时，GORM 只会查询非零值字段。这意味着如果您的字段值为 0、''、false 或其他 零值，该字段不会被用于构建查询条件
func QueryUserSellsWithPaginate(db *gorm.DB, condition *UserSell, pageSize, currentPage int) ([]*UserSell, int64, error) {
	res := make([]*UserSell, 0)
    var totalCount int64
    err = db.Where(condition).
        Offset((currentPage - 1) * pageSize).
		Find(&res).
        Offset(-1).
        Count(&totalCount).Error
	if err != nil {
		return res, totalCount, err
	}
	return res, totalCount, nil
}
 

func UpdateUserSellByNum(db *gorm.DB, num int, up map[string]interface{}) (int64, error) {
	if err := db.Where("num = ?", num).Updates(up).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func DeleteUserSellByNum(db *gorm.DB, num int) (int64, error) {
	if err := db.Where("num = ?", num).Delete(&UserSell{}).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func QueryUserSellsByNum(db *gorm.DB, num int) ([]UserSell, error) {
	res := make([]UserSell, 0)
	if err := db.Where("num = ?", num).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}



func UpdateUserSellByIDAndItem(db *gorm.DB, id int, item string, up map[string]interface{}) (int64, error) {
	if err := db.Where("id = ? And item = ? ", id, item).Updates(up).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}
func DeleteUserSellByIDAndItem(db *gorm.DB, id int, item string) (int64, error) {
	if err := db.Where("id = ? And item = ? ", id, item).Delete(&UserSell{}).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func QueryUserSellsByIDAndItem(db *gorm.DB, id int, item string) ([]UserSell, error) {
	res := make([]UserSell, 0)
	if err := db.Where("id = ? And item = ? ", id, item).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil

}

// https://gorm.cn/zh_CN/docs/query.html
// 注意 当使用结构作为条件查询时，GORM 只会查询非零值字段。这意味着如果您的字段值为 0、''、false 或其他 零值，该字段不会被用于构建查询条件
func QueryUserSellsByIDAndItemWithPaginate(db *gorm.DB, id int, item string, pageSize, currentPage int) ([]*UserSell, int64, error) {
	res := make([]*UserSell, 0)
    var totalCount int64
    err = db.Where("id = ? And item = ? ", id, item).
        Offset((currentPage - 1) * pageSize).
		Find(&res).
        Offset(-1).
        Count(&totalCount).Error
	if err != nil {
		return res, totalCount, err
	}
	return res, totalCount, nil
}
 
 