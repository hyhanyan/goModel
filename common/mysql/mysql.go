package mysql

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

type mysqlDb struct {
	User     string
	Password string
	Host     string
	Database string
	Port     int
	DB       *gorm.DB
}

type NullType byte

var MySqlDb *mysqlDb

func NewMySqlDb(user, passwd, host, db string, port int) *mysqlDb {
	return &mysqlDb{User: user, Password: passwd, Host: host, Database: db, Port: port}
}

func (Db *mysqlDb) DbConn() error {
	connArgs := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", Db.User, Db.Password, Db.Host, Db.Port, Db.Database)
	db, err := gorm.Open("mysql", connArgs)
	if nil != err {
		panic(err)
	}
	// 全局禁用表名复数
	db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	Db.DB = db
	return nil
}

func (Db *mysqlDb) Close() {
	Db.DB.Close()
}

func (Db *mysqlDb) Insert(tablename string, value interface{}) {
	Db.DB.Table(tablename).Create(value)
}

// value interface{} 为struct 或者 map
func (Db *mysqlDb) UpdateByConds(tablename string, updateVals interface{}, where interface{}) {
	Db.DB.Table(tablename).Where(where).Updates(updateVals)
}

func (Db *mysqlDb) SelectByConds(tablename string, findVals interface{}, where interface{}) {
	Db.DB.Table(tablename).Where(where).Find(findVals)
}

func (Db *mysqlDb) DeleteByConds(tablename string, deletevals interface{}, where interface{}) {
	Db.DB.Table(tablename).Where(where).Delete(deletevals)
}

func (Db *mysqlDb) SelectBySqlString(tablename string, findVals interface{}, where map[string]interface{}) {
	cond, vals, err := whereBuild(where)
	if nil != err {
		panic(err)
	}
	Db.DB.Table(tablename).Where(cond, vals...).Find(findVals)
}

func (Db *mysqlDb) UpdateBySqlString(tablename string, updateVals interface{}, where map[string]interface{}) {
	cond, vals, err := whereBuild(where)
	if nil != err {
		panic(err)
	}
	Db.DB.Table(tablename).Where(cond, vals...).Updates(updateVals)
}

func (Db *mysqlDb) DeleteBySqlString(tablename string, deletevalsVals interface{}, where map[string]interface{}) {
	cond, vals, err := whereBuild(where)
	if nil != err {
		panic(err)
	}
	Db.DB.Table(tablename).Where(cond, vals...).Delete(deletevalsVals)
}

// sql build where
func whereBuild(where map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}

		if whereSQL != "" {
			whereSQL += " AND "
		}
		switch len(ks) {
		case 1:
			//fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, "<>?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				vals = append(vals, v)
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)
			}
			break
		}
	}
	return
}
