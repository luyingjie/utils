package psql

import "time"

type Account struct {
	ID         string      `gorm:"column:id;type:varchar;PRIMARY_KEY "form:"id" json:"id"`
	Name       string      `gorm:"column:name;type:varchar;"json:"name"`
	Password   string      `gorm:"column:password;type:varchar;"json:"password"`
	UseStatus  string      `gorm:"column:use_status;type:varchar; "form:"use_status" json:"use_status"` //使用状态0：正常，-1删除
	Type       string      `gorm:"column:type;type:varchar"json:"type"`
	Mail       string      `gorm:"column:mail;type:varchar"json:"mail"`
	Phone      string      `gorm:"column:phone;type:varchar"json:"phone"`
	OrgID      interface{} `gorm:"column:org_id;type:varchar"json:"org_id"`
	GroupID    interface{} `gorm:"column:group_id;type:varchar"json:"group_id"`
	CreateTime time.Time   `gorm:"column:create_time;type:timestamp; "form:"create_time" json:"create_time"`
}

//TableName 设置表名
func (Account) TableName() string {
	return "t_account"
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
func (o *Account) SelectList(limit, offset int) (os []Account) {
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
	db.Raw(sql).Limit(limit).Offset(offset).Scan(&os)
	// db.Where(o).Limit(limit).Offset(offset).Scan(&os)
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
