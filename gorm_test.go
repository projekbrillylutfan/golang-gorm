package belajar_golang_gorm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:admin@tcp(localhost:3306)/belajar_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return db
}

var db = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, db)
}

func TestExecuteSQL(t *testing.T) {
	err := db.Exec("INSERT INTO sample(id, name) VALUES(?, ?)", "1", "John").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample(id, name) VALUES(?, ?)", "2", "ayam").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample(id, name) VALUES(?, ?)", "3", "hitam").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample(id, name) VALUES(?, ?)", "4", "legam").Error
	assert.Nil(t, err)
}

type Sample struct {
	Id   string
	Name string
}

func TestRawSQL(t *testing.T) {
	var sample Sample
	err := db.Raw("select * from sample where id = ?", "1").Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "John", sample.Name)

	var samples []Sample
	err = db.Raw("select id, name from sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(samples))
}

func TestSqlRow(t *testing.T) {
	var samples []Sample

	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	for rows.Next() {
		var id string
		var name string

		err := rows.Scan(&id, &name)
		assert.Nil(t, err)

		samples = append(samples, Sample{Id: id, Name: name})
	}

	assert.Equal(t, 4, len(samples))
}

func TestScanRow(t *testing.T) {
	var samples []Sample

	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	for rows.Next() {
		err := db.ScanRows(rows, &samples)
		assert.Nil(t, err)
	}

	assert.Equal(t, 4, len(samples))
}

func TestCreateUser(t *testing.T) {
	user := User{
		ID:       "1",
		Password: "password",
		Name: Name{
			FirstName:  "John",
			MiddleName: "Doe",
			LastName:   "ireng",
		},
		Information: "information",
	}

	response := db.Create(&user)
	assert.Nil(t, response.Error)
	assert.Equal(t, 1, int(response.RowsAffected))
}

func TestBatchInsert(t *testing.T) {
	var users []User
	for i := 2; i < 10; i++ {
		users = append(users, User{
			ID:       strconv.Itoa(i),
			Password: "password",
			Name: Name{
				FirstName:  "John" + strconv.Itoa(i),
				MiddleName: "Doe" + strconv.Itoa(i),
				LastName:   "ireng" + strconv.Itoa(i),
			},
		})
	}

	result := db.Create(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 8, int(result.RowsAffected))
}

func TestTransactionSuccses(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			ID:       "10",
			Password: "password",
			Name: Name{
				FirstName:  "John",
				MiddleName: "Doe",
				LastName:   "ireng",
			},
		}).Error
		if err != nil {
			return err
		}
		err = tx.Create(&User{
			ID:       "11",
			Password: "password",
			Name: Name{
				FirstName:  "John",
				MiddleName: "Doe",
				LastName:   "ireng",
			},
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       "12",
			Password: "password",
			Name: Name{
				FirstName:  "John",
				MiddleName: "Doe",
				LastName:   "ireng",
			},
		}).Error
		if err != nil {
			return err
		}
		return nil
	})

	assert.Nil(t, err)
}

func TestTransactionError(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			ID:       "13",
			Password: "password",
			Name: Name{
				FirstName:  "John",
				MiddleName: "Doe",
				LastName:   "ireng",
			},
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       "11",
			Password: "password",
			Name: Name{
				FirstName:  "John",
				MiddleName: "Doe",
				LastName:   "ireng",
			},
		}).Error
		if err != nil {
			return err
		}
		return nil
	})

	assert.NotNil(t, err)
}

func TestManualTransactionSuccses(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{
		ID:       "13",
		Password: "password",
		Name: Name{
			FirstName:  "John",
			MiddleName: "Doe",
			LastName:   "ireng",
		},
	}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{
		ID:       "14",
		Password: "password",
		Name: Name{
			FirstName:  "John",
			MiddleName: "Doe",
			LastName:   "ireng",
		},
	}).Error

	assert.Nil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestManualTransactionError(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{
		ID:       "15",
		Password: "password",
		Name: Name{
			FirstName:  "John",
			MiddleName: "Doe",
			LastName:   "ireng",
		},
	}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{
		ID:       "14",
		Password: "password",
		Name: Name{
			FirstName:  "John",
			MiddleName: "Doe",
			LastName:   "ireng",
		},
	}).Error

	assert.NotNil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestQuerySingleObject(t *testing.T) {
	user := User{}
	result := db.First(&user)
	assert.Nil(t, result.Error)
	assert.Equal(t, "1", user.ID)

	user = User{}
	result = db.Last(&user)
	assert.Nil(t, result.Error)
	assert.Equal(t, "9", user.ID)
}

func TestQuerySingleObjectInlineCondition(t *testing.T) {
	user := User{}
	result := db.First(&user, "id = ?", "5")
	assert.Nil(t, result.Error)
	assert.Equal(t, "5", user.ID)
}

func TestQueryAllObject(t *testing.T) {
	var users []User
	result := db.Find(&users, "id in ?", []string{"1", "2", "3", "4"})
	assert.Nil(t, result.Error)
	assert.Equal(t, 4, len(users))
}

func TestQueryCondition(t *testing.T) {
	var users []User
	result := db.Where("first_name like ?", "%Jo%").
		Where("passsword = ?", "password").Find(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 13, len(users))
}

func TestOrOperator(t *testing.T) {
	var users []User
	result := db.Where("first_name like ?", "%Jo%").Or("passsword = ?", "password").Find(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 13, len(users))
}

func TestNotOperator(t *testing.T) {
	var users []User
	result := db.Not("first_name like ?", "%Jo%").Where("passsword = ?", "password").Find(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 1, len(users))
}

func TestSelectFields(t *testing.T) {
	var users []User
	result := db.Select("id, name").Find(&users)
	assert.Nil(t, result.Error)

	for _, user := range users {
		assert.Equal(t, "id", user.ID)
		assert.Equal(t, "name", user.Name)
	}
	assert.Equal(t, 14, len(users))
}

func TestStructCondition(t *testing.T) {
	userCondition := User{
		Name: Name{
			FirstName: "John",
		},
	}

	var users []User
	result := db.Where(userCondition).Find(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 1, len(users))
}

func TestMapCondition(t *testing.T) {
	mapCondition := map[string]interface{}{
		"middle_name": "",
	}

	var users []User
	result := db.Where(mapCondition).Find(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 14, len(users))
}

func TestOrderLimitOffset(t *testing.T) {
	var users []User
	result := db.Order("id asc, first_name asc").Limit(5).Offset(5).Find(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 5, len(users))
	assert.Equal(t, "14", users[0].ID)
}

type UserResponse struct {
	ID        string
	FirstName string
	LastName  string
}

func TestQueryNonModel(t *testing.T) {
	var users []UserResponse

	result := db.Model(&User{}).Select("id, first_name, last_name").Find(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, 14, len(users))
	fmt.Println(users)
}

func TestUpdate(t *testing.T) {
	user := User{}

	result := db.First(&user, "id = ?", "1")
	assert.Nil(t, result.Error)

	user.Name.FirstName = "Joko"
	user.Name.MiddleName = "Ayam"
	user.Name.LastName = "Goreng"
	user.Password = "123"
	result = db.Save(&user)
	assert.Nil(t, result.Error)
}

func TestSelectedColumns(t *testing.T) {
	result := db.Model(&User{}).Where("id = ?", "1").Updates(map[string]interface{}{
		"middle_name": "",
		"last_name":   "Morro",
	})
	assert.Nil(t, result.Error)

	result = db.Model(&User{}).Where("id = ?", "1").Update("password", "000")
	assert.Nil(t, result.Error)

	result = db.Where("id = ?", "1").Updates(User{
		Name: Name{
			FirstName: "Eko",
			LastName:  "kanedi",
		},
	})

	assert.Nil(t, result.Error)
}

func TestAutoIncrement(t *testing.T) {
	for i := 0; i < 10; i++ {
		userLog := UserLog{
			UserID: "1",
			Action: "Test Action",
		}
		result := db.Create(&userLog)
		assert.Nil(t, result.Error)

		assert.NotEqual(t, 0, userLog.ID)
		fmt.Println(userLog.ID)
	}
}

func TestSaveOrUpdate(t *testing.T) {
	userLog := UserLog{
		UserID: "1",
		Action: "Test Action",
	}

	result:= db.Save(&userLog) // create
	assert.Nil(t, result.Error)

	userLog.UserID = "2"
	result = db.Save(&userLog) // update
	assert.Nil(t, result.Error)
}

func TestSaveOrUpdateNonAutoIncrement(t *testing.T) {
	user:= User{
		ID: "99",
		Name: Name{
			FirstName: "User 99",
		},
	}
	result := db.Save(&user) // create
	assert.Nil(t, result.Error)

	user.Name.FirstName = "User 99 Updated"
	result = db.Save(&user) // update
	assert.Nil(t, result.Error)
}
