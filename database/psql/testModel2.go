package psql

import (
	"time"
)

type App struct {
	ID          string    `gorm:"column:id;type:varchar;PRIMARY_KEY "form:"id" json:"id"`
	Name        string    `gorm:"column:name;type:varchar; "form:"name" json:"name"`
	Url         string    `gorm:"column:url;type:varchar; "form:"url" json:"url"`
	Icon        string    `gorm:"column:icon;type:varchar; "form:"icon" json:"icon"`
	Explain     string    `gorm:"column:explain;type:varchar; "form:"explain" json:"explain"`
	Secret      string    `gorm:"column:secret;type:varchar; "form:"secret" json:"secret"`
	Type        string    `gorm:"column:type;type:varchar; "form:"type" json:"type"`
	UseStatus   string    `gorm:"column:use_status;type:varchar; "form:"use_status" json:"use_status"` //使用状态0：正常，-1删除
	RedirectURL string    `gorm:"column:redirect_url;type:varchar; "form:"redirect_url" json:"redirect_url"`
	Scope       string    `gorm:"column:scope;type:varchar; "form:"scope" json:"scope"`
	CreateTime  time.Time `gorm:"column:create_time;type:timestamp; "form:"create_time" json:"create_time"`
}

//TableName 设置表名
func (App) TableName() string {
	return "t_app"
}

/*
Insert 插入数据
*/
func (o *App) Insert() error {
	tx := DB.Begin()
	err := tx.Create(o).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Delete : 删除数据
func (o *App) Delete() error {
	tx := DB.Begin()
	// err := tx.Where("id=?", o.ID).Delete(o).Error
	err := tx.Delete(o).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (o *App) Cancel() error {
	tx := DB.Begin()
	err := tx.Where(" secret=?", o.Secret).Delete(o).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

/*
UpdateByID 修改数据
*/
func (o *App) UpdateByID(app App) error {
	tx := DB.Begin()
	err := tx.Model(o).Updates(app).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

/*
SelectList 查询列表
*/
func (o *App) SelectList(limit, offset int) (os []App, count int) {
	db := DB
	os = make([]App, 0)
	sql := "from t_app "

	if o.Name != "" || o.Type != "" {
		sql = sql + " where "
	}
	if o.Name != "" {
		sql = sql + " name LIKE '%" + o.Name + "%' "
	}
	if o.Type != "" {
		if o.Name != "" {
			sql = sql + " and type LIKE '%" + o.Type + "%' "
		} else {
			sql = sql + " type LIKE '%" + o.Type + "%' "
		}
	}

	db.Raw("select count(*) " + sql).Count(&count)
	db.Raw("select * " + sql).Limit(limit).Offset(offset).Scan(&os)
	return
}

/*
SelectAll 查询全部
*/
func (o *App) SelectAll() (os []App) {
	db := DB
	os = make([]App, 0)
	sql := "select * from t_app "
	if o.Name != "" {
		sql = sql + " where name LIKE '%" + o.Name + "%'"
	}

	if o.Name != "" || o.Type != "" {
		sql = sql + " where "
	}
	if o.Name != "" {
		sql = sql + " name LIKE '%" + o.Name + "%' "
	}
	if o.Type != "" {
		if o.Name != "" {
			sql = sql + " and type LIKE '%" + o.Type + "%' "
		} else {
			sql = sql + " type LIKE '%" + o.Type + "%' "
		}
	}

	db.Raw(sql).Scan(&os)
	return
}

/*
SelectList 查询所有状态
*/
func (o *App) SelectListByStatus(useStatus string) (os []App) {
	db := DB
	os = make([]App, 0)
	sql := "select * from t_app "
	if useStatus != "" {
		sql = sql + " where use_status='" + useStatus + "'"
	}
	db.Raw(sql).Scan(&os)
	return
}

/*
SelectByID 根据id查询
*/
func (o *App) SelectByID() (p App) {
	db := DB
	db.Where("id=?", o.ID).Find(&p)
	return
}

// SelectAccountCollection : 获取当前用户的所有收藏
func (o *App) SelectAccountCollection(accountID string, limit, offset int) (os []App, count int) {
	db := DB
	os = make([]App, 0)
	sql := `from t_account 
			inner join t_collection on t_collection.account_id = t_account.id
			right join t_app on t_app.id = t_collection.app_id
			where t_account.id=? `

	if o.Name != "" {
		sql = sql + "and t_app.name LIKE '%" + o.Name + "%' "
	}
	if o.Type != "" {
		sql = sql + " and t_app.type LIKE '%" + o.Type + "%' "
	}

	db.Raw("select count(t_app.*) "+sql, accountID).Count(&count)
	db.Raw("select t_app.* "+sql+"order by t_app.name", accountID).Limit(limit).Offset(offset).Scan(&os)
	return
}

// SelectAccountApp : 获取当前用户的所有App, 可按类型查询。
func (o *App) SelectAccountApp(accountID string, limit, offset int) (os []App, count int) {
	db := DB
	os = make([]App, 0)
	sql := `from t_account 
			inner join t_account_to_group on t_account_to_group.account_id = t_account.id
			inner join t_app_to_group on t_app_to_group.group_id = t_account_to_group.group_id
			inner join t_app on t_app.id = t_app_to_group.app_id
			where t_account.id=? `

	if o.Name != "" {
		sql = sql + "and t_app.name LIKE '%" + o.Name + "%' "
	}
	if o.Type != "" {
		sql = sql + " and t_app.type LIKE '%" + o.Type + "%' "
	}

	// db.Raw("select count(t_app.*) "+sql, accountID).Count(&count)
	// db.Raw("select t_app.* "+sql, accountID).Limit(limit).Offset(offset).Scan(&os)
	db.Raw("select count(distinct t_app.*) "+sql, accountID).Count(&count)
	db.Raw("select distinct t_app.* "+sql+"order by t_app.name", accountID).Limit(limit).Offset(offset).Scan(&os)
	return
}

// 测试用，保留写法。Limit(-1)可以去掉前面的Limit，可以Find或者Scan后取消分页，然后获取count。但是又2个问题：
// 1. 仍然时查询2次，没有优化。
// 2. 查询总数只能默认count(*)
// func (o *App) SelectAccountApp(accountID string, limit, offset int) (os []App, count int) {
// 	db := DB
// 	os = make([]App, 0)
// 	db.Table("t_account").Select("distinct t_app.*").Joins(
// 		"inner join t_account_to_group on t_account_to_group.account_id = t_account.id",
// 	).Joins(
// 		"inner join t_app_to_group on t_app_to_group.group_id = t_account_to_group.group_id",
// 	).Joins(
// 		"inner join t_app on t_app.id = t_app_to_group.app_id",
// 	).Where("t_account.id=? ", accountID).Order("t_app.name").Limit(limit).Offset(offset).Find(&os).Limit(-1).Offset(-1).Count(&count)
// 	return
// }

// AppTypeList : 查询所有的应用类型。
func (o *App) AppTypeList() []string {
	db := DB
	apps := make([]App, 0)
	sql := `select distinct t_app.type from t_app
			where type != '' `

	db.Raw(sql).Find(&apps)

	types := make([]string, 0)
	for _, app := range apps {
		types = append(types, app.Type)
	}

	return types
}

// SelectByGroupID : 查询所有应用，按用户组.
func (o *App) SelectByGroupID(groupID string, limit, offset int) (os []App, count int) {
	db := DB
	os = make([]App, 0)
	sql := `from t_app
			inner join t_app_to_group on t_app_to_group.app_id = t_app.id
			where t_app_to_group.group_id = ? `

	if o.Name != "" {
		sql = sql + "and t_app.name LIKE '%" + o.Name + "%' "
	}
	if o.Type != "" {
		sql = sql + " and t_app.type LIKE '%" + o.Type + "%' "
	}

	db.Raw("select count(t_app.*) "+sql, groupID).Count(&count)
	db.Raw("select t_app.* "+sql+"order by t_app.name", groupID).Limit(limit).Offset(offset).Scan(&os)
	return
}
