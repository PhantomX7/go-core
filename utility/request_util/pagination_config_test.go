package request_util_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"github.com/PhantomX7/go-core/lib/scope"
	"github.com/PhantomX7/go-core/utility/request_util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

type TestPaginationConfigSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *gorm.DB
}

func TestPagination(t *testing.T) {
	suite.Run(t, new(TestPaginationConfigSuite))
}

func (suite *TestPaginationConfigSuite) SetupSuite() {
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

func (suite *TestPaginationConfigSuite) TestNewPaginationConfig() {
	suite.Run("with limit", func() {
		pagination := request_util.NewPaginationConfig(100, 0, "")

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(0, pagination.Offset())          // skip scope
		suite.Equal("", pagination.Order())          // skip scope
		suite.Equal(1, len(pagination.MetaScopes())) // total 1 meta scope
		suite.Nil(pagination.QueryMap())
	})

	suite.Run("with limit and offset", func() {
		pagination := request_util.NewPaginationConfig(100, 1, "")

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(1, pagination.Offset())          // add to scope
		suite.Equal("", pagination.Order())          // skip scope
		suite.Equal(2, len(pagination.MetaScopes())) // total 2 meta scope
		suite.Nil(pagination.QueryMap())
	})

	suite.Run("with limit, offset and order", func() {
		pagination := request_util.NewPaginationConfig(100, 1, "name asc")

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(1, pagination.Offset())          // add to scope
		suite.Equal("name asc", pagination.Order())  // add to scope
		suite.Equal(3, len(pagination.MetaScopes())) // total 3 meta scope
		suite.Nil(pagination.QueryMap())
	})

	suite.Run("with limit, offset, order and custom scope", func() {
		pagination := request_util.NewPaginationConfig(100, 1, "name asc",
			func(db *gorm.DB) *gorm.DB {
				return db.Where("test")
			},
		)

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(1, pagination.Offset())          // add to scope
		suite.Equal("name asc", pagination.Order())  // add to scope
		suite.Equal(3, len(pagination.MetaScopes())) // total 3 meta scope
		suite.Equal(1, len(pagination.Scopes()))     //  1 scope from args
		suite.Nil(pagination.QueryMap())
	})

	suite.Run("with limit, offset, order and pre made scope", func() {
		pagination := request_util.NewPaginationConfig(100, 1, "name asc",
			scope.WhereInScope("test", []interface{}{1, true, "4"}),
		)

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(1, pagination.Offset())          // add to scope
		suite.Equal("name asc", pagination.Order())  // add to scope
		suite.Equal(3, len(pagination.MetaScopes())) // total 3 meta scope
		suite.Equal(1, len(pagination.Scopes()))     // 1 extra scope from args
		suite.Nil(pagination.QueryMap())
	})

	suite.Run("with limit, offset, order and scope with add scope function", func() {
		pagination := request_util.NewPaginationConfig(100, 1, "name asc",
			scope.WhereInScope("test", []interface{}{1, true, "4"}),
		)

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(1, pagination.Offset())          // add to scope
		suite.Equal("name asc", pagination.Order())  // add to scope
		suite.Equal(3, len(pagination.MetaScopes())) // total 3 scope
		suite.Equal(1, len(pagination.Scopes()))     // 1 extra scope from args
		suite.Nil(pagination.QueryMap())

		pagination.AddScope(scope.OffsetScope(1))
		suite.Equal(2, len(pagination.Scopes()))
	})

	suite.Run("default pagination", func() {
		pagination := request_util.NewDefaultPaginationConfig()

		suite.Equal(20, pagination.Limit())          // add to scope
		suite.Equal(0, pagination.Offset())          // add to scope
		suite.Equal("", pagination.Order())          // add to scope
		suite.Equal(1, len(pagination.MetaScopes())) // total 1 meta scope
		suite.Equal(0, len(pagination.Scopes()))     // total 1 scope
		suite.Nil(pagination.QueryMap())
	})
}

func (suite *TestPaginationConfigSuite) TestNewRequestPaginationConfig() {
	suite.Run("with limit", func() {
		pagination := request_util.NewRequestPaginationConfig(
			map[string][]string{
				"limit": {"100"},
			},
			map[string]string{},
		)

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(0, pagination.Offset())          // skip scope
		suite.Equal("id desc", pagination.Order())   // add to scope because default order
		suite.Equal(2, len(pagination.MetaScopes())) // total 2 meta scope
		suite.NotNil(pagination.QueryMap())
	})

	suite.Run("with limit exceed 100", func() {
		pagination := request_util.NewRequestPaginationConfig(
			map[string][]string{
				"limit": {"200"},
			},
			map[string]string{},
		)

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(0, pagination.Offset())          // skip scope
		suite.Equal("id desc", pagination.Order())   // add to scope because default order
		suite.Equal(2, len(pagination.MetaScopes())) // total 2 meta scope
		suite.NotNil(pagination.QueryMap())
	})

	suite.Run("with limit and offset", func() {
		pagination := request_util.NewRequestPaginationConfig(
			map[string][]string{
				"limit":  {"100"},
				"offset": {"1"},
			},
			map[string]string{},
		)

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(1, pagination.Offset())          // add to scope
		suite.Equal("id desc", pagination.Order())   // add to scope because default order
		suite.Equal(3, len(pagination.MetaScopes())) // total 3 meta scope
		suite.NotNil(pagination.QueryMap())
	})

	suite.Run("with limit, offset and order", func() {
		pagination := request_util.NewRequestPaginationConfig(
			map[string][]string{
				"limit":  {"100"},
				"offset": {"1"},
				"sort":   {"name asc"},
			},
			map[string]string{},
		)

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(1, pagination.Offset())          // add to scope
		suite.Equal("name asc", pagination.Order())  // add to scope
		suite.Equal(3, len(pagination.MetaScopes())) // total 3 meta scope
		suite.NotNil(pagination.QueryMap())
	})

	suite.Run("with limit, offset, order and query", func() {
		pagination := request_util.NewRequestPaginationConfig(
			map[string][]string{
				"limit":  {"100"},
				"offset": {"1"},
				"sort":   {"name asc"},
				"name":   {"test"},
			},
			map[string]string{
				"name": request_util.StringType,
			},
		)

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(1, pagination.Offset())          // add to scope
		suite.Equal("name asc", pagination.Order())  // add to scope
		suite.Equal(3, len(pagination.MetaScopes())) // total 3 meta scope
		suite.Equal(1, len(pagination.Scopes()))     // 1 extra scope from query
		suite.NotNil(pagination.QueryMap())
	})

	suite.Run("with limit, offset, order and query with add scope function", func() {
		pagination := request_util.NewRequestPaginationConfig(
			map[string][]string{
				"limit":  {"100"},
				"offset": {"1"},
				"sort":   {"name asc"},
				"name":   {"test"},
			},
			map[string]string{
				"name": request_util.StringType,
			},
		)

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(1, pagination.Offset())          // add to scope
		suite.Equal("name asc", pagination.Order())  // add to scope
		suite.Equal(3, len(pagination.MetaScopes())) // total 3 scope
		suite.Equal(1, len(pagination.Scopes()))     // 1 extra scope from query
		suite.NotNil(pagination.QueryMap())

		pagination.AddScope(scope.OffsetScope(1))
		suite.Equal(2, len(pagination.Scopes()))
	})

	suite.Run("with limit, offset, order and lot of query", func() {
		conditions := map[string][]string{
			"limit":     {"100"},
			"offset":    {"1"},
			"sort":      {"name asc"},
			"name":      {"test"},
			"id":        {"1"},
			"price":     {"1000,2000"},
			"is_active": {"true"},
			"date":      {"2017-10-13,2017-10-13"},
			"datetime":  {"2017-10-13,2017-10-13"},
			"excluded":  {"test"},
		}
		pagination := request_util.NewRequestPaginationConfig(
			conditions,
			map[string]string{
				"name":      request_util.StringType,
				"id":        request_util.IdType,
				"price":     request_util.NumberType,
				"is_active": request_util.BoolType,
				"date":      request_util.DateType,
				"datetime":  request_util.DatetimeType,
			},
		)

		request_util.OverrideKey(conditions, "excluded", "included")

		suite.Equal(100, pagination.Limit())         // add to scope
		suite.Equal(1, pagination.Offset())          // add to scope
		suite.Equal("name asc", pagination.Order())  // add to scope
		suite.Equal(3, len(pagination.MetaScopes())) // total 3 meta scope
		suite.Equal(6, len(pagination.Scopes()))     // 6 extra scope from query
		suite.NotNil(pagination.QueryMap())
	})
}
