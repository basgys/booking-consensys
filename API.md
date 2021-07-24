# POST /auth/challenge

+ Request (application/json; charset=utf-8)

        {
            "address": "0x6c35683E9d599B715B57108296DA14Cb944316AD"
        }

+ Response 201 (application/json; charset=utf-8)

    + Headers

            Node-Name: 
            Request-Id: 94dd9bb0-9f5f-428e-9368-8437f95ffd70
            Node-Version: 

    + Body

            {"challenge":"login:2021-07-26T15:59:55.734247Z:nkaul/GdzSU="}
            


# POST /auth/authorise

+ Request (application/json; charset=utf-8)

        {
            "address": "0x6c35683E9d599B715B57108296DA14Cb944316AD",
            "signature": "0x1856d3a2c69c0153d272b4fc518d59ad15cc6cdc088c6b31221fd71cdac58055099c56b115247af7734b90adbb9a5269eeb61928d80df88cf2db45548b199c281b"
        }

+ Response 403 (application/json; charset=utf-8)

    + Headers

            Node-Name: 
            Node-Version: 
            Request-Id: eea40a7c-85b9-4a84-9ff6-0aed38aaf6be

    + Body

            {"error":{"message":"permission denied"}}
            


# GET /booking/rooms

+ Response 200 (application/json; charset=utf-8)

    + Headers

            Request-Id: b4be7624-0548-45ba-8c2e-9d3375dce9d5
            Node-Version: 
            Node-Name: 

    + Body

            {
                "rooms": [
                    {
                        "ref": "C01"
                    },
                    {
                        "ref": "C02"
                    },
                    {
                        "ref": "C03"
                    },
                    {
                        "ref": "C04"
                    },
                    {
                        "ref": "C05"
                    },
                    {
                        "ref": "C06"
                    },
                    {
                        "ref": "C07"
                    },
                    {
                        "ref": "C08"
                    },
                    {
                        "ref": "C09"
                    },
                    {
                        "ref": "C10"
                    },
                    {
                        "ref": "P01"
                    },
                    {
                        "ref": "P02"
                    },
                    {
                        "ref": "P03"
                    },
                    {
                        "ref": "P04"
                    },
                    {
                        "ref": "P05"
                    },
                    {
                        "ref": "P06"
                    },
                    {
                        "ref": "P07"
                    },
                    {
                        "ref": "P08"
                    },
                    {
                        "ref": "P09"
                    },
                    {
                        "ref": "P10"
                    }
                ]
            }


# GET /booking/rooms/P09/availabilities?from=2021-01-01T00:00:00Z&to=2021-12-31T23:59:59Z

+ Response 200 (application/json; charset=utf-8)

    + Headers

            Node-Name: 
            Node-Version: 
            Request-Id: ef8d5f25-1389-4470-9889-690e9a808e02

    + Body

            {"availabilities":[{"from":"2021-01-01T00:00:00Z","to":"2021-07-26T12:00:00Z"},{"from":"2021-07-26T13:00:00Z","to":"2021-12-31T23:59:59Z"}]}
            


# GET /booking/rooms/c10/reservations

+ Response 200 (application/json; charset=utf-8)

    + Headers

            Request-Id: d602d7c8-f1a6-41a3-aa43-803c5022423a
            Node-Version: 
            Node-Name: 

    + Body

            {"reservations":[{"id":"1vqs9KjdMS2v1zGZr18pqO4KT24","from":"2021-07-26T12:00:00Z","to":"2021-07-26T13:00:00Z","roomRef":"C10","userId":"5f6d4cae-3cd2-4903-a7e6-63a352be444e"}]}
            


# POST /booking/rooms/c10/reservations

+ Request (application/json; charset=utf-8)

    + Headers

            Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjk5MDcxNjgwMDU5NTEwMDAsImp0aSI6IjB4NmMzNTY4M0U5ZDU5OUI3MTVCNTcxMDgyOTZEQTE0Q2I5NDQzMTZBRCJ9.DqNZiwYiXb2nuOudK5k56FLpCql42UM56sogpi5ilyU

    + Body

            {
                "from": "2021-08-01T16:30:00Z",
                "hours": 1
            }

+ Response 201 (application/json; charset=utf-8)

    + Headers

            Request-Id: 8ae08238-540d-4749-adc9-ceef14078154
            Node-Version: 
            Node-Name: 

    + Body

            {"id":"1w8I4sdWEC1QfpHxRzvGSJq0mBy","from":"2021-08-01T16:00:00Z","to":"2021-08-01T18:00:00Z","roomRef":"C10","userId":"611b783e-7138-42bc-b2af-fb69ad810fa3"}
            


# DELETE /booking/rooms/C10/reservations/1vqs9KjdMS2v1zGZr18pqO4KT24

+ Request

    + Headers

            Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjk5MDcxNjgwMDU5NTEwMDAsImp0aSI6IjB4NmMzNTY4M0U5ZDU5OUI3MTVCNTcxMDgyOTZEQTE0Q2I5NDQzMTZBRCJ9.DqNZiwYiXb2nuOudK5k56FLpCql42UM56sogpi5ilyU



+ Response 204

    + Headers

            Request-Id: 4b117d4e-c865-4088-a34c-0c2c39f19f6e
            Node-Version: 
            Node-Name: 




