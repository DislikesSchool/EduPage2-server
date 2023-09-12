# Edupage

You begin with creating a `Credentials` struct

```golang
credentials, err := Login(username, password, server)
if err != nil {
    //Proper error handling...
}
```
This struct contains the important authorization and server information.

Then you can create the `EdupageClient` struct

```golang
client, err := CreateClient(credentials)
if err != nil {
    //Proper error handling...
}
```
This is the most important part, the `EdupageClient` struct. You will use this to interact with the Edupage API.

To update credentials of an existing client use `EdupageClient#UpdateCredentials(Credentials)`
## User
To retrieve the user information use function `EdupageClient#GetUser(bool)`. If the boolean is set to true, this function
will update the stored used structure and, otherwise return the said structure without updating it.
```golang
user, err := client.GetUser(true)
if err != nil {
    //Proper error handling...
}
```

## Timeline
To retrieve timeline information use function `EdupageClient#GetRecentTimeline` or `EdupageClient#GetTimeline(Time, Time)`

```golang
timeline, err := client.GetRecentTimeline()
if err != nil {
    //Proper error handling...
}
```
This will retrieve the last 30 days from the timeline.

Or you can specifiy your own interval, like so

```golang
timeline, err := client.GetTimeline(TIME_FROM, TIME_TO)
if err != nil {
    //Proper error handling...
}
```

## Results (grades)
To retrieve the results use function `EdupageClient#GetRecentResults`

```golang
results, err := client.GetRecentResults()
if err != nil {
	//Proper error handling...
}
```

This will load results from the current year, you can also specify your own interval
using `EdupageClient#GetResults(string, string)`

The arguments are halfyear and year.
Possible halfyears are (but are not limited to) these values:
- P1 first half
- P2 second half
- RX whole year

## Timetable (dayplan)
To retrieve the timetable use function `EdupageClient#GetRecentResults`
```golang
timetable, err := client.GetRecentTimetable()
if err != nil {
    //Proper error handling...
}
```

This will load timetable 7 days ahead, and 2 days before today (using local time). You can also specify your own interval using `EdupageClient#GetResults(Time, Time)`.

