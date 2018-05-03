# NeoTrade-MicroloansFilter
Batch code written in Golang which uses Prosper API to find microloans that match loan criteria.

It also uses Twilio API to send  text messages on loans found.

Used the following standard library go packages:
* time - displaying dates as well as manipulating date/time objects
* os - for operating system functionality such as accessing command line arguments and outputing to standard error file descriptors etc.
* strconv - only to code integer to string object conversions
* encoding/json - to implement unmarshalling of JSON data
