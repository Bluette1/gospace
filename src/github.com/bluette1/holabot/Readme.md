# holabot
- A twitterbot that responds to tweets.

## Running locally
### Set up
- Clone this  git repository
- Apply for a Twitter developer account
- Create a project if there'sn't an already existing one and create an application called `halobot` for this app
- Retrieve the application's user context keys and add then to a `.env` folder in the local project directory, and a `CREDENTIALS` file according to the given `CREDENTIALS_EXAMPLE` file
- Signup for [ngrok](https://dashboard.ngrok.com/get-started/setup) 
- [Create a twitter sandbox dev environment](https://developer.twitter.com/en/account/environments)
- Start the ngrok tunnel by running
`./ngrok http 9090`
- Copy and paste the `https` link from the ngrok tunnel in the  `.env` file

### Install dependencies
- Run `go install`
### Register and subsribe webhook
- To register and subscribe the webhook run the following
```go install github.com/bluette1/holabot
  holabot -register
  holabot -subscribe
``` 


## Authors

👤 **Marylene Sawyer**
- Github: [@Bluette1](https://github.com/Bluette1)
- Twitter: [@MaryleneSawyer](https://twitter.com/MaryleneSawyer)
- Linkedin: [Marylene Sawyer](https://www.linkedin.com/in/marylene-sawyer-b4ba1295/)


# Acknowledgements

- The content in this repository was retrieved from or inspired by the following sites
  - [How to create a Twitter bot from scratch with Golang](https://kofo.dev/how-to-create-a-twitter-bot-from-scratch-with-golang)
  - [Post, Retrieve and Engage with Tweets](https://developer.twitter.com/en/docs/twitter-api/v1/tweets/post-and-engage/api-reference/post-statuses-update)
  - [Account Activity Api: Premium](https://developer.twitter.com/en/docs/twitter-api/premium/account-activity-api/api-reference/aaa-premium#post-account-activity-all-env-name-subscriptions)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!
