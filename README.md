# PesenTiketOnlineGolang

This is a golang app about online ticket buy with goroutine and http. It runs on port 8080 and have 15 end points that can be used.

## First this is the endpoint to get all event and tickets

![Event Get All Screenshot](./images/EventGetAll.png)

## Next is the andpoint to add user and view users

![User Add](./images/UserAdd.png)
![User View All](./images/UserViewAll.png)

## Last is the endpoint for buy and creating order transaction details

![Create Order](./images/CreateOrder.png)

## The respond json should be like this.

### If the order status success then the ticket stock and also the balance of the user is reduced

![Result of Create Order](./images/ResultCreateOrder.png)

## The log view should be like this everytime you hit an endpoint it will show whether it success or not

![Log View](./images/LogView.png)

## Next is the update where I added database for storing events, users, and orders (transaction details)
## Event Endpoints
### First is the endpoint to create an Event. It will looks like this
![Event Create DB](./images/using_database/event_api/eventCreate.png)

### Next is the endpoint to View Events (All events and by id)
![Event View All DB](./images/using_database/event_api/eventGet.png)
![Event View by ID](./images/using_database/event_api/eventGetById.png)

### Next is the endpoint to Add a new ticket to the event
![Event Add New Ticket DB](./images/using_database/event_api/eventAddTicket.png)

### Last is the endpoint to delete an event in database
![Event Delete DB](./images/using_database/event_api/eventDelete.png)

## User Endpoints
### First is the endpoint to create User. It will looks like this
![User Create DB](./images/using_database/user_api/userCreate.png)

### Next is the endpoint to View Users (All users and by id)
![User View All DB](./images/using_database/user_api/userGet.png)
![User View by id DB](./images/using_database/user_api/userGetById.png)

### Next is the endpoint to update user's balance
![User Update Balance DB](./images/using_database/user_api/userUpdateBalance.png)

### Last is the endpoint to delete user DB
![User Delete DB](./images/using_database/user_api/userDelete.png)

## Order (Transaction Detail) Endpoints
### First is the endpoint to create an Order. It will looks like this
![Order Create Request DB](./images/using_database/order_api/orderCreateRequest.png)
![Order Create Response DB](./images/using_database/order_api/orderCreateResponse.png)

### Next is the endpoint to View Order (All orders and by User ID)
#### This is Get All Orders
![Order View All 1 DB](./images/using_database/order_api/orderView1.png)
![Order View All 2 DB](./images/using_database/order_api/orderView2.png)

#### This is Get Orders by User ID
![Order View By ID 1 DB](./images/using_database/order_api/orderViewByUserId.png)
![Order View By ID 1 DB](./images/using_database/order_api/orderViewByUserId2.png)

## This is What The log looks like
### There are 2 log, the success one and the failed(error) one. The only difference is that the failed log (log error) is telling us the error messages, while the success one is not
![Log View 2](./images/using_database/LogView2.png)

### There are other end point that can be used to perform CRUD to User and Events. You can see it on file main.go
