# Receipt Processor

## Steps to Run the Application

1. Pull this repository 
```
git clone https://github.com/vishisth29/Fetch-Project
```
2. Build the Docker image

```
docker build -t fetch .
```

4. Run the docker the container that you built

```
docker run -p 8080:8080 fetch
```

5. You can now use either Postman or simply use curl from the terminal

```
curl --data '{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}' \
  http://localhost:8080/receipts/process
```

Respone will look something like 
```
{"id":"92fe2e82-7fb3-464c-9529-e42b9f8542ec"}
```
Then check the /getPoints running 
```
curl http://localhost:8080/receipts/<id>/points
```
