package storage

import (
	"apistream/utils"
	"time"
	"strings"
	"github.com/go-sql-driver/mysql"
	"reflect"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"errors"
)

type ChannelTable struct {
	Id 					int64			`json:"id"`
	ChannelName 		string			`json:"channel_name"`
	ChannelAliasName 	string			`json:"channel_alias_name,omitempty"`
	OwnerId				int64			`json:"owner_id"`
	TimeStart			int64			`json:"time_start,omitempty"`
	TimeEnd				int64			`json:"time_end,omitempty"`
	Storage				string			`json:"storage"`
	Password			string			`json:"password,omitempty"`
}

type Count struct {
	Count int64	`json:"count"`
}

type MySqlClient struct {
	Client *sql.DB
}

type Table struct {
	Name                      string
	IgnoreColumns             []string
	DateTimeColumns           []string
	AutoUpdateDateTimeColumns []string
	ForeignKeys               []string
	NotNullColumns            []string
	ConvertJSONKey2Column     map[string]string
}

type SQLTool struct {
	db               *sql.DB
	table            Table
	columns          []string
	column2FieldName map[string]string
	column2Kind      map[string]reflect.Kind
	values           []interface{}
}

func NewMySqlClient(username string, password string, db string, hostname string, port string) MySqlClient {	
	client, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?autocommit=true", username, password, hostname, port,db))
	if err != nil {
		panic(err.Error())
	}
	// defer client.Close()
	return MySqlClient{
		Client: client,
	}
}

// PrepareSelect -- prepare select query
func PrepareSelect(db *sql.DB, table Table, protoObj interface{}) (st SQLTool) {
	st.db = db
	st.table = table
	st.parseColumns(protoObj)
	return
}


// PrepareInsert -- prepare insert query
func PrepareInsert(db *sql.DB, table Table, protoObj interface{}) (st SQLTool) {
	st.db = db
	st.table = table
	st.parseColumns(protoObj)
	st.fillValues(protoObj)
	return
}

// PrepareInsert -- prepare insert query
func PrepareCount(db *sql.DB, table Table) (st SQLTool) {
	st.db = db
	st.table = table
	return
}

func (st *SQLTool) selectQuery(query string, args ...interface{}) (rows *sql.Rows, err error) {
	stmt, err := st.db.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(args...)
	if err != nil {
		return
	}

	return
}

func (st *SQLTool) ExecuteInsert() sql.Result{
	// Prepare statement for inserting data
	stt := fmt.Sprintf("INSERT INTO %v (%v) VALUES(%v);", st.table.Name, st.GetQueryColumnList(), st.GetQueryValueList())
	stmtIns, err := st.db.Prepare(stt)
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()
	
	res, er := stmtIns.Exec(st.values...)
	if er != nil {
		panic(er.Error())
	}
	return res
}

func (st *SQLTool) ExecuteCount(query string, args ...interface{}) (ret int64, err error){
	stmt, err := st.db.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(args...).Scan(&ret)
	if err != nil {
		return 
	}	
	return ret, nil
}

func (st *SQLTool) ExecuteSelect(dest interface{}, query string, args ...interface{}) (err error) {
	rows, err := st.selectQuery(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	var vp reflect.Value

	value := reflect.ValueOf(dest)

	// json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	direct := reflect.Indirect(value)

	slice, err := baseType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}

	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := Deref(slice.Elem())

	for rows.Next() {
		vp = reflect.New(base)
		err = st.Scan(rows, vp.Interface())
		if err != nil {
			fmt.Println("Scan Errr:", err)
			continue
		}
		// append
		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
		}
	}

	return
}

// Scan -- scan row to object
func (st *SQLTool) Scan(rows *sql.Rows, dest interface{}) (err error) {
	st.values = make([]interface{}, 0)
	for _, column := range st.columns {
		kind, _ := st.column2Kind[column]

		switch kind {
		case reflect.String:
			var nstr = &sql.NullString{}
			st.values = append(st.values, nstr)
		case reflect.Bool:
			var nbool = &sql.NullBool{}
			st.values = append(st.values, nbool)
		case reflect.Float64:
			var nfloat64 = &sql.NullFloat64{}
			st.values = append(st.values, nfloat64)
		case reflect.Int64:
			if utils.IsStringSliceContains(st.table.DateTimeColumns, column) {
				var ntime = &mysql.NullTime{}
				st.values = append(st.values, ntime)
			} else {
				var nint64 = &sql.NullInt64{}
				st.values = append(st.values, nint64)
			}
		}
	}

	err = rows.Scan(st.values...)
	if err != nil {
		return
	}

	v := reflect.ValueOf(dest)
	// log.Print(v.Kind())
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}

	ve := v.Elem()
	// log.Print(ve.Kind())
	for index, column := range st.columns {
		// get field name
		fieldName, _ := st.column2FieldName[column]
		kind, _ := st.column2Kind[column]
		value := st.values[index]

		switch kind {
		case reflect.String:
			ve.FieldByName(fieldName).SetString(value.(*sql.NullString).String)
		case reflect.Bool:
			ve.FieldByName(fieldName).SetBool(value.(*sql.NullBool).Bool)
		case reflect.Float64:
			ve.FieldByName(fieldName).SetFloat(value.(*sql.NullFloat64).Float64)
		case reflect.Int64:
			if utils.IsStringSliceContains(st.table.DateTimeColumns, column) {
				var vtime = value.(*mysql.NullTime)
				if vtime.Valid && vtime.Time.Unix() >= 0 {
					ve.FieldByName(fieldName).SetInt(vtime.Time.Unix())
				}
			} else {
				ve.FieldByName(fieldName).SetInt(value.(*sql.NullInt64).Int64)
			}
		}
	}

	return
}
// GetQueryColumnList -- use for SELECT, INSERT query
func (st *SQLTool) GetQueryColumnList() string {
	return strings.Join(st.columns, ",")
}

// GetQueryValueList -- use for INSERT query
func (st *SQLTool) GetQueryValueList() string {
	questionMarks := make([]string, 0)
	for index := 0; index < len(st.columns); index++ {
		questionMarks = append(questionMarks, "?")
	}
	return strings.Join(questionMarks, ",")
}
func baseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = Deref(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}

// Deref is Indirect for reflect.Types
func Deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func (st *SQLTool) parseColumns(dest interface{}) {
	// Make a slice for the values
	st.values = make([]interface{}, 0)
	st.columns = make([]string, 0)
	st.column2Kind = make(map[string]reflect.Kind)
	st.column2FieldName = make(map[string]string, 0)
	
	t := reflect.TypeOf(dest).Elem()

	for index := 0; index < t.NumField(); index++ {
		f := t.Field(index)
		jsonKeys := f.Tag.Get("json")
		s := strings.Split(jsonKeys, ",")[0]
		st.column2FieldName[s]	= f.Name
		st.columns = append(st.columns, s)
		k := f.Type.Kind()
		st.column2Kind[s] = k
		switch k {
			case reflect.String:
				var nstr = &sql.NullString{}
				st.values = append(st.values, nstr)
				
			case reflect.Bool:
				var nbool = &sql.NullBool{}
				st.values = append(st.values, nbool)
			case reflect.Float64:
				var nfloat64 = &sql.NullFloat64{}
				st.values = append(st.values, nfloat64)
			case reflect.Int64:
				var in = false				
				for i := 0; i < len(st.table.DateTimeColumns); i++ {
					if st.table.DateTimeColumns[i] == f.Name{
						in = true
					}
				}
				if in {
					var nTime = &mysql.NullTime{}
					st.values = append(st.values, nTime)
				}else{
					var nint64 = &sql.NullInt64{}
					st.values = append(st.values, nint64)
				}				
		}
	}
}
// IsZeroOfUnderlyingType --
func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
func (st *SQLTool) fillValues(i interface{}) {
	st.values = make([]interface{}, 0)
	v := reflect.ValueOf(i).Elem()
	for _, column := range st.columns {
		// get field name
		fieldName, _ := st.column2FieldName[column]
		fieldValue := v.FieldByName(fieldName)
		fieldValueInterface := fieldValue.Interface()
		convertedValue := fieldValueInterface

		// if nullable column, check if zero
		if !utils.IsStringSliceContains(st.table.NotNullColumns, column) {
			if IsZeroOfUnderlyingType(fieldValueInterface) {
				convertedValue = nil
			}
		}

		// check if column is foreign key
		if utils.IsStringSliceContains(st.table.ForeignKeys, column) {
			if IsZeroOfUnderlyingType(fieldValueInterface) {
				convertedValue = nil
			}
		}

		// check if column is datetime
		if utils.IsStringSliceContains(st.table.DateTimeColumns, column) {
			if utils.IsStringSliceContains(st.table.AutoUpdateDateTimeColumns, column) {
				convertedValue = time.Now()
			} else {
				if IsZeroOfUnderlyingType(fieldValueInterface) {
					convertedValue = nil
				} else {
					convertedValue = time.Unix(fieldValue.Int(), 0)
				}
			}
		}

		st.values = append(st.values, convertedValue)
	}
}