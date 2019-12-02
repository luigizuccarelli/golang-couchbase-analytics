# Go simple service

Simple utility for testing web services on [Microlib](https://github.com/microlib).
Uses a generic script `script.sh` to simulate start & stop for Linux & MacOS based systems. 

## Usage 

```bash
# cd to project directory and build executable
$ go build .
$ chmod u+x script.sh

# start the service
$ ./script.sh start

# stop the service
$ ./script.sh stop
```

Replace the `EXEC` variable in `script.sh` with the name of your executable, if it's different.

## Note
The http server by @luigizuccarelli uses signals to allow for graceful shutdown. Use this as a standard pattern when creating all web services. 

## Queries

### All affiliates and campaigns
SELECT DISTINCT  `utm_affiliate`,`utm_campaign` from analytics ;

### All sources and ad variants
SELECT DISTINCT `utm_source`,`utm_content` from analytics where `utm_affiliate` = 'SBR-01' and `utm_campaign` = 'WinBig' ;

### Pages
SELECT DISTINCT  `from`.`pagename` as source ,`to`.`pagename` as destination from analytics where `utm_affiliate` = 'SBR-01' and `utm_campaign` = 'WinBig' group by `from`.`pagename`, `to`.`pagename` ;

### Paths
SELECT `from`.`pagename` as source, `to`.`pagename` as destination,`trackingid`  from analytics where `utm_source` = 'google' and `utm_content` = 'advariantA' ;

### TrackingId per path
SELECT DISTINCT `trackingid`  from analytics where `utm_source` = 'google' and `utm_content` = 'advariantA' ;

### Count per node
SELECT `to`.`pagename` , `to`.`pagetype`, count(*) from SBR where `event`.`type` = 'load' and utm_campaign = 'WinBig' and utm_source = 'mail' and `to`.`pagename` in ['page-one', 'page-six', 'page-seven'] group by `to`.`pagename` , `to`.`pagetype`
