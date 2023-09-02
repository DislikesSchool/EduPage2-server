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
This is the most important part, the `EdupageClient` struct takes care of data management and fetching.

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
