#!/bin/sh
##add this to crontab to autorun:
## 30 2 * * 1-5  cd /<path to NeoTrade-MicroloansFilter> && sh run_prosper_loans_filter example.sh 2>> /tmp/neotrade_prosper_loans_filter.err 1>> /tmp/neotrade_prosper_loans_filter.log

cd ../bin
./prosper_loans_filter.exe "../src/prosper_marketplace/creds/creds.json" "../src/prosper_marketplace/filters/filter_typeC1.json" "incoming_phone_number_for_sms"