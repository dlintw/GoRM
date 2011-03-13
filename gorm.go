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
	conn, err := sqlite.Open("test.db")
	return &Conn{conn: conn}, err
}

func getTableName(obj interface{}) string {
	return pluralizeString(snakeCasedName(getTypeName(obj)))
}

func (c *Conn) getResultsForQuery(tableName, condition string) (resultsSlice []map[string][]byte, err os.Error) {
    s, err := c.conn.Prepare(fmt.Sprintf("select * from %v %v", tableName, condition))
    if err != nil {
		return nil, err
    }

    defer s.Finalize()
    err = s.Exec()
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

func (c *Conn) Save(rowStruct interface{}) os.Error {
	results, _ := scanStructIntoMap(reflect.NewValue(rowStruct))
	
	id := results["id"]
	results["id"] = 0, false
	
	// log.Fatalf("%v\n", id)
	var updates []string
	
	for key, val := range results {
		escStr, err := escapeString(fmt.Sprintf("%v = ?", key), val)
		if err != nil {
			return err
		}
		updates = append(updates, escStr)
	}
	
	updatesStr := strings.Join(updates, ", ")
	
    s, err := c.conn.Prepare(fmt.Sprintf("update %v set %v where id = %v", getTableName(rowStruct), updatesStr, id))
    if err != nil {
		return err
    }

    defer s.Finalize()
    err = s.Exec()
    if err != nil {
        return err
    }

	return nil
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
	
	conditionStr, err := escapeString(conditionStr, args...)
	if err != nil {
		return err
	}
	
	conditionStr = fmt.Sprintf("where %v", conditionStr)
	
	resultsSlice, err := c.getResultsForQuery(getTableName(rowStruct), conditionStr)
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
	
	condition, err := escapeString(condition, args...)
	if err != nil {
		return err
	}
	
	condition = fmt.Sprintf("where %v", condition)
	
	resultsSlice, err := c.getResultsForQuery(getTableName(rowsSlicePtr), condition)
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
