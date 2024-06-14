   # deploy.sh
   #!/bin/bash

   # Set the function name
   FUNCTION_NAME="trek-file-list"

   ENTRY_POINT="GetTrekFileList"

   # Set the runtime
   RUNTIME="go121"

   # Set the trigger type
   TRIGGER="--trigger-http"

   # Set the environment variables file
   ENV_VARS_FILE="env.yaml"
   

   echo "Deploying function $FUNCTION_NAME with runtime $RUNTIME"

   # Deploy the function
   gcloud functions deploy $FUNCTION_NAME \
     --runtime $RUNTIME \
     $TRIGGER \
     --allow-unauthenticated \
     --env-vars-file $ENV_VARS_FILE \
     --entry-point $ENTRY_POINT