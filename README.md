# Trailleo Project

## Overview

This is the backend for trailleo. It was migrated away from a monolith written in Nodejs and hosted on AWS EC2. It was noticed that it was becoming expensive to keep this project running given the number of MAUs were seeing. So for a couple of years it was taken down due to AWS Bills being too high. Untill I decided to rewrite it in Golang (becuase I wanted to learn golang and this was the perfect opportunity)

## Features

- **Trail Discovery**: Users can explore a wide range of trekking paths, from easy walks to challenging climbs.
- **Personalized Recommendations**: Based on user preferences and past activities, Trailleo offers personalized trail suggestions.
- **Community Engagement**: Connect with fellow hikers, share experiences, and join community-led hiking events.

## Technology Stack

- **Database**: MongoDB Atlas is used for storing and retrieving data. I had chosen MongoDB becuase of its generous free tier, text search capabilities via Atlas search indexes & also geo spatial indexing needed to get all data points with a specified region
- **Backend**: Golang is used because it has better performance than nodejs or python (the only languages I knew at this point), since Golang is a compiled language. This makes a differnce in serverless functions becuase it can be cheaper (takes lesser time) to startup and run.
- **Cloud Infra**: I chose Google cloud as cloud functions had a better free tier than AWS lambda functions. (Yes, I know. I am cheap).

## Project Folder Structure

Modules and packages are central to any go project. In this project I have divided each API into its own module, since we make changes and deploy to each API seperately. Each API module has its own folder under the `functions` folder. Common functions like DB clients, common structs and parsers are grouped under one module called `common`. But since env files are required while deploying each function will have its own env var file even though the contents are mostly same across API modules.

## Environment variables

env file is `env.yaml` which is present in every API module folder.

## Deploying

Before deploying make sure
