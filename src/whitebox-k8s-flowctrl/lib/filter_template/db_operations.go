package filter_template

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// db operations

var (
	ErrDBConnectionInfoConfigMapHasInvalidKey = errors.New("DB connection info ConfigMap has invalid key")
	ErrDBConnectionInfoSecretHasInvalidKey    = errors.New("DB connection info Secret has invalid key")
)

type db_ConnectionInfoStruct struct {
	DriverName string
	IP         string
	Port       string
	DB_Name    string
	RoleName   string
	Password   string
	DB         *sql.DB
}

// Open DataBase and return
func (dbConInfo *db_ConnectionInfoStruct) Open() (err error) {
	openOption := fmt.Sprintf(`host=%v
						port=%v
						user=%v
						password=%v
						dbname=%v
						sslmode=disable`, dbConInfo.IP, dbConInfo.Port, dbConInfo.RoleName, dbConInfo.Password, dbConInfo.DB_Name)
	dbConInfo.DB, err = sql.Open(dbConInfo.DriverName, openOption)
	return err
}

func (dbConInfo *db_ConnectionInfoStruct) Close() error {
	return dbConInfo.DB.Close()
}

func (dbConInfo *db_ConnectionInfoStruct) Query(q string) ([]map[string]interface{}, error) {
	rows, err := dbConInfo.DB.Query(q)
	if err != nil {
		var tmpRes []map[string]interface{}
		return tmpRes, err
	}
	return rowsToMapSlice(rows)
}

func fetchDB_ConnectionInfo(
	r client.Reader,
	ctx context.Context,
	dbInfoConfigMapName string,
	dbInfoConfigMapNamespaceCandidates []string,
	roleName string,
	dbRoleInfoSecretNamespaceCandidates []string,
) (*db_ConnectionInfoStruct, error) {

	l := log.FromContext(ctx)

	var dbConInfo db_ConnectionInfoStruct

	DB_ConnectionInfoConfigMap, err := tryFetchConfigMapFromSeveralNameSpaceCandidates(
		r, ctx, dbInfoConfigMapName,
		dbInfoConfigMapNamespaceCandidates,
	)

	if err != nil {
		return &dbConInfo, err
	}

	for key, val := range DB_ConnectionInfoConfigMap.Data {
		switch key {
		case "IP":
			dbConInfo.IP = val
		case "Port":
			dbConInfo.Port = val
		case "DB_Name":
			dbConInfo.DB_Name = val
		case "DB_DriverName":
			dbConInfo.DriverName = val
		default:
			l.Error(ErrDBConnectionInfoConfigMapHasInvalidKey, fmt.Sprintf("Invalid ConfigMap Paramter is inputted. key : %v", key))
			return &dbConInfo, ErrDBConnectionInfoConfigMapHasInvalidKey
		}
	}

	dbRoleInfoSecret := &corev1.Secret{}
	err = tryFetchResourceFromSeveralNameSpaceCandidates(
		r, ctx, fmt.Sprintf("db-role-info-%v", roleName),
		dbRoleInfoSecretNamespaceCandidates,
		dbRoleInfoSecret,
	)

	if err != nil {
		return &dbConInfo, err
	}

	dbConInfo.RoleName = roleName
	for key, val := range dbRoleInfoSecret.Data {
		switch key {
		case "Password":
			dbConInfo.Password = string(val)
		default:
			l.Error(ErrDBConnectionInfoSecretHasInvalidKey, "Invalid Secret Paramter is inputted")
			return &dbConInfo, ErrDBConnectionInfoSecretHasInvalidKey
		}
	}

	return &dbConInfo, nil
}

func rowsToMapSlice(rows *sql.Rows) ([]map[string]interface{}, error) {

	var results []map[string]interface{}
	cols, err := rows.Columns()

	if err != nil {
		return results, err
	}

	for rows.Next() {
		var row = make([]interface{}, len(cols))
		var rowp = make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			rowp[i] = &row[i]
		}

		rows.Scan(rowp...)

		rowMap := make(map[string]interface{})
		for i, col := range cols {
			switch row[i].(type) { //nolint:gocritic // 'type switches' is simple way to compare type
			case []byte:
				row[i] = string(row[i].([]byte))
				num, err := strconv.Atoi(row[i].(string))
				if err == nil {
					row[i] = num
				}
			}
			rowMap[col] = row[i]
		}

		results = append(results, rowMap)
	}
	return results, nil

}
