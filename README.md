# jump-run
Very basic go application that maintains running average

#####Add Action

This endpoint adds an action in the form of 

```json
{
      "action" : "jump",
      "time" : 100
}
```
* Endpoint : POST /action
* action is a string representation of an action.
* time is an int that represents time. Value in the range [0-10000].
* _A negative time or malformed json request will return a bad request._

#####Get Stats Action
A GET request for /stats will return the average for all of the actions. This endpoint will also accept a param in the form of
```html
/stats?action=jump
```
which will only output the average for that action.


#####Delete Action
This endpoint deletes an action in the form of
```json
{
      "action" : "jump"
}
```
* Endpoint : DELETE /delete
* action is a string representation of an action. An empty action will clear all of the action

##Installation

For installation of Go. https://golang.org/doc/install

To Run the application
```html
go build main.go

```