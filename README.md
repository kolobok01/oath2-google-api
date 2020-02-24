
# Google OAuth 

## Setting Env

```
vim ~/.bashrc

# Get from GCP Console, after [Add Client ID], you will get it 
# [Add Client ID] look following 

export GOOGLE_CLIENT_ID="xxxxx"
export GOOGLE_CLIENT_SECRET=="xxxxx"
```

## http OAuth

[Add Client ID] from GCP Console

* API & Service -> Credentials -> Cretate credentials -> OAuth Client ID -> Web application

```
go run google_http.go
```

1. http://localhost:8080/

2. http://localhost:8080/login

3. User Login

4. Redirect to http://localhost:8080/callback

5. Get Token

6. https://www.googleapis.com/oauth2/v2/userinfo?access_token=[token.access_token]

7. Get User Info

Refer: https://medium.com/@pliutau/getting-started-with-oauth2-in-go-2c9fae55d187


## Offline OAuth

Add Client ID from GCP Console

* API & Service -> Credentials -> Cretate credentials -> OAuth Client ID -> Other

```
go run google_offline_other_client.go
```

NOTE: redirectURL := **"urn:ietf:wg:oauth:2.0:oob"**

1. Retrive Code: Copy & Paste the AuthURL in browser to retrive the code

2. User Login

3. Paste the code to terminal

4. Get Token

5. https://www.googleapis.com/oauth2/v2/userinfo?access_token=[token.access_token]

6. Get User Info

Refer: https://cloud.google.com/docs/authentication/?hl=zh_TW&_ga=2.126343757.-1956126422.1565750607
