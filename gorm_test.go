package belajar_golang_gorm

//belajar lagi golang gorm

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

//open connection di gorm
func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:admin@tcp(localhost:3306)/belajar_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{
		//add logger
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

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
	fmt.Println(samples)
}

func TestSqlRow(t *testing.T) {
	var samples []*Sample

	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	for rows.Next() {
		var id string
		var name string

		err := rows.Scan(&id, &name)
		assert.Nil(t, err)

		samples = append(samples, &Sample{Id: id, Name: name})
	}

	assert.Equal(t, 4, len(samples))
	fmt.Println(samples)
}

func TestScanRow(t *testing.T) {
	var samples []*Sample

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
		ID:       "80",
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
	fmt.Println(response)
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
	fmt.Println(user)
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
	result := db.Order("id asc, first_name asc").Limit(5).Offset(5).Find(&users).Error
	assert.Nil(t, result)
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


// update di gorm
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


// ini patch
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

//save or update
func TestSaveOrUpdate(t *testing.T) {
	userLog := UserLog{
		UserID: "1",
		Action: "Test Action",
	}

	result := db.Save(&userLog) // create
	assert.Nil(t, result.Error)

	userLog.UserID = "2"
	result = db.Save(&userLog) // update
	assert.Nil(t, result.Error)
}

func TestSaveOrUpdateNonAutoIncrement(t *testing.T) {
	user := User{
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

func TestConflict(t *testing.T) {
	user := User{
		ID: "88",
		Name: Name{
			FirstName: "User 88",
		},
	}
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user) // create
	assert.Nil(t, result.Error)
}

func TestDelete(t *testing.T) {
	var user User
	result := db.First(&user, "id = ?", "88")
	assert.Nil(t, result.Error)
	result = db.Delete(&user)
	assert.Nil(t, result.Error)

	result = db.Delete(&User{}, "id = ?", "99")
	assert.Nil(t, result.Error)

	result = db.Where("id = ?", "77").Delete(&User{})
	assert.Nil(t, result.Error)
}

func TestSoftDelete(t *testing.T) {
	todo := Todo{
		UserID: "1",
		Title:  "Todo 1",
		Description: "Isi todo 1",
	}
	result:= db.Create(&todo)
	assert.Nil(t, result.Error)

	result = db.Delete(&todo)
	assert.Nil(t, result.Error)
	assert.NotNil(t, todo.DeletedAt)

	var todos []Todo
	result = db.Find(&todos)
	assert.Nil(t, result.Error)
	assert.Equal(t, 0, len(todos))
}

func TestUnscoped(t *testing.T) {
	var todo Todo
	result := db.Unscoped().First(&todo, "id = ?", "2")
	assert.Nil(t, result.Error)

	result = db.Unscoped().Delete(&todo)
	assert.Nil(t, result.Error)

	var todos []Todo
	result = db.Unscoped().Find(&todos)
	assert.Nil(t, result.Error)
	assert.Equal(t, 0, len(todos))
}

func TestLock(t *testing.T) {
	err:= db.Transaction(func (tx *gorm.DB) error {
		var user User
		err := tx.Clauses(clause.Locking{Strength: "UPDATE",}).First(&user, "id = ?", "1").Error
		if err != nil {
			return err
		}
		user.Name.FirstName = "Joko"
		user.Name.LastName = "Morro"
		return tx.Save(&user).Error
	})

	assert.Nil(t, err)
}

func TestCreateWallet(t *testing.T) {
	wallet := Wallet{
		ID: "1",
		UserID: "1",
		Balance: 1000000,
	}

	err:= db.Create(&wallet).Error
	assert.Nil(t, err)
}

func TestRetrieveRelation(t *testing.T) {
	var user User
	err:= db.Model(&user).Preload("Wallet").First(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "1", user.Wallet.ID)
}

func TestRetrieveRelationJoin(t *testing.T) {
	var users []User
	err := db.Model(&User{}).Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)

	assert.Equal(t, 14, len(users))
	fmt.Println(users)
}

func TestAutoCreateUpdate(t *testing.T) {
	user:= User{
		ID: "20",
		Password: "password",
		Name: Name{
			FirstName:  "User 20",
		},
		Wallet: Wallet{
			ID: "20",
			UserID: "20",
			Balance: 1000000,
		},
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestSkipAutoCreateUpdate(t *testing.T) {
	user:= User{
		ID: "21",
		Password: "password",
		Name: Name{
			FirstName:  "User 21",
		},
		Wallet: Wallet{
			ID: "21",
			UserID: "21",
			Balance: 1000000,
		},
	}
	err := db.Omit(clause.Associations).Create(&user).Error
	assert.Nil(t, err)
}

func TestUserAndAddresses(t *testing.T) {
	user := User {
		ID: "2",
		Password: "password",
		Name: Name{
			FirstName:  "User 2",
		},
		Wallet: Wallet{
			ID: "2",
			UserID: "2",
			Balance: 1000000,
		},
		Addresses: []Address{
			{
				UserID: "2",
				Address: "Jl. Raya No. 2",
			},
			{
				UserID: "2",
				Address: "Jl. Raya No. 51",
			},
		},
	}

	err := db.Save(&user).Error
	assert.Nil(t, err)
}

func TestPreloadJoinOneToMany(t *testing.T) {
	var userPreload []User
	err:= db.Model(&User{}).Preload("Addresses").Joins("Wallet").Find(&userPreload).Error
	assert.Nil(t, err)
	fmt.Println(userPreload)
}

func TestTakePreloadJoinOneToMany(t *testing.T) {
	var user User
	err:= db.Model(&User{}).Preload("Addresses").Joins("Wallet").Take(&user, "users.id = ?", "50").Error
	assert.Nil(t, err)
	fmt.Println(user)
}

func TestBelongsTo(t *testing.T) {
	fmt.Println("preload")
	var addresses []Address
	err:= db.Preload("User").Find(&addresses).Error
	assert.Nil(t, err)

	fmt.Println("joins")
	addresses = []Address{}
	err = db.Joins("User").Find(&addresses).Error
	assert.Nil(t, err)
}

func TestBelongsToOneToOne(t *testing.T) {
	fmt.Println("preload")
	var wallets []Wallet
	err:= db.Model(&Wallet{}).Preload("User").Find(&wallets).Error
	assert.Nil(t, err)

	fmt.Println("joins")
	wallets = []Wallet{}
	err = db.Model(&Wallet{}).Joins("User").Find(&wallets).Error
	assert.Nil(t, err)
}

func TestCreateManyToMany(t *testing.T) {
	product := Product{
		ID: "P002",
		Name: "Product 2",
		Price: 100000,
	}
	err:= db.Create(&product).Error
	assert.Nil(t, err)

	err = db.Table("user_like_products").Create(map[string]interface{}{
		"user_id": "1",
		"product_id": "P002",
	}).Error
	assert.Nil(t, err)

	err = db.Table("user_like_products").Create(map[string]interface{}{
		"user_id": "2",
		"product_id": "P002",
	}).Error
	assert.Nil(t, err)
}

func TestPreloadManyToMany(t *testing.T) {
	var product Product
	err:= db.Preload("LikedByUsers").First(&product, "id = ?", "P002").Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(product.LikedByUsers))
	fmt.Println(product)
}

func TestPreloadManyToManyUser(t *testing.T) {
	var user User
	err:= db.Preload("LikedProducts").Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(user.LikedProducts))
}

func TestAssociationFind(t *testing.T) {
	var product Product
	err:= db.First(&product, "id = ?", "P002").Error
	assert.Nil(t, err)

	var users []User
	err = db.Model(&product).Where("users.first_name LIKE ?", "J%").Association("LikedByUsers").Find(&users)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestAssociationAdd(t *testing.T) {
	var user User
	err:= db.First(&user, "id = ?", "3").Error
	assert.Nil(t, err)

	var product Product
	err = db.First(&product, "id = ?", "P002").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Append(&user)
	assert.Nil(t, err)
}

func TestAssociationRepalace(t *testing.T) {
	err:= db.Transaction(func (tx *gorm.DB) error  {
		var user User
		err:= tx.First(&user, "id = ?", "1").Error
		assert.Nil(t, err)

		Wallet := Wallet{
			ID: "01",
			UserID: "1",
			Balance: 1000000,
		}
		err = tx.Model(&user).Association("Wallet").Replace(&Wallet)
		return err
	})
	assert.Nil(t, err)
}

func TestAssociationDelete(t *testing.T) {
	var user User
	err:= db.First(&user, "id = ?", "3").Error
	assert.Nil(t, err)

	var product Product
	err = db.First(&product, "id = ?", "P002").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Delete(&user)
	assert.Nil(t, err)
}

func TestAssociationClear(t *testing.T) {
	var product Product
	err := db.First(&product, "id = ?", "P002").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Clear()
	assert.Nil(t, err)
}

func TestPreloadingWithCondition(t *testing.T) {
	var user User
	err := db.Preload("Wallet", "balance > ?", 10000).First(&user, "id = ?", "1").Error
	assert.Nil(t, err)
	fmt.Println(user)
}

func TestNestedPreload(t *testing.T) {
	var wallet Wallet
	err:= db.Preload("User.Addresses").Find(&wallet, "id = ?", "2").Error
	assert.Nil(t, err)
	fmt.Println(wallet)
	fmt.Println(wallet.User)
	fmt.Println(wallet.User.Addresses)
}

func TestPreloadAll(t *testing.T) {
	var user User
	err := db.Preload(clause.Associations).First(&user, "id = ?", "1").Error
	assert.Nil(t, err)
	fmt.Println(user)
}

func TestJoinQuery(t *testing.T) {
	var users []User
	err := db.Joins("join wallets on wallets.user_id = users.id").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))

	users = []User{}
	err = db.Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 17, len(users))
}

func TestJoinQueryCondition(t *testing.T) {
	var users []User
	err := db.Joins("join wallets on wallets.user_id = users.id AND wallets.balance > ?", 10000).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))

	users = []User{}
	err = db.Joins("Wallet").Where("balance > ?", 10000).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))
}

func TestCount(t *testing.T) {
	var count int64
	err := db.Model(&User{}).Joins("Wallet").Where("Wallet.balance > ?", 10000).Count(&count).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(4), count)
}

type AggregationResult struct {
	TotalBalance int64
	MinBalance int64
	MaxBalance int64
	AvgBalance int64
}

func TestAggregation(t *testing.T) {
	var result AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance" , "min(balance) as min_balance", "max(balance) as max_balance", "avg(balance) as avg_balance").Scan(&result).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(4000000), result.TotalBalance)
	assert.Equal(t, int64(1000000), result.MinBalance)
	assert.Equal(t, int64(1000000), result.MaxBalance)
	assert.Equal(t, int64(500000), result.AvgBalance)
}

func TestGroupByHaving(t *testing.T) {
	var result []AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance" , "min(balance) as min_balance", "max(balance) as max_balance", "avg(balance) as avg_balance").Joins("User").Group("User.id").Having("sum(balance) > ?", 2000000).Find(&result).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestContext(t *testing.T) {
	ctx := context.Background()

	var users []User
	err := db.WithContext(ctx).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 17, len(users))
}

func BrokenWalletBalance(db *gorm.DB) *gorm.DB{
	return db.Where("balance = ?", 0)
}

func SulatanWalletBalance(db *gorm.DB) *gorm.DB{
	return db.Where("balance > ?", 100000)
}

func TestScopes(t *testing.T) {
	var wallets []Wallet
	err:= db.Scopes(BrokenWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)

	wallets = []Wallet{}
	err = db.Scopes(SulatanWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)
}

func TestMigrator(t *testing.T) {
	err := db.Migrator().AutoMigrate(&GuestBook{})
	assert.Nil(t, err)
}

func TestHook(t *testing.T) {
	user := User{
		Password: "rahasia",
		Name: Name{
			FirstName: "User 100",
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
	assert.NotEqual(t, "", user.ID)

	fmt.Println(user.ID)
}
