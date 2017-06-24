# Postgres Database
Build with:
```
sudo docker build -t some-postgres .
```

Run with:
```
sudo docker run -p 5432:5432 -e POSTGRES_PASSWORD=password postgres
sudo docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
```
