# Sagittarius

This service implements a simple CRUD to manage users callback for asynchronous responses.

We implements a service interface, that let us communicate other services to make asyncrouns calls to other services over internet.

### Service

```
service Saggitarius {
  rpc Create (Callback) returns (Response);
  rpc Update (Callback) returns (Response);
  rpc Delete (Callback) returns (Response);
  rpc Dispatch (Request) returns (Response);
}

```

This service must be implemented by our programa in order to connect postgres to retrieve callbacks to deliver responses to our register users.

#### Service message

```
message Callback {
  string id = 1;
  string user_id = 2;
  string url = 3;
}

message Request {
  string id = 1;
}

message Response {
  string status = 1;
  string message = 2;
}

```
##### Create
##### Read
##### Update
##### Delete
