package psql

import "time"

type Account struct {
	ID             string         `gorm:"column:id;type:varchar;PRIMARY_KEY "form:"id" json:"id"`
	Name           string         `gorm:"column:name;type:varchar;unique"json:"name";`
	Password       string         `gorm:"column:password;type:varchar;"json:"password"`
	UseStatus      string         `gorm:"column:use_status;type:varchar; "form:"use_status" json:"use_status"` //使用状态0：正常，-1删除
	Type           string         `gorm:"column:type;type:varchar"json:"type"`
	Mail           string         `gorm:"column:mail;type:varchar"json:"mail"`
	Phone          string         `gorm:"column:phone;type:varchar"json:"phone"`
	Orgs           []Org          `gorm:"many2many:t_account_to_org;" json:"orgs"`
	Groups         []AccountGroup `gorm:"many2many:t_account_to_group;" json:"groups"`
	CreateTime     time.Time      `gorm:"column:create_time;type:timestamp; "form:"create_time" json:"create_time"`
	Apps           []App          `gorm:"-" json:"apps"` //`gorm:"many2many:t_account_group;"`
	CollectionApps []App          `gorm:"many2many:t_collection;" json:"collections"`
}

//TableName 设置表名
func (Account) TableName() string {
	return "t_account"
}

/*
SelectAllByID 根据ID查询用户全部数据, 备用方法，不建议使用。
*/
// ti 1 不查询收藏  2 不查询app 3 收藏和app都不查询. 其他字段都表示正常查询。
func (o *Account) SelectAllByID(ti string) *Account {
	db := DB
	account := Account{}
	// db.Raw("SELECT * from t_account WHERE id=? ", o.ID).Scan(&user)
	// db.Table("t_account").Select("t_account.*, t_app.*, t_account_group.*, t_organization.*").Joins("left join t_account_group on t_account_group.id = t_account.group_id").Joins("left join t_organization on t_organization.id = t_account.org_id").Joins("left join t_collection on t_collection.account_id = ?", o.ID).Joins("left join t_app on t_app.id = t_collection.app_id").Where("t_account.id=?", o.ID).Related(&account)

	// db.Table("t_account").Select("*").Joins(
	// 	"left join t_account_to_group on t_account_to_group.account_id = ?", o.ID,
	// ).Joins(
	// 	"left join t_account_to_org on t_account_to_org.account_id = ?", o.ID,
	// ).Joins(
	// 	"left join t_collection on t_collection.account_id = ?", o.ID,
	// ).Joins(
	// 	"left join t_app on t_app.id = t_collection.app_id",
	// ).Where("t_account.id=?", o.ID).Find(&account)

	// db.Raw(`SELECT * from t_account
	// left join t_collection on t_collection.account_id = ?
	// left join t_app on t_app.id = t_collection.app_id
	//  where t_account.id=? `, o.ID, o.ID).Scan(&account)

	// 这里可以一个使用join查询多对多的数据，但是问题是 gorm 多对多模型取值取不到，这里先做多次查询， 后面需要优化。 实在无法在gorm中实现就写view吧, view也有问题，没有分类。
	db.Table("t_account").Select("*").Where("t_account.id=?", o.ID).Find(&account)

	// 获取分组
	var group []AccountGroup
	db.Table("t_account").Select("t_account_group.*").Joins(
		"inner join t_account_to_group on t_account_to_group.account_id = t_account.id",
	).Joins(
		"right join t_account_group on t_account_group.id = t_account_to_group.group_id",
	).Where("t_account.id=?", o.ID).Scan(&group)

	account.Groups = group

	// 获取机构
	var org []Org
	db.Table("t_account").Select("t_organization.*").Joins(
		"inner join t_account_to_org on t_account_to_org.account_id = t_account.id",
	).Joins(
		"right join t_organization on t_organization.id = t_account_to_org.org_id",
	).Where("t_account.id=?", o.ID).Find(&org)

	account.Orgs = org

	if ti != "3" && ti != "1" {
		// 二次查询所有收藏的用户, 这里多一次数据库查询会多消耗 200多毫秒
		var apps []App
		db.Table("t_account").Select("t_app.*").Joins(
			"inner join t_collection on t_collection.account_id = t_account.id",
		).Joins(
			"right join t_app on t_app.id = t_collection.app_id",
		).Where("t_account.id=?", o.ID).Order("t_app.name").Scan(&apps)

		account.CollectionApps = apps
	}

	if ti != "3" && ti != "2" {
		// 二次查询用户所有应用, 这里多一次数据库查询会多消耗 200多毫秒
		var apps2 []App
		db.Table("t_account").Select("distinct t_app.*").Joins(
			"inner join t_account_to_group on t_account_to_group.account_id = t_account.id",
		).Joins(
			"inner join t_app_to_group on t_app_to_group.group_id = t_account_to_group.group_id",
		).Joins(
			"inner join t_app on t_app.id = t_app_to_group.app_id",
		).Where("t_account.id=?", o.ID).Order("t_app.name").Scan(&apps2)

		account.Apps = apps2
	}

	return &account
}

/*
Insert 插入数据
*/
func (o *Account) Insert() error {
	tx := DB.Begin()
	err := tx.Create(o).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

/*
UpdatePwdByAccount 更新密码
*/
func (o *Account) UpdatePwdByAccount(newPassword string) error {
	tx := DB.Begin()
	err := tx.Exec("UPDATE t_account set password=? WHERE (name=? OR phone=? OR mail=?) and password=?", newPassword, o.Name, o.Name, o.Name, o.Password).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

/*
UpdateStatusByAccount 更新用户状态
*/
func (o *Account) UpdateStatusByAccount() error {
	tx := DB.Begin()
	err := tx.Exec("UPDATE t_account set use_status=? WHERE name = ?", o.UseStatus, o.Name).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

/*
SelectByAccount 根据ID查询用户数据
*/
func (o *Account) SelectByID() *Account {
	db := DB
	user := Account{}
	// db.Raw("SELECT * from t_account WHERE id=? ", o.ID).Scan(&user)
	db.Where("id=?", o.ID).Find(&user)
	return &user
}

/*
SelectByName 根据用户名查询用户数据
*/
func (o *Account) SelectByName() {
	db := DB
	db.Raw("SELECT * from t_account WHERE name=?", o.Name).Scan(&o)
}

/*
SelectByAccountPwd 根据用户名和密码查下用户信息
*/
func (o *Account) SelectByAccountPwd() *Account {
	db := DB
	user := Account{}
	db.Raw("SELECT * from t_account WHERE (name=? OR phone=? OR mail=?) AND password=? ", o.Name, o.Name, o.Name, o.Password).Scan(&user)
	return &user
}

/*
SelectByAccountPwd 根据用户名和密码和类型查下用户信息，主要用户管理员登录。
*/
func (o *Account) SelectByAccountPwdType() *Account {
	db := DB
	user := Account{}
	db.Raw("SELECT * from t_account WHERE (name=? OR phone=? OR mail=?) AND password=? AND type=? ", o.Name, o.Name, o.Name, o.Password, o.Type).Scan(&user)
	return &user
}

/*
SelectList 查询列表
*/
func (o *Account) SelectListByStatus(useStatus string) (account []Account) {
	db := DB
	account = make([]Account, 0)
	sql := "select * from t_account "
	if useStatus != "" {
		sql = sql + " where use_status='" + useStatus + "'"
	}
	db.Raw(sql).Scan(&account)
	return
}

/*
SelectList 查询列表
*/
func (o *Account) SelectList(limit, offset int) (os []Account, count int) {
	db := DB
	os = make([]Account, 0)
	sql := "from t_account "
	// 这里处理查询项比较粗暴，要是有其他的添加的查询项就修改掉。
	if o.Name != "" || o.Mail != "" || o.Phone != "" {
		sql = sql + " where "
	}
	if o.Name != "" {
		sql = sql + " name LIKE '%" + o.Name + "%' "
	}
	if o.Mail != "" {
		if o.Name != "" {
			sql = sql + " and mail LIKE '%" + o.Mail + "%' "
		} else {
			sql = sql + " mail LIKE '%" + o.Mail + "%' "
		}
	}
	if o.Phone != "" {
		if o.Name != "" || o.Mail != "" {
			sql = sql + " and phone LIKE '%" + o.Phone + "%' "
		} else {
			sql = sql + " phone LIKE '%" + o.Phone + "%' "
		}
	}

	db.Raw("select count(*) " + sql).Count(&count)
	db.Raw("select * " + sql).Limit(limit).Offset(offset).Scan(&os)
	// db.Where(o).Limit(limit).Offset(offset).Scan(&os)
	return
}

/*
SelectAll 查询全部
*/
func (o *Account) SelectAll() (os []Account) {
	db := DB
	os = make([]Account, 0)
	sql := "select * from t_account "
	// 这里处理查询项比较粗暴，要是有其他的添加的查询项就修改掉。
	if o.Name != "" || o.Mail != "" || o.Phone != "" {
		sql = sql + " where "
	}
	if o.Name != "" {
		sql = sql + " name LIKE '%" + o.Name + "%' "
	}
	if o.Mail != "" {
		if o.Name != "" {
			sql = sql + " and mail LIKE '%" + o.Mail + "%' "
		} else {
			sql = sql + " mail LIKE '%" + o.Mail + "%' "
		}
	}
	if o.Phone != "" {
		if o.Name != "" || o.Mail != "" {
			sql = sql + " and phone LIKE '%" + o.Phone + "%' "
		} else {
			sql = sql + " phone LIKE '%" + o.Phone + "%' "
		}
	}

	db.Raw(sql).Scan(&os)
	return
}

/*
UpdateByID 修改数据
*/
func (o *Account) UpdateByID(account Account) error {
	tx := DB.Begin()
	err := tx.Model(o).Updates(account).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Delete : 删除数据
func (o *Account) Delete() error {
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

// SelectByGroupID : 查询所有用户，按用户组.
func (o *Account) SelectByGroupID(groupID string, limit, offset int) (os []Account, count int) {
	db := DB
	os = make([]Account, 0)
	sql := `from t_account
			inner join t_account_to_group on t_account_to_group.account_id = t_account.id
			where t_account_to_group.group_id = ? `

	if o.Name != "" {
		sql = sql + "and name LIKE '%" + o.Name + "%' "
	}
	if o.Mail != "" {
		sql = sql + " and mail LIKE '%" + o.Mail + "%' "
	}
	if o.Phone != "" {
		sql = sql + " and phone LIKE '%" + o.Phone + "%' "
	}

	db.Raw("select count(t_account.*) "+sql, groupID).Count(&count)
	db.Raw("select t_account.* "+sql, groupID).Limit(limit).Offset(offset).Scan(&os)
	return
}
