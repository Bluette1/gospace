## holabot
- A twitter bot that responds to tweets.

### Running locally
#### Set up
- Clone this  git repository
- Apply for a Twitter developer account
- Create a project if there'sn't an already existing one and create an application called `halobot` for this app
- Retrieve the application's user context keys and add then to a `.env` folder in the local project directory, and a `CREDENTIALS` file according to the given `CREDENTIALS_EXAMPLE` file
- Signup for [ngrok](https://dashboard.ngrok.com/get-started/setup) 
- [Create a twitter sandbox dev environment](https://developer.twitter.com/en/account/environments)
- Start the ngrok tunnel by running
`./ngrok http 9090`
- Copy and paste the `https` link from the ngrok tunnel in the  `.env` file

#### Install dependencies
- Run `go install github.com/bluette1/holabot` inside the root 
#### Register and subsribe webhook
- To register and subscribe the webhook run the following
```go install github.com/bluette1/holabot
  holabot -register
  holabot -subscribe
``` 

#### Twitter signin

-To grant permission from a Twitter user, you have to implement [Twitter Signin](https://developer.twitter.com/en/docs/authentication/guides/log-in-with-twitter) 
- Edit app settings to enable [3-legged OAuth](https://developer.twitter.com/en/docs/authentication/oauth-1-0a/obtaining-user-access-tokens)
   - Make sure you register the callback url you'll be using in your implementation of twitter signin.
- Start the bot by running `holabot`
- Hit the root http://localhost:9090/signin
   - Follow the instructions to sign into Twitter for a successful login process
#### Authors

👤 **Marylene Sawyer**
- Github: [@Bluette1](https://github.com/Bluette1)
- Twitter: [@MaryleneSawyer](https://twitter.com/MaryleneSawyer)
- Linkedin: [Marylene Sawyer](https://www.linkedin.com/in/marylene-sawyer-b4ba1295/)


#### Acknowledgements

- The content in this repository was retrieved from or inspired by the following sites
  - [How to create a Twitter bot from scratch with Golang](https://kofo.dev/how-to-create-a-twitter-bot-from-scratch-with-golang)
  - [Post, Retrieve and Engage with Tweets](https://developer.twitter.com/en/docs/twitter-api/v1/tweets/post-and-engage/api-reference/post-statuses-update)
  - [Account Activity Api: Premium](https://developer.twitter.com/en/docs/twitter-api/premium/account-activity-api/api-reference/aaa-premium#post-account-activity-all-env-name-subscriptions)
  
[Twittergo Examples](https://github.com/kurrik/twittergo-examples/blob/master/sign_in/main.go)
[Outh1a](https://golangrepo.com/repo/kurrik-oauth1a-go-authentication-oauth)
[Twitter Signin](https://developer.twitter.com/en/docs/authentication/guides/log-in-with-twitter)
#### 🤝 Contributing

Contributions, issues and feature requests are welcome!
