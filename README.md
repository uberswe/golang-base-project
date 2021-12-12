# Golang Base Project

A minimal Golang project with user authentication ready out of the box. All frontend assets should be less than 100 kB on every page load. 

See a live example at: [https://www.golangbase.com](https://www.golangbase.com)

Functionality includes:

 - Login
 - Logout
 - Register
 - User Activation
 - Resend Activation Email
 - Forgot Password
 - Admin Dashboard
 - Search
 - Throttling

This easiest way for me to achieve this was with a database. I decided to use [GORM](https://gorm.io/docs/) which should fully support MySQL, PostgreSQL, SQLite, SQL Server and Clickhouse or any other databases compatible with these dialects.

The frontend is based off of examples from [https://getbootstrap.com/docs/5.0/examples/](https://getbootstrap.com/docs/5.0/examples/).

## Getting started

Simply run `go run cmd/base/main.go` and the entire project should run using an sqlite in-memory database.

## Project structure

This is the latest way I like to organize my projects. It's something that is always evolving and I know some will like this structure while others may not and that is ok. 

I have mixed in the frontend assets with the go project. The `/src` folder has js and css assets for the frontend which are compiled and put into `/dist` where all the html templates also are.

I create packages when I feel there is a need for them to be shared at different levels of the project such as in another package or the routes, my rules here are very flexible and I have yet to come up with a good rule. Packages are in folders like `/middleware`, `email`, `util` and `config`.

You can run this project with a single go file by typing `go run cmd/base/main.go`.

There is a `/models` package which contains all the database models.

The `/routes` package contains all the route functions and logic. Typically, I try to break of the logic into other packages when functions become too big but I have no strict rule here.

All in all I have tried to keep the project simple and easy to understand. I want this project to serve as a template for myself and perhaps others when you want to create a new website.

## Docker

A dockerfile and docker compose file is provided to make running this project easy. Simply run `docker-compose up`.

You will need to change the env variables for sending email, when testing locally I recommend [Mailtrap.io](https://mailtrap.io/).

## Dependencies

I have tried to keep the dependencies low, there is always a balance here in my opinion and I have included the golang vendor folder and compiled assets so that there is no need to download anything to use this project other than the project itself.

### Go Dependencies

The following dependencies are used with Go.

 - [Gin](https://github.com/gin-gonic/gin) - A web framework which makes routes, middleware and static assets easier to use and handle.
 - [GORM](https://gorm.io/index.html) - An ORM library to make writing queries easier.

### NPM Dependencies

I personally dislike NPM because I there seems to be so many dependencies in javascript projects often with vulnerabilites which are hard to fix. However, npm is also one of the easiest ways to build an optimized frontend today which can be modified and built upon by others. I have tried to keep dependencies low and most are only used for the compiling of the static assets.

 - [Bootstrap 5](https://getbootstrap.com/docs/5.0/getting-started/introduction/) - Bootstrap 5 is a frontend framework that makes it easy to create a good looking website which is responsive.
 - [Webpack](https://webpack.js.org/) - An easy way to bundle and compile assets.
 - [sass-loader](https://www.npmjs.com/package/sass-loader) - To compile scss files from bootstrap together with custom styles added by me.
 - [PurgeCss](https://purgecss.com/) - PurgeCss removes unused CSS so that we only load what is needed. Since we have Bootstrap our compiled css would be around 150kB without PurgeCSS compared to around 10kB as it is now.

There are some more dependencies but I have focused on mentioning those that have the greatest impact on the project as a whole.

## GitHub Actions

There is a workflow to deploy to my personal server whenever there is a merge to master. This is the way I like to deploy for personal projects. The steps are:

 - Make a new image with docker
 - Push this image to a private container registry on my server, you can see the registry here https://registry.beubo.com/repositories/20
 - Then I use docker-compose to pull the latest image from the private registry

I use [supervisor](http://supervisord.org/) with [docker-compose](https://docs.docker.com/compose/production/) to run my containers. [Caddy](https://caddyserver.com/) handles the SSL configuration and routing. I use [Ansible](https://docs.ansible.com/ansible/latest/user_guide/playbooks.html) to manage my configurations.

## Contributions

Contributions are welcome and greatly appreciated. Please note that I am not looking to add any more features to this project but I am happy to take care of bugfixes, updates and other suggestions. If you have a question or suggestion please feel free to [open an issue](https://github.com/uberswe/golang-base-project/issues/new). To contribute code, please fork this repository, make your changes on a separate branch and then [open a pull request](https://github.com/uberswe/golang-base-project/compare).

For security related issues please see my profile, [@uberswe](https://github.com/uberswe), for ways of contacting me privately. 

## License

Please see the `LICENSE` file in the project repository.