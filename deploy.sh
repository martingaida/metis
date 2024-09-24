#!/bin/bash

# Function to check and install dependencies
check_dependencies() {
    # Check for AWS CLI
    if ! command -v aws &> /dev/null; then
        echo "AWS CLI is not installed. Installing..."
        curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
        unzip awscliv2.zip
        sudo ./aws/install
        rm -rf aws awscliv2.zip
    fi

    # Check for SAM CLI
    if ! command -v sam &> /dev/null; then
        echo "SAM CLI is not installed. Installing..."
        pip install aws-sam-cli
    fi

    # Check for Python
    if ! command -v python3 &> /dev/null; then
        echo "Python 3 is not installed. Please install Python 3 and try again."
        exit 1
    fi

    # Check for Go (for backend)
    if ! command -v go &> /dev/null; then
        echo "Go is not installed. Please install Go and try again."
        exit 1
    fi

    # Check for zip
    if ! command -v zip &> /dev/null; then
        echo "zip is not installed. Installing..."
        sudo apt-get update && sudo apt-get install -y zip
    fi
}

# Function to deploy a component
deploy_component() {
    component=$1
    echo "Deploying $component..."
    cd $component

    # Check for .env file and source it if it exists
    if [ -f .env ]; then
        echo "Found .env file. Loading environment variables..."
        export $(grep -v '^#' .env | xargs)
    fi

    # For backend (Go), we need to build the binary
    if [ "$component" == "backend" ]; then
        # Check if LLM_SERVICE_URL is set
        if [ -z "$LLM_MICROSERVICE_URL" ]; then
            echo "Error: LLM_SERVICE_URL is not set in .env file"
            exit 1
        fi

        # For backend (Go), we need to build the binary
        echo "Building Go Lambda function..."
        cd src
        GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
        cd ..
        
        # Use SAM to package and deploy
        sam package --template-file template.yaml --output-template-file packaged.yaml --s3-bucket metis-backend-w8ej736hj9
        sam deploy --template-file packaged.yaml --stack-name metis-backend --capabilities CAPABILITY_IAM --region us-east-1 --no-confirm-changeset --parameter-overrides LLMServiceUrl=$LLM_MICROSERVICE_URL
    else
        # For other components (like microservices), use standard SAM commands
        sam build
        
        # Check if OPENAI_API_KEY is set and add it to parameter overrides if it exists
        if [ -n "$OPENAI_API_KEY" ]; then
            echo "Using OPENAI_API_KEY from .env file"
            sam deploy --stack-name metis-$component --capabilities CAPABILITY_IAM --region us-east-1 --no-confirm-changeset --parameter-overrides OpenAIApiKey=$OPENAI_API_KEY
        else
            echo "OPENAI_API_KEY not found in .env file. Deployment may fail."
            sam deploy --stack-name metis-$component --capabilities CAPABILITY_IAM --region us-east-1 --no-confirm-changeset
        fi
    fi
    
    cd ..
}

# Check and install dependencies
check_dependencies

# Check command line argument
if [ "$1" = "frontend" ]; then
    echo "Frontend deployment is not implemented in this script."
elif [ "$1" = "backend" ]; then
    deploy_component "backend"
elif [ "$1" = "microservices" ]; then
    deploy_component "microservices"
elif [ "$1" = "all" ]; then
    deploy_component "backend"
    deploy_component "microservices"
    echo "Frontend deployment is not implemented in this script."
else
    echo "Usage: ./deploy.sh [backend|microservices|all]"
    exit 1
fi

echo "Deployment completed successfully!"