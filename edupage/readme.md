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


