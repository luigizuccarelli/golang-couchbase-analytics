# Go aggregation service for sankey charts & analytics

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

### Sankey in time order
SELECT `page`.`referrer` AS source,
       `page`.`url` AS destination,
       COUNT(`journey_id`) AS count,
       `TEST`.`timestamp` as ts
FROM TEST
WHERE `spec` = 'page'
GROUP BY `TEST`.`timestamp` AS ts,
         `page`.`referrer` AS source,
         `page`.`url` AS destination
ORDER BY `TEST`.`timestamp`

### Sankey old
"select `from`.`pagename` as source,`to`.`pagename` as destination, count(`trackingid`) as count  from " +
vars["affiliate"] + " where `event`.`type` = 'load'  group by `from`.`pagename` as source, `to`.`pagename` as destination"


