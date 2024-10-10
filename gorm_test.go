package belajar_golang_gorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:admin@tcp(localhost:3306)/belajar_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{})
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
