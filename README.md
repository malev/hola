# HOLA

`hola` is an http client that can be configured with files. It's Github friendly and provides flexibility when managing secrets.

## Story

Every time I try to send a cURL request, I have to open a text editor, construct the command, and then copy-paste it into the terminal. This repetitive process has led me to create pre-written cURL command files for frequent use. Inspired by tools like rest.nvim and kulala, I've decided to develop a similar solution without relying on Neovim. While I'm a big fan of Neovim, I aim to create a tool accessible to users who may not be familiar with it yet.

`hola` enables you to define cURL requests in `.http` files, which can then be executed from the terminal. Sensitive information like API keys can be securely stored in a configuration file or injected via environment variables. By encouraging users to organize their requests and secrets in separate files, `hola` facilitates team collaboration. Teams can push their .http files to a version control repository and use pull requests to review and merge changes, streamlining the workflow and ensuring consistency across the team.

## Usage

### Setup a file with your requests

users.http:
```
### Get all users
GET https://api.malev.xyz/api/users

### Create a new user
POST https://api.malev.xyz/api/users
X-API-KEY: {{env|X_API_KEY}}
X-API-SECRET: {{env|X_API_SECRET}}

{
  "username": "malev",
  "email": "malev@lol.com"
}
```

* Export your secrets:
  * `export X_API_KEY=THIS-IS-MY-KEY`
  * `export X_API_SECRET=GIMME-ACCESS`
* `hola users.http --index 0` to send the GET request
* `hola users.http --index 1` to send the POST request (2nd request in the file)

### Make use of configuration files

config.json:
```
{
  "host": "api.malev.xyz"
}
```

users.http
```
### Get all users
GET https://{{host}}/api/users
```

* `hola users.http --config config.json`

