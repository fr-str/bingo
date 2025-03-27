# set session cookie
GET http://localhost:2138/
HTTP 200
[Asserts]
cookie "session" exists

# set timesheets crash
GET http://localhost:2138/api/square/click
[Query]
field: timesheets crash
HTTP 302
Location: /

# check if set
GET http://localhost:2138/api/stats
HTTP 200
Content-Type: application/json
[Asserts]
jsonpath "$.[0].Count" != null
jsonpath "$.[0].Count" == 1
jsonpath "$.[0].Field" == "timesheets crash"
jsonpath "$.[0].Date" != null
jsonpath "$.[0].Date" != ""

# set timesheets crash again
GET http://localhost:2138/api/square/click
[Query]
field: timesheets crash
HTTP 302
Location: /

# should not be set
GET http://localhost:2138/api/stats
HTTP 200
Content-Type: application/json
[Asserts]
bytes count == 5

# check if bingo-cells are rendered
GET http://localhost:2138/
HTTP 200
Content-Type: text/html; charset=utf-8
[Asserts]
xpath "//div[@class='bingo-cell']" count != 0

# set timesheets crash again
GET http://localhost:2138/api/square/click
[Query]
field: timesheets crash
HTTP 302
Location: /

# check if timesheets crash is set
GET http://localhost:2138/
HTTP 200
Content-Type: text/html; charset=utf-8
[Asserts]
xpath "//div[@class='bingo-cell selected']" count == 1
xpath "string(//div[@class='bingo-cell selected'])" == "timesheets crash"
