# Trailleo backend

## Overview

This is the backend for trailleo. It was migrated away from a monolith written in Nodejs and hosted on AWS EC2. We noticed that it was becoming expensive to keep this project running given the number of MAUs we were seeing. So for a couple of years it was taken down due to AWS Bills being too high. Untill I decided to rewrite it in Golang (becuase I wanted to learn golang and this was the perfect opportunity)

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

env file is `env.yaml` which is present in every API module folder. Add the following to env.yaml:

```
DB_URI: <your mongo db uri>
```

## Deploying

Before deploying make sure you have the glocud CLI installed and configured to the required project. Also make sure google cloud functions are enabled.
Do deploy a function

- navigate to the folder `./functions/<function folder>`.
- make the `deploy.sh` bash script executable by running `chmod +x deploy.sh` from the root of the function module folder
- run `./deploy.sh`

## Running & debugging locally

- create a script file called `run-local.sh` with the following content

```
DB_URI="<your mongo db uri>"
FUNCTION_NAME="<the function name as specified inside the init() function inside function.go pkg>"
DB_URI=$DB_URI FUNCTION_TARGET=$FUNCTION_NAME LOCAL_ONLY=true go run cmd/main.go
```

- make the `run-local.sh` bash script executable by running `chmod +x run-local.sh` from the root of the function module folder
- run `run-local.sh`

## Updating CORS on Google Storage Bucket

To update CORS on a Google Storage Bucket, follow these steps:

1. **Create a CORS configuration file**: Create a JSON file named `cors.json` with the following content:

```
[
	{
		"maxAgeSeconds": 3600,
		"method": ["GET", "HEAD", "OPTIONS"],
		"origin": ["*"],
		"responseHeader": [
			"Content-Type",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods",
			"Access-Control-Allow-Headers"
		]
	}
]

```

This configuration allows GET, HEAD, and OPTIONS requests from any origin, with a maximum age of 3600 seconds.

2. **Upload the CORS configuration file to your bucket**: Use the `gsutil` command-line tool to upload the `cors.json` file to your Google Storage Bucket.

**Using `glcoud` & `gsutils`:**

- Run the following command in your terminal:

```
gcloud storage buckets update gs://<your-bucket-name>  --cors-file=<path_to_file>
```

Replace `<your-bucket-name>` & `path_to_file` with the name of your Google Storage Bucket and cors file path.

3. **Verify the CORS configuration**: After uploading the CORS configuration file, verify that it has been successfully applied to your bucket. You can do this by checking the bucket's CORS configuration in the Google Cloud Console or by using the `gsutil` command-line tool.

- Run the following command in your terminal:

```
gsutil cors get gs://<your-bucket-name>
```

This should display the CORS configuration for your bucket.
