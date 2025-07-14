# Seo crawler

This is an application where we can pass the URL in the frontnend which we want to analyse and the backend will find out the h1 title description and type of links and save those in database. we also have login and signup options.

Backend: A Go-based web crawler service that analyzes websites for SEO metrics, following best practices.
Frontend: A react based frontend where we can provide a URL that we want to crawl.

### How to run ?

#### Option 1 : docker-compose

- On the root run : `docker-compose up`

#### Option 2 : manually

- For backend:

  - Install dependencies: ` go mod tidy`
  - Build the app `go build -o seo-crawler ./cmd/server/main.go`
  - Find your mysql username password and keep it handy we will need it.
    - Login : `mysql -u root -p` change user and password
    - create a database that we will need in next step
    - `CREATE DATABASE seo_crawler_check;`
  - Run the app: (update values)

    ```
    export DB_HOST=127.0.0.1
    export DB_PORT=3306
    export DB_USER=root
    export DB_PASSWORD=password
    export DB_NAME=seo_crawler_check
    export SERVER_PORT=8080
    export API_KEY=seo-crawler-api-key-2025

    ./seo-crawler
    ```

- For frontend
  - Install dependencies: `npm install`
  - copy `.env.example` and rename it to `.env` (change env in prod)
  - run the app in dev mode : `npm run dev`
  - To deploy on prod :
    - build the app : `yarn build`
    - run it via : `yarn preview` (and deploy this static site as s3 with cloudfront distribution)

### Future improvements

- Add test
- Improve design
- add more features
- Load Test backend and profile it
- Deploy it
