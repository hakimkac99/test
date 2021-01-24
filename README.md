# Test technique stage chez DataImpact

## Build the project 
    
    go build 

## Run the CRUD project

    ./test-CRUD

## Users files 

    /files

## API Documentation 

### 1. Create (POST /api/add/users)

- Add users to MongoDB database, wrangling data (normalize the age attribute), hash password, create a file for each user containing data. (the name of the file is the id of user) and not inserting the existing users.

Payload example :

    [
        {
            "id": "1t5VsIBXpGl4s8C4CAXTsAlIZISYEOlicj14obz3CwFXCCvaRyuhDI10fah5IfdMS3VblW51my8xt6aQvJI3qNg5as0yqoTCvdZd",
            "password": "CGUsfQkS06lo7LM2guSV",
            "isActive": false,
            "balance": "$2,547.50",
            "age": "27",
            "name": "Nikki Farley",
            "gender": "female",
            "company": "ANIVET",
            "email": "nikkifarley@anivet.com",
            "phone": "+1 (868) 439-2675",
            "address": "588 Schaefer Street, Falconaire, Missouri, 9457",
            "about": "Commodo minim fugiat est fugiat sunt duis consectetur fugiat Lorem sunt. Dolor nulla nisi quis deserunt occaecat tempor reprehenderit. Amet sunt incididunt cillum magna id. Mollit reprehenderit irure non sunt aute consectetur adipisicing. Pariatur veniam proident ut nulla culpa eiusmod ipsum occaecat.\r\n",
            "registered": "2014-10-12T04:14:20 -02:00",
            "latitude": -10.647121,
            "longitude": 126.04006,
            "tags": [
            "consectetur",
            "ipsum",
            "dolor",
            "consectetur",
            "laborum",
            "incididunt",
            "consequat"
            ],
            "friends": [
            {
                "id": 0,
                "name": "Petty Warren"
            },
            {
                "id": 1,
                "name": "Elma Schneider"
            },
            {
                "id": 2,
                "name": "Duffy Mcbride"
            }
            ],
            "data":"G7UxLg4fDL3wifdVAEvRL9EcRyy92FyZCnmLtynifMFcfouQ3TKEKx8HSolWWzbupM2vHIL64khX1yGxYnF7lcNiC02KyF1OdXjJDjbrGTIhpMP1HKuUNQti3dpT12UwzECxiIvyVHKMZiS8PerTEgnqD9CmemQxTKlZvaWJAK569hvQpBzTkUnZETjxs5y66mFyXpu7lvtVXpX6NLfwuKHFLaUycImMuqH97u0njuyg7SGHYFuBYBuYcYSE5thvkXn4AZXrCtFSs5r8pGYBX96UTXnLboiVqR5ph2KyclGtHkMsgFQCrs2wirdmrWbtxY5Y2SHHyjwfRXCGYuy9MtpQDjgjjjf3Euobb6ypcYw6IJ7BUkiNrRCwlZkI2ky9KiCii4aw8jHrb6D5ZUijoTf9MUB7qfnH13ijp6596NWePxH0caGelaTwXUzswKVXyGXDf65CgxUffKrcqPFGCRPFZFFfKVbJOlj3CJPcW0pGMTt6Uyt69xicHIpv4ty0FQNVHCxdklVePXwkJp63lggkBx5rEp0lSvWukbbXefna8oMsjP3twu11AnPSUECX8mdemqy85LUHuVNCICBaTzYofkQtORGYUdjUq9r4BSwk5Tvqn9oHok6JcVLSLnQHVvxBBoS3a3sQQ4yuwdUBRq4c5OyUvWErDNvWKGo97kFQE7vNFq1GIouenZj22sfcEKzzq84aWjP6KsNoh77wuk4xIupWpQvuUU7LWCTdzy3eDoJLxOu0RdprtMKOxPrPaWtTdYJx5kYOkVMQgy2KGx06Jule3j5rF575cLBTiLXb1ag6nCUbo3Xs7PhUPz0EcFlc1QwTWnqIPs7XRJrMNq4V6nEpOOPFpx2LFE"
        },
        {
            "id": "mAdo98L1nvCzdw4CanEmah8PrPqzFNQmO509HrdKCxcLEcluy1zwm9PLvSIOWhZxDFtnM3rLvjK4cVKVpC0BjYbeU6KpDe14Eh0T",
            "password": "VMRmmHi0NvcGHQHI5wrA",
            "isActive": false,
            "balance": "$2,076.63",
            "age": 31,
            "name": "Daugherty Finch",
            "gender": "male",
            "company": "EXOVENT",
            "email": "daughertyfinch@exovent.com",
            "phone": "+1 (949) 540-2258",
            "address": "806 Wortman Avenue, Urie, Puerto Rico, 8196",
            "about": "Ullamco culpa minim pariatur commodo consectetur labore adipisicing anim in voluptate. Pariatur anim laboris nostrud adipisicing dolor est dolor ullamco amet. Commodo amet dolor in nulla laborum laborum culpa et eiusmod. Consequat enim ad nulla ad excepteur elit excepteur non adipisicing eiusmod. Pariatur dolore ea aliquip incididunt anim magna.\r\n",
            "registered": "2015-06-27T01:48:28 -02:00",
            "latitude": -80.569767,
            "longitude": -119.440838,
            "tags": [
            "enim",
            "veniam",
            "irure",
            "laboris",
            "commodo",
            "deserunt",
            "laborum"
            ],
            "friends": [
            {
                "id": 0,
                "name": "Sonia Silva"
            },
            {
                "id": 1,
                "name": "Ford Stewart"
            },
            {
                "id": 2,
                "name": "Mendoza Andrews"
            }
            ],
            "data":"5FcjywRgchSVwBKyVEXXbXzPlw6ONOr5NBn9CyscF2XaNDkiuVb4UU080zeAWpBM4USK4VEDXjipKB11ztqo1wBMHMm3vlFRqIWyISPcOCZNWoIZphGD1JqnykBQ6Ra6QKSwVJ9w8I3R6fztjmhVg25QiV6IaKnK8gSCc3hczZyjq8yUMRzo1eHvqdNlyWkjEGfZcEFyeEMtsvseOtlv1UxpQ8BRLrUwGj0Pn8wkfXlRXv6lYjY7k5265bBOOuhmVeS6fdTlJsFR8ZpAaxpdfaAorSObuNJuwJ9WgCrNNOABVImrkK4WDFf795VAYw85sMgzrRqrcs2tWZG2xR3qzG6ZBjzKfwlgHX5JVjwRwFO0mUefp9fo8e3hqkWpLs5dYrPkzZafSBJ0lB5HZia6qxDd5Bo1iEM5WbbxfJFJGVb94C5CswLq9d7eJL0jZfGnZJap4V0bAMxt4nxZgjm4VTGgVAVfPz1ZgZmFsi2vnssNSguHXMaJX6Aab7V6cq63PLMDxuHUDu1emf4JgXZgPLZnuayg0JZlUX8BOvhfTrkV3MsIbszLWeNrlj4vIki5iKsyXpTV2iFR1HhOUsmglatlufCSvydoxZKHylthzXDSQU0xAAqO4dG7PJLggZPrrmjv24MHGlP7aRtW3KeRuMRxW8IFGmhFS3pmJtwgqGOW5A7NVQSHuulaF0ZtDZMm9l3omHNYqvIJZscBQXVJE6ruJX9tCeq8bsjV8H2H7xzSzHwdX5MkdBIjh5ecyC08rkuO3lgR9B8z0ZEdZs6WPpo6P549uuXVEoG6l1TYv87H2GjDEInM0D9bqhanlVbKCdvsO95acKoDcqHFDUbJcszXprbVQJMhuwjB1iA1oTbqDEXqssNhYwMM5JWhU3YbGlwxyjok4FkkMCyOf0ajL23GzldEv7CZIcJwXZn"
        }
    ]

### 2. Login (Post /api/login)

- Login user 

payload example :

    {
        "id" : "test",
        "password" : "test"
    } 


### 3. Delete (DELETE /api/delete/user/:id)

- Delete the user and the file.


### 4. Read (GET /api/users/list)

- Retrieve all users. 


### 5. Read (GET /api/user/:id)

- Retrieve user with id . 


### 6. Update (UPDATE /api/user/:id)

- update user data with id. (If existing changes in data attribute, then update the file)

payload example :

    {
        "data" : "New data"
    }

## Environment 

- Libs : 
    github.com/gin-gonic/gin
    golang.org/x/crypto
    gopkg.in/mgo.v2

- Developement OS : 
    Windows 10

- Developemnt IDE :
    VS Code

- DB name : go-mongo

- Collection name : user 


