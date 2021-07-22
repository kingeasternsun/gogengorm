package taskdata

import (
	"errors"

	"gorm.io/gorm"
)

func AddUserLink(db *gorm.DB, m *UserLink) error {
	return db.Save(m).Error
}

// https://gorm.cn/zh_CN/docs/query.html
// 注意 当使用结构作为条件查询时，GORM 只会查询非零值字段。这意味着如果您的字段值为 0、''、false 或其他 零值，该字段不会被用于构建查询条件
func QueryUserLinksWithPaginate(db *gorm.DB, condition *UserLink, pageSize, currentPage int) ([]*UserLink, int64, error) {
	res := make([]*UserLink, 0)
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

func UpdateUserLinkByFromUser(db *gorm.DB, fromUser string, up map[string]interface{}) (int64, error) {
	if err := db.Where("from_user = ?", fromUser).Updates(up).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func DeleteUserLinkByFromUser(db *gorm.DB, fromUser string) (int64, error) {
	if err := db.Where("from_user = ?", fromUser).Delete(&UserLink{}).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func QueryUserLinksByFromUser(db *gorm.DB, fromUser string) ([]UserLink, error) {
	res := make([]UserLink, 0)
	if err := db.Where("from_user = ?", fromUser).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateUserLinkByToUser(db *gorm.DB, toUser string, up map[string]interface{}) (int64, error) {
	if err := db.Where("to_user = ?", toUser).Updates(up).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func DeleteUserLinkByToUser(db *gorm.DB, toUser string) (int64, error) {
	if err := db.Where("to_user = ?", toUser).Delete(&UserLink{}).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func QueryUserLinksByToUser(db *gorm.DB, toUser string) ([]UserLink, error) {
	res := make([]UserLink, 0)
	if err := db.Where("to_user = ?", toUser).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateUserLinkByFromUserAndToUser(db *gorm.DB, fromUser string, toUser string, up map[string]interface{}) (int64, error) {
	if err := db.Where("from_user = ? And to_user = ? ", fromUser, toUser).Updates(up).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}
func DeleteUserLinkByFromUserAndToUser(db *gorm.DB, fromUser string, toUser string) (int64, error) {
	if err := db.Where("from_user = ? And to_user = ? ", fromUser, toUser).Delete(&UserLink{}).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func GetUserLinkByFromUserAndToUser(db *gorm.DB, fromUser string, toUser string) (UserLink, bool, error) {
	var res UserLink
	if err := db.Where("from_user = ? And to_user = ? ", fromUser, toUser).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, false, nil
		}
		return res, false, err
	}
	return res, false, nil
}

func UpdateUserLinkByID(db *gorm.DB, id int, up map[string]interface{}) (int64, error) {
	if err := db.Where("id = ?", id).Updates(up).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func DeleteUserLinkByID(db *gorm.DB, id int) (int64, error) {
	if err := db.Where("id = ?", id).Delete(&UserLink{}).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func GetUserLinkByID(db *gorm.DB, id int) (UserLink, bool, error) {
	var res UserLink
	if err := db.Where("id = ?", id).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, false, nil
		}
		return res, false, err
	}
	return res, false, nil
}

func UpdateUserLinkByScore(db *gorm.DB, score int, up map[string]interface{}) (int64, error) {
	if err := db.Where("score = ?", score).Updates(up).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func DeleteUserLinkByScore(db *gorm.DB, score int) (int64, error) {
	if err := db.Where("score = ?", score).Delete(&UserLink{}).Error; err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func QueryUserLinksByScore(db *gorm.DB, score int) ([]UserLink, error) {
	res := make([]UserLink, 0)
	if err := db.Where("score = ?", score).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
