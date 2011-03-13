package gorm

import (
	"os"
	"sdegutis/sqlite"
	"fmt"
	"strings"
	// "log"
	"reflect"
)

type Conn struct {
	conn *sqlite.Conn
}

func (c *Conn) Close() os.Error {
	return c.conn.Close()
}

func OpenDB(filename string) (*Conn, os.Error) {
	conn, err := sqlite.Open(filename)
	return &Conn{conn: conn}, err
}

func getTableName(obj interface{}) string {
	return pluralizeString(snakeCasedName(getTypeName(obj)))
}

func (c *Conn) getResultsForQuery(tableName, condition string, args []interface{}) (resultsSlice []map[string][]byte, err os.Error) {
	s, err := c.conn.Prepare(fmt.Sprintf("select * from %v %v", tableName, condition))
	if err != nil {
		return nil, err
	}

	defer s.Finalize()
	err = s.Exec(args...)
	if err != nil {
		return nil, err
	}

	for s.Next() {
		results, err := s.ResultsAsMap()
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, results)
	}

	return
}

func (c *Conn) insert(tableName string, properties map[string]interface{}) (int, os.Error) {
	var keys []string
	var placeholders []string
	var args []interface{}

	for key, val := range properties {
		keys = append(keys, key)
		placeholders = append(placeholders, "?")
		args = append(args, val)
	}

	statement := fmt.Sprintf("insert into %v (%v) values (%v)",
		tableName,
		strings.Join(keys, ", "),
		strings.Join(placeholders, ", "))

	err := c.conn.Exec(statement, args...)
	if err != nil {
		return -1, err
	}

	s, err := c.conn.Prepare("select last_insert_rowid()")
	if err != nil {
		return -1, err
	}

	defer s.Finalize()
	err = s.Exec()
	if err != nil {
		return -1, err
	}

	id := -1

	if s.Next() {
		err := s.Scan(&id)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}

func (c *Conn) Save(rowStruct interface{}) os.Error {
	results, _ := scanStructIntoMap(reflect.NewValue(rowStruct))
	tableName := getTableName(rowStruct)

	id := results["id"]
	results["id"] = 0, false

	if id == 0 {
		id, err := c.insert(tableName, results)
		if err != nil {
			return nil
		}

		structPtr := reflect.NewValue(rowStruct).(*reflect.PtrValue)
		structVal := structPtr.Elem().(*reflect.StructValue)
		structField := structVal.FieldByName("Id")
		structField.SetValue(reflect.NewValue(id))

		return nil
	}

	var updates []string
	var args []interface{}

	for key, val := range results {
		updates = append(updates, fmt.Sprintf("%v = ?", key))
		args = append(args, val)
	}

	statement := fmt.Sprintf("update %v set %v where id = %v",
		tableName,
		strings.Join(updates, ", "),
		id)

	return c.conn.Exec(statement, args...)
}

func (c *Conn) Get(rowStruct interface{}, condition interface{}, args ...interface{}) os.Error {
	conditionStr := ""

	switch condition := condition.(type) {
	case string:
		conditionStr = condition
	case int:
		conditionStr = "id = ?"
		args = append(args, condition)
	}

	conditionStr = fmt.Sprintf("where %v", conditionStr)

	resultsSlice, err := c.getResultsForQuery(getTableName(rowStruct), conditionStr, args)
	if err != nil {
		return err
	}

	switch len(resultsSlice) {
	case 0:
		return os.NewError("did not find any results")
	case 1:
		results := resultsSlice[0]
		scanMapIntoStruct(reflect.NewValue(rowStruct), results)
	default:
		return os.NewError("more than one row matched")
	}

	return nil
}

func (c *Conn) GetAll(rowsSlicePtr interface{}, condition string, args ...interface{}) os.Error {
	rowsPtrValue, _ := reflect.NewValue(rowsSlicePtr).(*reflect.PtrValue)
	rowsPtrType, ok := reflect.Typeof(rowsSlicePtr).(*reflect.PtrType)
	if !ok {
		return os.NewError("needs a *pointer* to a slice")
	}

	sliceValue, _ := rowsPtrValue.Elem().(*reflect.SliceValue)
	sliceType, ok := rowsPtrType.Elem().(*reflect.SliceType)
	if !ok {
		return os.NewError("needs a pointer to a *slice*")
	}

	sliceElementType := sliceType.Elem()

	condition = strings.TrimSpace(condition)
	if len(condition) > 0 {
		condition = fmt.Sprintf("where %v", condition)
	}

	resultsSlice, err := c.getResultsForQuery(getTableName(rowsSlicePtr), condition, args)
	if err != nil {
		return err
	}

	for _, results := range resultsSlice {
		newValue := reflect.MakeZero(sliceElementType)
		scanMapIntoStruct(newValue.Addr(), results)
		sliceValue.SetValue(reflect.Append(sliceValue, newValue))
	}

	return nil
}
