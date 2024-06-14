   # deploy.sh
   #!/bin/bash

   # Set the function name
   FUNCTION_NAME="trek-detail"

   ENTRY_POINT="GetTrekDetail"

   # Set the runtime
   RUNTIME="go121"

   # Set the trigger type
   TRIGGER="--trigger-http"

   # Set the environment variables file
   ENV_VARS_FILE="env.yaml"

    # Run go mod tidy and go mod vendor
    echo "Tidying and vendoring Go modules"
    go mod tidy
    go mod vendor
   

   echo "Deploying function $FUNCTION_NAME with runtime $RUNTIME"

   # Deploy the function
   gcloud functions deploy $FUNCTION_NAME \
     --runtime $RUNTIME \
     $TRIGGER \
     --allow-unauthenticated \
     --env-vars-file $ENV_VARS_FILE \
     --entry-point $ENTRY_POINT