# Go Space Voyagers

## Setup

1. Clone this repository.
2. Check go version and then install dependencies:
   ```bash
   go version
   go mod tidy
   ```
3. Add port in a `.env` file.
4. Run the server (You can also use `air` for live reloading):
   ```bash
   go run .
   ```
5. Run unit tests and check code coverage:
   ```bash
   go test ./... -coverprofile=coverage.out
   go tool cover -html=coverage.out
   ```

## API Endpoints

- GET /planets: Retrieves all the planets  
  ![Get Planets](https://github.com/kaitou-1412/Go-Space-Voyagers/blob/main/media/readall.png)
- GET /planets/:id: Retrieves a planet by its ID  
  ![Get Planet By Id](https://github.com/kaitou-1412/Go-Space-Voyagers/blob/main/media/read.png)
- GET /planets/getFuelCost/:id: Retrieves a planet fuel cost by its ID and crew capacity
  ![Get Planet Fuel Cost By Id](https://github.com/kaitou-1412/Go-Space-Voyagers/blob/main/media/fuelcost.png)
- POST /planets: Creates a new planet  
  ![Create Planet](https://github.com/kaitou-1412/Go-Space-Voyagers/blob/main/media/create.png)
- PUT /planets/:id: Updates a planet by its ID  
  ![Update Planet By Id](https://github.com/kaitou-1412/Go-Space-Voyagers/blob/main/media/update.png)
- DELETE /planets/:id: Deletes a planet by its ID  
  ![Delete Planet By Id](https://github.com/kaitou-1412/Go-Space-Voyagers/blob/main/media/delete.png)
