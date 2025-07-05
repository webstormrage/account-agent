package main

import (
	accountdb "account-agent/pkg/account-db"
	"account-agent/pkg/appContext"
	"fmt"
	"slices"
	"strconv"
	"strings"
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

func main(){
	err := appContext.Init()
    if err !=  nil {
		panic(err)
	}
	amds, err := accountdb.GetAccountsByCurrency("amd")
	if err !=  nil {
		panic(err)
	}
	fmt.Println("\n")
	PrintAccounts(amds)
    usds, err := accountdb.GetAccountsByCurrency("usd")
	if err !=  nil {
		panic(err)
	}
	fmt.Println("\n")
	PrintAccounts(usds)
	rubs, err := accountdb.GetAccountsByCurrency("rub")
	if err !=  nil {
		panic(err)
	}
	fmt.Println("\n")
	PrintAccounts(rubs)

	reportAmd, err := accountdb.GetReportByCurrency("amd")
	if err !=  nil {
		panic(err)
	}
	fmt.Println("\n")
	fmt.Printf("%-15s %-15s %-10s\n", "СУММА", "НА ДЕНЬ", "ДНЕЙ")
	fmt.Printf(
		"%-15s %-15s %-10s\n",
		  formatMoney(reportAmd.Total, "amd"),
		   formatMoney(reportAmd.Daily, "amd"),
		    strconv.Itoa(reportAmd.Days) + " д.",
		)
}