# DailyNewsPodcast

DailyNewsPodcast is a web application that automatically generates daily podcast episodes discussing top news headlines. It leverages a Go backend, AWS Lambda, Meta Llama for conversation generation, and AWS Polly for text-to-speech, delivering fresh news podcasts every day.

## Features

- **Automated News Podcast Generation:**  
  Fetches top headlines from [NewsAPI](https://newsapi.org/), uses Meta Llama to create podcast-style discussions, and converts them to audio with AWS Polly.

- **Cloud Storage:**  
  Generated podcast episodes are stored in an AWS S3 bucket and served to users. User data is stored in an AWS RDS PostgreSQL database.

- **Go Backend with GraphQL:**  
  Handles API requests, user authentication (JWT-based), and podcast retrieval.

- **User Login & Signup:**  
  Secure authentication and session management via JWT tokens.

- **Frontend with React:**  
  Modern, responsive web UI for browsing and listening to podcasts.

- **Customizable News Experience:**  
  Users can select their preferred country and topics for personalized news podcasts (in progress).

## Tech Stack

- **Backend:** Go, GraphQL
- **Frontend:** React, JavaScript, HTML, CSS
- **Database** PostgreSQL
- **Cloud Services:** AWS S3, AWS Polly, AWS Lambda, AWS RDS
- **AI:** Meta Llama (for conversation generation)
- **Authentication:** JWT tokens

## Getting Started

1. **Clone the repository**
    ```sh
    git clone https://github.com/ShavaizKhan/DailyNewsPodcast.git
    ```

2. **Backend Setup**
    - Install Go dependencies.
    - Set up AWS credentials and environment variables for Lambda, S3, and Polly.
    - Configure NewsAPI and Meta Llama access.

3. **Frontend Setup**
    - Install Node.js dependencies:
      ```sh
      cd frontend
      npm install
      npm start
      ```

4. **Run Locally**
    - Start the Go backend server.
    - Launch the React frontend.

## Usage

- Sign up or log in.
- Browse and listen to daily news podcasts.
- (Soon) Choose your country and preferred topics for a tailored podcast experience.
