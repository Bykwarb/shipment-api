# Shipment api
This API provides endpoints for managing delivery and checking barcode availability.

## Endpoints
- POST /api/v1/shipment
  * Request body:
  ```json
  {
    "barcode": "EXAMPLEBARCODE",
    "sender": "Sender",
    "receiver": "Receiver",
    "is_delivered": false,
    "origin": "Origin",
    "destination": "Destination"
  }
   ```
  * Response body if success:
  ```json
  {
    "message": "shipment successfully created"
  }
  ```
  * Response body if there are some errors
  ```json
  {
    "message": "failed to decode JSON"
  }
  ```
  ```json
  {
    "message": "barcode length must be <= then 25 and >= 13"
  }
  ```
  ```json
  {
    "message": "failed to save shipment"
  }
  ```
  * Response status code:
    - **200 OK**

- DELETE /api/v1/shipment/{id}
  * Path variable: id (int)
  * Response body if success:
  ```json
  {
    "message": "shipment successfully deleted"
  }
  ```
  * Response body if there are some errors
  ```json
  {
    "message": "failed to delete shipment"
  }
  ```
   ```json
  {
    "message": "failed to parse id parameter"
  }
  ```
  * Response status code:
    - **200 OK**
   
- GET /api/v1/shipment/{id}
  * Path variable: id (int)
  * Response body if success:
  ```json
    {
    "Id":585683,
    "Barcode":"J55548763064048LS",
    "Sender":"pZcHxJhlkF",
    "Receiver":"sxylDdTdoL",
    "IsDelivered":false,
    "Origin":"jAXTrp0ir5",
    "Destination":"lCUzDtDtES",
    "CreatedAt":"2023-04-12T15:03:18.745512Z"
    }
  ```
  * Response body if there are some errors
  ```json
  {
    "message": "failed to parse id parameter"
  }
  ```
  ```json
  {
    "message": "failed to get shipment with Id: 585683"
  }
  ```
  * Response status code:
    - **200 OK**

- GET /api/v1/shipment?barcode=barcode
  * Response body:
  ```json
    {
    "Id":585683,
    "Barcode":"J55548763064048LS",
    "Sender":"pZcHxJhlkF",
    "Receiver":"sxylDdTdoL",
    "IsDelivered":false,
    "Origin":"jAXTrp0ir5",
    "Destination":"lCUzDtDtES",
    "CreatedAt":"2023-04-12T15:03:18.745512Z"
    }
  ```
   * Response body if there are some errors
  ```json
  {
    "message": "missing barcode parameter"
  }
  ```
  ```json
  {
    "message": "failed to get shipment with barcode: J55548763064048LS"
  }
  ```
  * Response status code:
    - **200 OK**

- GET /api/v1/barcodes/{barcode}/availability
  * Path variable: barcode (string)
  * Response body:
  ```json
    {
    "unavailable":true
    }
  ```
  * Response status code:
    - **200 OK**
