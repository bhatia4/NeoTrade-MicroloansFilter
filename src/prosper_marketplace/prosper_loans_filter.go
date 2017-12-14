package main

import (
  "fmt"
  "time"
  "io/ioutil"
  "os"
  "encoding/json"
  "github.com/bhatia4/gofn-prosper/prosper"
  "github.com/bhatia4/gofn-prosper/prosper/auth"
  "github.com/bhatia4/gofn-prosper/interval"
)

type ProsperCreds struct {
    ClientID		string `json:"clientID"`
    ClientSecret	string `json:"clientSecret"`
    Username		string `json:"username"`
    Password		string `json:"password"`
}

func main() {
	raw, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

	var creds ProsperCreds
    json.Unmarshal(raw, &creds)

	client := prosper.NewClient(auth.ClientCredentials{
  		ClientID:     creds.ClientID,
  		ClientSecret: creds.ClientSecret,
  		Username:     creds.Username,
  		Password:     creds.Password,
	})

	account, err := client.Account(prosper.AccountParams{})
	if err != nil {
  		fmt.Printf("Failed to retrieve account information: %v", err)
  		return
	}
	fmt.Printf("Your account has $%.2f in cash and a total value of $%.2f\n",
	  	account.AvailableCashBalance, account.TotalAccountValue)

	searchResp, err := client.Search(prosper.SearchParams{
  		Offset: 0,
 		Limit:  20,
  		Filter: prosper.SearchFilter{
    		Rating:      		[]prosper.Rating{prosper.RatingC},
			ListingStartDate: 	interval.NewTimeRange(time.Now().AddDate(0, 0, -5), time.Now()),//going back atmost 2 days
			ListingStatus: 		[]prosper.ListingStatus{prosper.ListingActive},
			DtiWprosperLoan: 	interval.NewFloat64Range(0.0, 0.225),
			EstimatedReturn: 	interval.NewFloat64Range(0.0851, 0.2),
			IncomeRange:		[]prosper.IncomeRange{prosper.Between75kAnd100k, prosper.Over100k},
			ListingTerm:		[]prosper.ListingTerm{prosper.ListingTerm60Months},
  		},
	})
	if err != nil {
  		fmt.Printf("Failed to search available note listings: %v\n", err)
  		return
	}

	fmt.Printf("Found %d matching notes, processing first %d\n",
  		searchResp.TotalCount, searchResp.ResultCount)
//	fmt.Printf("Second filter - remove loans based on income range and listing term\n")

	for i, listing := range searchResp.Results {
		fmt.Printf("%2d: ID:%v; Rating:%s; Status:%s; Amount:$%5.0f; Listed on:%s; Delinquencies last 7yrs:%d; Est Return:%.2f%%; Term:%d; Income Range:%d (%s); Last 6 mos. Inquiries:%d; Debt-to-Income Ratio:%.2f%%; Prior Prosper Loans(Late Payments 1 mo+:%d; Bal. Outstanding:%.2f)\n",
				i+1, 
				listing.ListingNumber, 
				prosper.RatingToString(listing.Rating),
				listing.ListingStatusReason,
				listing.ListingAmount,
				listing.ListingStartDate.Format("2/Jan/06"),
				listing.DelinquenciesLast7Years,
				listing.EstimatedReturn*100.0,
				listing.ListingTerm,
				listing.IncomeRange,
				listing.IncomeRangeDescription,
				listing.InquiriesLast6Months,
				listing.DtiWprosperLoan*100.0,
				listing.PriorProsperLoansLatePaymentsOneMonthPlus,
				listing.PriorProsperLoansBalanceOutstanding,)
	}
}
