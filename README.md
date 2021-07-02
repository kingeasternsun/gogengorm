# gogengorm
Go tool to generate gorm code from struct defination
inspired by [gomodifytags](https://github.com/kingeasternsun/gomodifytags)
- generate Update,Delete,Query function for each column.
- generate Update,Delete,Get function for each unique index.

the Query function return multiple rows like this
```go
func QueryUserLinksByFromUser(db *gorm.DB, fromUser string) ([]UserLink, error) {
	res := make([]UserLink, 0)
	if err := db.Where("from_user = ?",fromUser).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
```

the Get function return single row and if exist like this
```go
func GetUserLinkByID(db *gorm.DB, id int) (UserLink,bool, error) {
	var res :=UserLink
	if err := db.Where("id = ?",id).First(&res).Error; err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return res, false, nil
		}
		return res, false, err
	}
	return res, false, nil
}
```

# how to use

first we define a struct like this without tags
```go
type UserLink struct {
	ID       int
	FromUser string
	ToUser   string
	Score    int
}

```
then we can add `gorm` tags manually or  use [gomodifytags](https://github.com/kingeasternsun/gomodifytags) add tags ,we got this

```go
type UserLink struct {
	ID       int    `json:"id,omitempty" form:"id" binding:"required" gorm:"column:id"`
	FromUser string `json:"fromUser,omitempty" form:"fromUser" binding:"required" gorm:"column:from_user"`
	ToUser   string `json:"toUser,omitempty" form:"toUser" binding:"required" gorm:"column:to_user"`
	Score    int    `json:"score,omitempty" form:"score" binding:"required" gorm:"column:score"`
}

```
then we add more options to the gorm tag ,like `type`, `primaryKey` 

```go
type UserLink struct {
	ID       int    `json:"id,omitempty" form:"id" binding:"required" gorm:"column:id;primaryKey"`
	FromUser string `json:"fromUser,omitempty" form:"fromUser" binding:"required" gorm:"column:from_user;uniqueIndex:idx_from_to"`
	ToUser   string `json:"toUser,omitempty" form:"toUser" binding:"required" gorm:"column:to_user;uniqueIndex:idx_from_to""`
	Score    int    `json:"score,omitempty" form:"score" binding:"required" gorm:"column:score"`
}
```

finally we use command below to generate the golang code 
```shell
gogengorm -file ./testdata/user.go -struct UserLink
```

you can also specify your own template file to generate code 
```shell
gogengorm -file ./testdata/user.go -struct UserLink -template youowntemplate
```