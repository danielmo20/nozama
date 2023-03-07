terraform {
    required_providers {
        aws = {
            source = "hashicorp/aws"
            version = "~> 4.16"
        }
    }
}

provider "aws" {
    profile = "default"
    region = "us-east-2"
}

//nozama DynamoDB definitions for Orders
resource "aws_dynamodb_table" "orders_table" {
    name = "orders"
    hash_key = "order_id"
    billing_mode = "PAY_PER_REQUEST"
    attribute {
        name = "order_id"
        type = "S"
    }
}

//nozama DynamoDB definitions for Orders
resource "aws_dynamodb_table" "payments_table" {
    name = "payments"
    hash_key = "payment_id"
    billing_mode = "PAY_PER_REQUEST"
    attribute {
        name = "payment_id"
        type = "S"
    }
    attribute {
        name = "order_id"
        type = "S"
    }
    global_secondary_index {
        name = "order_id-index"
        hash_key = "order_id"    
        projection_type = "ALL"
    }
}

//Nozama API GATEWAY definition
resource "aws_api_gateway_rest_api" "nozama-api" {
    name = "nozama-api"
    description = "The Nozama AWS API Gateway"
    endpoint_configuration {
        types = ["REGIONAL"]
    }
}
// Resource - POST orders/
resource "aws_api_gateway_resource" "orders" {
    rest_api_id = aws_api_gateway_rest_api.nozama-api.id
    parent_id = aws_api_gateway_rest_api.nozama-api.root_resource_id
    path_part = "orders"
}
resource "aws_api_gateway_method" "createOrder" {
    rest_api_id = aws_api_gateway_rest_api.nozama-api.id
    resource_id = aws_api_gateway_resource.orders.id
    http_method = "POST"
    authorization = "NONE"
}

// Resource - PATCH payments/
resource "aws_api_gateway_resource" "payments" {
    rest_api_id = aws_api_gateway_rest_api.nozama-api.id
    parent_id   = aws_api_gateway_rest_api.nozama-api.root_resource_id
    path_part   = "payments"
}
resource "aws_api_gateway_method" "updatePayments" {
    rest_api_id = aws_api_gateway_rest_api.nozama-api.id
    resource_id = aws_api_gateway_resource.payments.id
    http_method = "PATCH"
    authorization = "NONE"
}

//Role definition for all resources
 resource "aws_iam_role" "nozama-role" {
    name = "nozama-role"
    assume_role_policy = <<EOF
    {
    "Version": "2012-10-17",
    "Statement": [
        {
        "Action": "sts:AssumeRole",
        "Principal": {
            "Service": [
                "apigateway.amazonaws.com",
                "lambda.amazonaws.com"
            ]
        },
        "Effect": "Allow",
        "Sid": ""
        }
    ]
    }
    EOF
}
data "template_file" "nozama-policy-file" {
    template = "${file("${path.module}/policy.json")}"
}
resource "aws_iam_policy" "nozama-policy" {
    name = "nozama-policy"
    path = "/"
    description = "IAM policy for Nozama lambda functions"
    policy = data.template_file.nozama-policy-file.rendered
}
resource "aws_iam_role_policy_attachment" "nozama-role" {
    role = aws_iam_role.nozama-role.name
    policy_arn = aws_iam_policy.nozama-policy.arn
} 

//Lambda - Create Orders (http-post) -> (sqs-payments-producer)
resource "aws_lambda_function" "CreateOrdersLambda" {
    function_name = "CreateOrdersLambda"
    filename = "../bin/create-orders/main.zip"
    handler = "main"
    runtime = "go1.x"
    source_code_hash = filebase64sha256("../bin/create-orders/main.zip")
    role = aws_iam_role.nozama-role.arn
    timeout = "5"
    memory_size = "128"
}
resource "aws_api_gateway_integration" "create-order-lambda" {
    depends_on = [
      aws_lambda_function.CreateOrdersLambda
    ]
    rest_api_id = aws_api_gateway_rest_api.nozama-api.id
    resource_id = aws_api_gateway_method.createOrder.resource_id
    http_method = aws_api_gateway_method.createOrder.http_method
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = aws_lambda_function.CreateOrdersLambda.invoke_arn
}
resource "aws_lambda_permission" "apigw-createOrder" {
    depends_on = [
      aws_lambda_function.CreateOrdersLambda
    ]
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.CreateOrdersLambda.function_name
    principal = "apigateway.amazonaws.com"
    source_arn = "${aws_api_gateway_rest_api.nozama-api.execution_arn}/*/POST/orders"
}

//Lambda - Update Payments (http-patch) -> (sqs-orders-producer)
resource "aws_lambda_function" "UpdatePaymentsLambda" {
    function_name = "UpdatePaymentsLambda"
    filename = "../bin/update-payments/main.zip"
    handler = "main"
    runtime = "go1.x"
    source_code_hash = filebase64sha256("../bin/update-payments/main.zip")
    role = aws_iam_role.nozama-role.arn
    timeout = "5"
    memory_size = "128"
}
resource "aws_api_gateway_integration" "update-payment-lambda" {
    depends_on = [
      aws_lambda_function.UpdatePaymentsLambda
    ]
    rest_api_id = aws_api_gateway_rest_api.nozama-api.id
    resource_id = aws_api_gateway_method.updatePayments.resource_id
    http_method = aws_api_gateway_method.updatePayments.http_method
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = aws_lambda_function.UpdatePaymentsLambda.invoke_arn
}
resource "aws_lambda_permission" "apigw-updatePayment" {
    depends_on = [
      aws_lambda_function.UpdatePaymentsLambda
    ]
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.UpdatePaymentsLambda.function_name
    principal = "apigateway.amazonaws.com"
    source_arn = "${aws_api_gateway_rest_api.nozama-api.execution_arn}/*/PATCH/payments"
}

resource "aws_cloudwatch_log_group" "nozama-api-cw-group" {
  name = "API-Gateway-Execution-Logs_${aws_api_gateway_rest_api.nozama-api.id}/dev"
 
}

resource "aws_api_gateway_account" "api_gateway_account" {
  cloudwatch_role_arn = aws_iam_role.nozama-role.arn
}

//Deploy nozama-api
resource "aws_api_gateway_deployment" "nozama-api-stage-dev" {
  depends_on = [
    aws_api_gateway_integration.create-order-lambda,
    aws_api_gateway_integration.update-payment-lambda,
    aws_cloudwatch_log_group.nozama-api-cw-group
  ]
  rest_api_id = aws_api_gateway_rest_api.nozama-api.id
  stage_name  = "dev"
}

//Lambda - Create Payments (sqs-payments-consumer)
resource "aws_lambda_function" "CreatePaymentsLambda" {
    function_name = "CreatePaymentsLambda"
    filename = "../bin/create-payments/main.zip"
    handler = "main"
    runtime = "go1.x"
    source_code_hash = filebase64sha256("../bin/create-payments/main.zip")
    role = aws_iam_role.nozama-role.arn
    timeout = "5"
    memory_size = "128"
}

//Lambda - Update Orders (sqs-orders-consumer)
resource "aws_lambda_function" "UpdateOrdersLambda" {
    function_name = "UpdateOrdersLambda"
    filename = "../bin/update-orders/main.zip"
    handler = "main"
    runtime = "go1.x"
    source_code_hash = filebase64sha256("../bin/update-orders/main.zip")
    role = aws_iam_role.nozama-role.arn
    timeout = "5"
    memory_size = "128"
}

//SQSPayments for OnOrderCreated events 
resource "aws_sqs_queue" "SQSPayments" {
    name = "SQSPayments"
    max_message_size = 5120
    message_retention_seconds = 86400
    receive_wait_time_seconds = 10
}

//SQSOrders for OnPaymentProcessed events
resource "aws_sqs_queue" "SQSOrders" {
    name = "SQSOrders"
    max_message_size = 5120
    message_retention_seconds = 86400
    receive_wait_time_seconds = 10
}

resource "aws_lambda_event_source_mapping" "sqs-payments-lambda" {
    event_source_arn = aws_sqs_queue.SQSPayments.arn
    enabled = true
    function_name = aws_lambda_function.CreatePaymentsLambda.arn
    batch_size = 1
    depends_on = [
        aws_lambda_function.CreatePaymentsLambda
    ]
}

resource "aws_lambda_event_source_mapping" "sqs-orders-lambda" {
    event_source_arn = aws_sqs_queue.SQSOrders.arn
    enabled = true
    function_name = aws_lambda_function.UpdateOrdersLambda.arn
    batch_size = 1
    depends_on = [
        aws_lambda_function.UpdateOrdersLambda
    ]
}