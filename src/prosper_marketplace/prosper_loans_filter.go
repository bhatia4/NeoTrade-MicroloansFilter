package main

import (
  "time"
  "io"
  "io/ioutil"
  "os"
  "log"
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

type Filter struct {
    ShortDesc		string `json:"shortDesc"`
    LongDesc		string `json:"LongDesc"`
	DtiWprosperLoan	interval.Float64Range `json:"dtiWprosperLoan"`
	EstimatedReturn	interval.Float64Range `json:"estimatedReturn"`
	IncomeRange		[]int8 `json:"incomeRange"`
	Limit			int8   `json:limit`
	DaysFromCurrentTime int8 `json:daysFromCurrentTime`
	ListingStatus 	[]int8 `json:listingStatus`
	ListingTerm		[]int8 `json:listingTerm`
	Offset			int8   `json:offset`
	Rating			[]string `json:rating`
}

type MyFloat64Range interval.Float64Range
func (o *MyFloat64Range) UnmarshalJSON(data []byte) error {
	Trace.Printf("%s\n", data)
	var v [2]float64
	if err := json.Unmarshal(data, &v); err != nil {
		Error.Println(err.Error())
		return err
	}
	Trace.Println("%+v\n", v)
	o.Min = &v[0]
	o.Max = &v[1]
	return nil
}

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func Init(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}


func readFromFile(filename string) []byte {
	raw, err := ioutil.ReadFile(filename)
    if err != nil {
        Error.Println(err.Error())
        os.Exit(1)
    }
	return raw
}

func main() {
	//setup logging
	Init(ioutil.Discard /*when testing, replace w/ os.Stdout*/, 
			os.Stdout, os.Stdout, os.Stderr)
	
	//read creds json file
	var creds ProsperCreds
    json.Unmarshal(readFromFile(os.Args[1]), &creds)

	//read filter json file
	var currFilter Filter
    json.Unmarshal(readFromFile(os.Args[2]), &currFilter)

	Info.Printf("%+v\n", currFilter)
	
	client := prosper.NewClient(auth.ClientCredentials{
  		ClientID:     creds.ClientID,
  		ClientSecret: creds.ClientSecret,
  		Username:     creds.Username,
  		Password:     creds.Password,
	})
	account, err := client.Account(prosper.AccountParams{})
	if err != nil {
  		Error.Printf("Failed to retrieve account information: %v", err)
  		return
	}
	
	Trace.Printf("Your account has $%.2f in cash and a total value of $%.2f\n",
	  	account.AvailableCashBalance, account.TotalAccountValue)


	var inputRatings []prosper.Rating = make([]prosper.Rating, len(currFilter.Rating));
	for index := range currFilter.Rating {
		curr, err := prosper.ParseRating(currFilter.Rating[index])
		if err != nil {
			Error.Println(err.Error())
			os.Exit(1)
		}
		inputRatings[index] = curr
	}
	
	var inputListingStatus []prosper.ListingStatus = make([]prosper.ListingStatus, len(currFilter.ListingStatus));
	for index := range currFilter.ListingStatus {
		curr, err := prosper.ParseListingStatus(int64(currFilter.ListingStatus[index]))
		if err != nil {
			Error.Println(err.Error())
			os.Exit(1)
		}
		inputListingStatus[index] = curr
	}
	
	var inputIncomeRange []prosper.IncomeRange = make([]prosper.IncomeRange, len(currFilter.IncomeRange));
	for index := range currFilter.IncomeRange {
		curr, err := prosper.ParseIncomeRange(int64(currFilter.IncomeRange[index]))
		if err != nil {
			Error.Println(err.Error())
			os.Exit(1)
		}
		inputIncomeRange[index] = curr
	}
	
	var inputListingTerm []prosper.ListingTerm = make([]prosper.ListingTerm, len(currFilter.ListingTerm));
	for index := range currFilter.ListingTerm {
		curr, err := prosper.ParseListingTerm(int64(currFilter.ListingTerm[index]))
		if err != nil {
			Error.Println(err.Error())
			os.Exit(1)
		}
		inputListingTerm[index] = curr
	}
	
	searchResp, err := client.Search(prosper.SearchParams{
  		Offset: 0,
 		Limit:  20,
  		Filter: prosper.SearchFilter{
    		Rating:      		inputRatings,
			ListingStartDate: 	interval.NewTimeRange(time.Now().AddDate(0, 0, int(-1*currFilter.DaysFromCurrentTime)), time.Now()),
			ListingStatus: 		inputListingStatus,
			DtiWprosperLoan: 	currFilter.DtiWprosperLoan,
			EstimatedReturn: 	currFilter.EstimatedReturn,
			IncomeRange:		inputIncomeRange,
			ListingTerm:		inputListingTerm,
  		},
	})
	if err != nil {
  		Error.Printf("Failed to search available note listings: %v\n", err)
  		return
	}

	Info.Printf("Found %d matching notes, processing first %d\n",
  		searchResp.TotalCount, searchResp.ResultCount)

	for i, listing := range searchResp.Results {
		Info.Printf("%2d: ID:%v; Rating:%s; Status:%s; Amount:$%5.0f; Listed on:%s; Delinquencies last 7yrs:%d; Est Return:%.2f%%; Term:%d; Income Range:%d (%s); Last 6 mos. Inquiries:%d; Debt-to-Income Ratio:%.2f%%; Prior Prosper Loans(Late Payments 1 mo+:%d; Bal. Outstanding:%.2f)\n",
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
