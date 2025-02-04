# Mysql Connector package

This package connects to MySQL.

The configuration required is to be set in the env variable `MYSQL_CONFIG` in the format:

```json
{
    "Username": "your_username",
    "Password": "your_password",
    "Host": "localhost",
    "Port": 3306,
    "Database": "your_database"
}
```