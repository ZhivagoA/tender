# Как запустить
Сбилдить docker-compose.yaml
Сервис доступен на 8080

# API Documentation

В этом файле расписаны пошаговые запросы для ручек вместе с тестовым телом запросов
Важно учесть что я использую UUID и id генерируются автоматически, поэтому при проверке обращайте внимание на внешние ключи которые передается при запросе.

## 1. Health Check

### `GET /ping`
A simple health check endpoint to verify if the server is running.

**Example:**
curl -X GET http://localhost:8080/ping

---

## 2. Organization Endpoints

### `POST /organizations/new`
Creates a new organization.

**Request Body:**
{
  "name": "Tech Corp",
  "type": "LLC"
}

**Example:**
curl -X POST http://localhost:8080/organizations/new \
-H "Content-Type: application/json" \
-d '{"name": "Tech Corp", "type": "LLC"}'

---

### `DELETE /organizations/{organization_id}`
Deletes an organization by its ID.

**Example:**
curl -X DELETE http://localhost:8080/organizations/3330a0f8-3c22-4c87-9e64-461543bca758

---

### `POST /organizations/responsibles/new`
Assigns a user as responsible for an organization.

**Request Body:**
{
  "organization_id": "3330a0f8-3c22-4c87-9e64-461543bca758",
  "user_id": "29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1"
}

**Example:**
curl -X POST http://localhost:8080/organizations/responsibles/new \
-H "Content-Type: application/json" \
-d '{"organization_id": "3330a0f8-3c22-4c87-9e64-461543bca758", "user_id": "29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1"}'

---

### `DELETE /organizations/responsibles/{responsible_id}`
Removes a responsible user from an organization.

**Example:**
curl -X DELETE http://localhost:8080/organizations/responsibles/d4f77f0b-f102-4480-8225-8eee848d79d8

---

## 3. User Endpoints

### `POST /users/new`
Creates a new user.

**Request Body:**
{
  "username": "john_doe",
  "first_name": "John",
  "last_name": "Doe"
}

**Example:**
curl -X POST http://localhost:8080/users/new \
-H "Content-Type: application/json" \
-d '{"username": "john_doe", "first_name": "John", "last_name": "Doe"}'

---

### `DELETE /users/{user_id}`
Deletes a user by its ID.

**Example:**
curl -X DELETE http://localhost:8080/users/29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1

---

## 4. Tender Endpoints

### `POST /tenders/new`
Creates a new tender.

**Request Body:**
{
  "title": "New Tender",
  "description": "Description of the new tender",
  "organization_id": "3330a0f8-3c22-4c87-9e64-461543bca758",
  "responsible_user_id": "d4f77f0b-f102-4480-8225-8eee848d79d8"
}

**Example:**
curl -X POST http://localhost:8080/tenders/new \
-H "Content-Type: application/json" \
-d '{
  "title": "New Tender",
  "description": "Description of the new tender",
  "organization_id": "3330a0f8-3c22-4c87-9e64-461543bca758",
  "responsible_user_id": "d4f77f0b-f102-4480-8225-8eee848d79d8"
}'

---

### `GET /tenders/my`
Retrieves tenders where the authenticated user is responsible.

**Example:**
curl -X GET http://localhost:8080/tenders/my \
-H "X-User-ID: d4f77f0b-f102-4480-8225-8eee848d79d8"

---

### `GET /tenders/status/{tender_id}`
Retrieves the status of a specific tender.

**Example:**
curl -X GET http://localhost:8080/tenders/status/6a1a93c2-21e3-43cd-9256-1b94111a7cb6

---

### `PUT /tenders/edit`
Edits an existing tender.

**Request Body:**
{
  "id": "6a1a93c2-21e3-43cd-9256-1b94111a7cb6",
  "title": "Updated Tender Title",
  "description": "Updated description",
  "organization_id": "3330a0f8-3c22-4c87-9e64-461543bca758",
  "responsible_user_id": "d4f77f0b-f102-4480-8225-8eee848d79d8",
  "status": "CREATED",
  "version": 2
}

**Example:**
curl -X PUT "http://localhost:8080/tenders/edit" \
-H "Content-Type: application/json" \
-d '{
  "id": "6a1a93c2-21e3-43cd-9256-1b94111a7cb6",
  "title": "Updated Tender Title",
  "description": "Updated description",
  "organization_id": "3330a0f8-3c22-4c87-9e64-461543bca758",
  "responsible_user_id": "d4f77f0b-f102-4480-8225-8eee848d79d8",
  "status": "CREATED",
  "version": 2
}'

---

### `POST /tenders/rollback`
Rolls back to a previous version of a tender.

**Request Body:**
{
  "tender_id": "6a1a93c2-21e3-43cd-9256-1b94111a7cb6",
  "version": 2,
  "responsible_user_id": "d4f77f0b-f102-4480-8225-8eee848d79d8"
}

**Example:**
curl -X POST "http://localhost:8080/tenders/rollback" \
-H "Content-Type: application/json" \
-d '{
  "tender_id": "6a1a93c2-21e3-43cd-9256-1b94111a7cb6",
  "version": 2,
  "responsible_user_id": "d4f77f0b-f102-4480-8225-8eee848d79d8"
}'

---

## 5. Bid Endpoints

### `POST /bids/new`
Creates a new bid for a tender.

**Request Body:**
{
  "amount": 1000.50,
  "tender_id": "6a1a93c2-21e3-43cd-9256-1b94111a7cb6",
  "user_id": "29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1",
  "status": "CREATED",
  "version": 1
}

**Example:**
curl -X POST http://localhost:8080/bids/new \
-H "Content-Type: application/json" \
-d '{
  "amount": 1000.50,
  "tender_id": "6a1a93c2-21e3-43cd-9256-1b94111a7cb6",
  "user_id": "29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1",
  "status": "CREATED",
  "version": 1
}'

---

### `POST /bids/submit_decision`
Submits a decision (approve/reject) on a bid.

**Request Body:**
{
  "bid_id": "a3be1115-7576-4eb4-bfb9-2eddebe17d2d",
  "tender_id": "6a1a93c2-21e3-43cd-9256-1b94111a7cb6",
  "responsible_user_id": "29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1",
  "decision": "approve"
}

**Example:**
curl -X POST http://localhost:8080/bids/submit_decision \
-H "Content-Type: application/json" \
-d '{
  "bid_id": "a3be1115-7576-4eb4-bfb9-2eddebe17d2d",
  "tender_id": "6a1a93c2-21e3-43cd-9256-1b94111a7cb6",
  "responsible_user_id": "29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1",
  "decision": "approve"
}'

---

### `GET /bids/list`
Retrieves a list of all bids.

**Example:**
curl -X GET http://localhost:8080/bids/list

---

### `GET /bids/my`
Retrieves the bids made by the authenticated user.

**Example:**
curl -X GET http://localhost:8080/bids/my \
-H "X-User-ID: 29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1"

---

### `GET /bids/status/{bid_id}`
Retrieves the status of a specific bid.

**Example:**
curl -X GET http://localhost:8080/bids/status/a3be1115-7576-4eb4-bfb9-2eddebe17d2d

---

### `PUT /bids/edit/{bid_id}`
Edits an existing bid.

**Request Body:**
{
  "amount": 1200.75,
  "status": "PUBLISHED",
  "user_id": "29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1"
}

**Example:**
curl -X PUT http://localhost:8080/bids/edit/a3be1115-7576-4eb4-bfb9-2eddebe17d2d \
-H "Content-Type: application/json" \
-d '{
  "amount": 1200.75,
  "status": "PUBLISHED",
  "user_id": "29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1"
}'

---

### `POST /bids/rollback`
Rolls back to a previous version of a bid.

**Request Body:**
{
  "bid_id": "a3be1115-7576-4eb4-bfb9-2eddebe17d2d",
  "version": 1
}

**Example:**
curl -X POST http://localhost:8080/bids/rollback \
-H "Content-Type: application/json" \
-d '{
  "bid_id": "a3be1115-7576-4eb4-bfb9-2eddebe17d2d",
  "version": 1
}'

---

### `GET /bids/reviews/{bid_id}`
Retrieves the reviews of a specific bid.

**Example:**
curl -X GET http://localhost:8080/bids/reviews/a3be1115-7576-4eb4-bfb9-2eddebe17d2d

---

### `POST /bids/feedback`
Submits feedback for a bid.

**Request Body:**
{
  "bid_id": "a3be1115-7576-4eb4-bfb9-2eddebe17d2d",
  "user_id": "29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1",
  "comment": "Great bid!",
  "rating": 5
}

**Example:**
curl -X POST http://localhost:8080/bids/feedback \
-H "Content-Type: application/json" \
-d '{
  "bid_id": "a3be1115-7576-4eb4-bfb9-2eddebe17d2d",
  "user_id": "29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1",
  "comment": "Great bid!",
  "rating": 5
}'

---

### `DELETE /organizations/{organization_id}`
Deletes an organization by its ID.

**Example:**
curl -X DELETE http://localhost:8080/organizations/3330a0f8-3c22-4c87-9e64-461543bca758

---

### `DELETE /users/{user_id}`
Deletes a user by its ID.

**Example:**
curl -X DELETE http://localhost:8080/users/29b433fb-b0ac-4f00-a5ad-7b5b8d3e7ff1


