# Taxi Data

Data is taken and modified from _Distributed Event-Based Systems conferences and workshops (DEBS)_.
The challenge is from 2015 "Taxi Trips" http://debs.org/debs-2015-grand-challenge-taxi-trips/

The file contains types of events:
* Taxi picks up a passenger
* Taxi drops passenger off

## Event Schema



|Field|Type|Info|
|-|-|-|
|Date|String|2013-01-01 00:21:00|
|Type|String|`pickup` or `dropoff`|
|TaxiID|String|taxi-x|
|LicenseID|String|license-n|
|Latitude|float|40.753086|
|Longitude|float|-73.978867|
|Charge|float|8.62 in dollars, only for type `dropoff`|
|Tip|float|1.62 in dollars, only for type `dropoff`|
|Duration|int|240 in seconds, only for type `dropoff`|
|Distance|float|3.23 in miles, only for type `dropoff`|
