package scope_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/PhantomX7/go-core/lib/scope"
)

type TestScopeSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *gorm.DB
}

func TestScope(t *testing.T) {
	suite.Run(t, new(TestScopeSuite))
}

func (suite *TestScopeSuite) SetupSuite() {
	db, mock := SetupDB()

	suite.mock = mock
	suite.db = db
}

func SetupDB() (*gorm.DB, sqlmock.Sqlmock) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		panic("setup mock database failed")
		return nil, nil
	}

	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(mock.
			NewRows([]string{"version()"}).
			AddRow("test_version"))
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: mockDb,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})

	return db, mock
}

func (suite *TestScopeSuite) TestLimitScope() {
	mock := suite.mock

	mock.ExpectQuery("SELECT `id` FROM `test` LIMIT 1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var result int
	err := suite.db.Table("test").
		Scopes(scope.LimitScope(1)).
		Pluck("id", &result).Error

	suite.Nil(err)
	suite.Equal(1, result)
}

func (suite *TestScopeSuite) TestOffsetScope() {
	mock := suite.mock

	mock.ExpectQuery("SELECT `id` FROM `test` OFFSET 1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var result int
	err := suite.db.Table("test").
		Scopes(scope.OffsetScope(1)).
		Pluck("id", &result).Error

	suite.Nil(err)
	suite.Equal(1, result)
}

func (suite *TestScopeSuite) TestOrderScope() {
	mock := suite.mock

	mock.ExpectQuery("SELECT `id` FROM `test` ORDER BY id").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var result int
	err := suite.db.Table("test").
		Scopes(scope.OrderScope("id")).
		Pluck("id", &result).Error

	suite.Nil(err)
	suite.Equal(1, result)
}

func (suite *TestScopeSuite) TestWhereNotInScope() {
	mock := suite.mock

	mock.ExpectQuery("SELECT `id` FROM `test` WHERE id NOT IN \\(\\?,\\?,\\?\\)").
		WithArgs(2, 3, 4).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var result int
	err := suite.db.Table("test").
		Scopes(scope.WhereNotInScope("id", []interface{}{2, 3, 4})).
		Pluck("id", &result).Error

	suite.Nil(err)
	suite.Equal(1, result)
}

func (suite *TestScopeSuite) TestWhereInScope() {
	mock := suite.mock

	mock.ExpectQuery("SELECT `id` FROM `test` WHERE id IN \\(\\?,\\?,\\?\\)").
		WithArgs(1, 2, 3).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var result int
	err := suite.db.Table("test").
		Scopes(scope.WhereInScope("id", []interface{}{1, 2, 3})).
		Pluck("id", &result).Error

	suite.Nil(err)
	suite.Equal(1, result)
}

func (suite *TestScopeSuite) TestWhereIsScope() {
	mock := suite.mock

	mock.ExpectQuery("SELECT `id` FROM `test` WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var result int
	err := suite.db.Table("test").
		Scopes(scope.WhereIsScope("id", 1)).
		Pluck("id", &result).Error

	suite.Nil(err)
	suite.Equal(1, result)
}

func (suite *TestScopeSuite) TestWhereIsNotScope() {
	mock := suite.mock

	mock.ExpectQuery("SELECT `id` FROM `test` WHERE id <> \\?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var result int
	err := suite.db.Table("test").
		Scopes(scope.WhereIsNotScope("id", 1)).
		Pluck("id", &result).Error

	suite.Nil(err)
	suite.Equal(1, result)
}

func (suite *TestScopeSuite) TestWhereLikeScope() {
	mock := suite.mock

	mock.ExpectQuery("SELECT `id` FROM `test` WHERE name LIKE \\?").
		WithArgs("%test%").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var result int
	err := suite.db.Table("test").
		Scopes(scope.WhereLikeScope("name", "test")).
		Pluck("id", &result).Error

	suite.Nil(err)
	suite.Equal(1, result)
}

func (suite *TestScopeSuite) TestWhereBetweenScope() {
	mock := suite.mock

	mock.ExpectQuery("SELECT `id` FROM `test` WHERE id BETWEEN \\? AND \\?").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var result int
	err := suite.db.Table("test").
		Scopes(scope.WhereBetweenScope("id", 1, 2)).
		Pluck("id", &result).Error

	suite.Nil(err)
	suite.Equal(1, result)
}

func (suite *TestScopeSuite) TestMultipleScope() {
	mock := suite.mock

	mock.ExpectQuery("SELECT `id` FROM `test` WHERE id = \\? AND name = \\? ORDER BY id asc LIMIT 1 OFFSET 1").
		WithArgs(1, "test").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var result int
	err := suite.db.Table("test").
		Scopes(
			scope.WhereIsScope("id", 1),
			scope.WhereIsScope("name", "test"),
			scope.LimitScope(1),
			scope.OrderScope("id asc"),
			scope.OffsetScope(1),
		).
		Pluck("id", &result).Error

	suite.Nil(err)
	suite.Equal(1, result)
}
