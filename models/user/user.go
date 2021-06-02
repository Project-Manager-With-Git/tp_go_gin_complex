package user

import (
	"errors"
	"fmt"

	log "github.com/Golang-Tools/loggerhelper"
	"github.com/xormplus/xorm"
)

// User 用户类
// Actived字段为空表示未被激活
// OTPCode字段为空表示未激活二次验证功能
// 默认构造用户admin,它会加入admin群组
type User struct {
	ID   int32  `xorm:"pk autoincr 'id' comment('用户id')" json:"ID,omitempty" uri:"uid" binding:"required"`
	Name string `xorm:"varchar(25) notnull unique 'name' comment('用户名')" json:"Name,omitempty"`
}

func (u *User) TableName() string {
	return "user"
}

//GetByID 通过用户ID获得用户在数据库中的信息
func GetByID(db xorm.EngineInterface, id int32) (*User, error) {
	newu := User{}
	ok, err := db.Alias("T").Where("T.uid = ?", id).Get(&newu)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("用户id为%d的用户不存在", id)
	}
	return &newu, nil
}

//Count 查看数据库种的对象个数
func Count(db xorm.EngineInterface) (int64, error) {
	newu := User{}
	res, err := db.Alias("T").Count(newu)
	if err != nil {
		return 0, err
	}
	return res, nil
}

//Sync 根据用户id同步对应数据到当前对象
func (u *User) Sync(db xorm.EngineInterface) error {
	if u.ID > 0 {
		_, err := db.ID(u.ID).Get(u)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("User ID <= 0")
	}

}

//Delete 软删除用户,如果ID为>0的数则删除对应id的用户数据
func (u *User) Delete(db xorm.EngineInterface) error {
	if u.ID > 0 {
		user := User{}
		_, err := db.ID(u.ID).Delete(&user)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("User ID <= 0")
	}

}

//Save 保存对象到目标数据库,如果有ID则更新覆盖,没有则创建
func (u *User) Save(db xorm.EngineInterface) error {
	if u.ID > 0 {
		_, err := db.ID(u.ID).Update(u)
		return err
	} else {
		_, err := db.Insert(u)
		return err
	}
}

// 回调操作
// RegistCallback 需要注册到db代理上的回调
func RegistCallback(db xorm.EngineInterface) error {
	emptyrow := User{}
	ok, err := db.IsTableExist(&emptyrow)
	if err != nil {
		log.Error("user init error", log.Dict{
			"err":            err.Error(),
			"place":          "IsTableExist",
			"RegistCallback": "User",
		})
		return err
	}
	if !ok {
		err := db.CreateTables(&emptyrow)
		if err != nil {
			log.Error("user init error", log.Dict{
				"err":            err.Error(),
				"place":          "CreateTables",
				"RegistCallback": "User",
			})
			return err
		}
		err = db.CreateUniques(&emptyrow)
		if err != nil {
			log.Error("user init error",
				log.Dict{
					"err":            err.Error(),
					"place":          "CreateUniques",
					"RegistCallback": "User",
				})
			return err
		}
	}
	ok, err = db.IsTableEmpty(&emptyrow)
	if err != nil {
		log.Error("user init error", log.Dict{
			"err":            err.Error(),
			"place":          "IsTableEmpty",
			"RegistCallback": "User",
		})
		return err
	}
	if ok {
		adminUser := User{
			Name: "admin",
		}
		err := adminUser.Save(db)
		if err != nil {
			log.Error("user init error", log.Dict{
				"err":            err.Error(),
				"place":          "Insert admin",
				"RegistCallback": "User",
			})
			return err
		}
	}
	return nil
}
