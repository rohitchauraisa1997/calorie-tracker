db = db.getSiblingDB("calorie-tracker-db");
db.createUser({
    "user": "calorie-user",
    "pwd": "calorie123",
    "roles": [{
            "role": "readWrite",
            "db": "calorie-tracker-db"
        },
        {
            "role": "dbAdmin",
            "db": "calorie-tracker-db"
        }
    ]
});