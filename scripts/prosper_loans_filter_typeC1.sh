#!/bin/sh
##add this to crontab to autorun:
## 30 2 * * 1-5  cd /root/NeoTrade-MicroloansFilter/scripts && sh prosper_loans_filter_typeC.sh 2>> /tmp/neotrade_prosper_loans_filter.err 1>> /tmp/neotrade_prosper_loans_filter.log

cd ../bin
./prosper_loans_filter "../src/prosper_marketplace/creds/creds.json" "../src/prosper_marketplace/filters/filter_typeC1_8.501min.json" "+12482472871"
