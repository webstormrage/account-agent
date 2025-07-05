package accountdb

import (
	appContext "account-agent/pkg/appContext"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Account struct {
	Name string
	Value int
	Currency string
	Id int
}

func GetAccountsByCurrency(currency string)([]Account, error) {
	ctx := appContext.Get()
	db, err := sql.Open("postgres", ctx.DataSourceName)
    if err != nil {
      return nil, err
    }
	defer db.Close()
	rows, err := db.Query("SELECT * FROM accounts WHERE currency=$1 ORDER BY id", currency)
	if err != nil {
		return nil, err
	}
	accounts := []Account{}
	for rows.Next() {
		var account Account
		err = rows.Scan(&account.Name, &account.Value, &account.Currency, &account.Id)
		if err != nil {
		  return nil, err
	    }
		accounts = append(accounts, account)
	}
	return accounts, nil
}

type Report struct{
	Total int
	Daily int
	Days int
}

func getDayEnd()time.Time{
    dayStartTime := time.Now()
	dayStartTime = dayStartTime.AddDate(0, 1, 0)
	dayStartTime = time.Date(dayStartTime.Year(), dayStartTime.Month(), 1,
		6, 0, 0, 0, dayStartTime.Location())
	return dayStartTime	
}

func GetReportByCurrency(currency string)(*Report, error) {
	ctx := appContext.Get()
	db, err := sql.Open("postgres", ctx.DataSourceName)
    if err != nil {
      return nil, err
    }
	defer db.Close()
	query := "SELECT " +
    "SUM(value) AS total, " +
    "SUM(value)/(SELECT DATE '2025-08-01' - CURRENT_DATE AS days) AS daily, " + 
    "(SELECT DATE '"+ getDayEnd().Format("2006-01-02")  + "' - CURRENT_DATE AS days_until) as days " + 
    "FROM accounts WHERE currency = $1"
	row := db.QueryRow(query, currency)
	var report Report
    err = row.Scan(&report.Total, &report.Daily, &report.Days)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &report, err
}

func UpdateRecord(name string, value int64)error{
	ctx := appContext.Get()
	db, err := sql.Open("postgres", ctx.DataSourceName)
    if err != nil {
      return err
    }
	defer db.Close()
	// findQuery := `SELECT id FROM accounts WHERE name=`+name;
	// fmt.Printf("Поиск: %s\n", findQuery)
	row := db.QueryRow(`SELECT id FROM accounts WHERE name=$1`, name)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return err
	}
	fmt.Printf("Найдено совпадение: %d %s\n", id, name)
	query := "UPDATE accounts " +
    "SET value=$1 " +
    "WHERE id=$2"
	_, err = db.Exec(query, value, id)
	if err != nil {
		return err
	}
	return nil
}