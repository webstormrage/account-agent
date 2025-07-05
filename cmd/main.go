package main

import (
	accountdb "account-agent/pkg/account-db"
	"account-agent/pkg/appContext"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"os"
)

func formatMoney(value int, currency string)string{
	if value == 0 {
		return "0 " + currency
	}
	parts := []string{}
	for value > 0 {
		if value >= 1000 {
			parts = append(parts, fmt.Sprintf("%03d", value % 1000))
		} else {
			parts = append(parts,  fmt.Sprintf("%d", value % 1000))
		}
		
		value /= 1000
	}
	slices.Reverse(parts)
	return strings.Join(parts, " ") + " " + strings.ToUpper(currency)
}

func PrintAccounts(accounts []accountdb.Account){
	for _, a := range accounts {
		fmt.Printf("%-20s %-20s\n",  strings.ToUpper(a.Name), formatMoney(a.Value, a.Currency))
	}
}

func HandleEntry(query string){
	entries := strings.Split(query, " ")
	fmt.Printf("Обновляю: %s\n", query);
	value, err := strconv.ParseInt(entries[0], 10, 64)
	if err != nil {
		fmt.Printf("Ошибка обновления: %w\n", err);
		return
	}
	name := entries[1]
	err = accountdb.UpdateRecord(name, value)
	if err != nil {
		fmt.Printf("Ошибка обновления: %w\n", err);
		fmt.Printf("%s %d\n", name, len(name));
	} else {
		fmt.Printf("Обновлено: %d %s\n", value, name);
		fmt.Printf("%s %d\n", name, len(name));
	}
}

func HandleAccountList(currency string){
	accounts, err := accountdb.GetAccountsByCurrency(currency)
	if err !=  nil {
		panic(err)
	}
	fmt.Println("\n")
	PrintAccounts(accounts)
}

func HandleAccountReport(currency string){
	reportAmd, err := accountdb.GetReportByCurrency(currency)
	if err !=  nil {
		panic(err)
	}
	fmt.Println("\n")
	fmt.Printf("%-15s %-15s %-10s\n", "СУММА", "НА ДЕНЬ", "ДНЕЙ")
	fmt.Printf(
		"%-15s %-15s %-10s\n",
		  formatMoney(reportAmd.Total, currency),
		   formatMoney(reportAmd.Daily, currency),
		    strconv.Itoa(reportAmd.Days) + " д.",
		)
}

func HandleCommand(command string){
	entries := strings.Split(command, " ")
	currency := entries[0]
	operation := entries[1]
	switch operation {
	case "*":
		HandleAccountList(currency)
		break;
	case "+":
		HandleAccountReport(currency)
		break;
	default:
		panic("unknown command")	
	}
}

func main(){
	err := appContext.Init()
    if err !=  nil {
		panic(err)
	}
	fileData, err := os.ReadFile("logs/input")
    if err != nil {
        panic(err) 
    }
	queiries := strings.Split(string(fileData), "\r\n")
	for _, query := range queiries {
		HandleEntry(query)
	}
	formatData, err := os.ReadFile("logs/config")
	if err != nil {
        panic(err) 
    }
	commands := strings.Split(string(formatData), "\r\n")
	for _, command := range commands {
		HandleCommand(command)
	}
}