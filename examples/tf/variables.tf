variable "bucket_name" {
    description = "The name of the S3 bucket"
    type        = string
    default = "wplite-test"
}

variable "aws_region" {
    description = "The AWS region"
    type        = string
    default = "us-west-2"
}

variable "cloudfront_aliases" {
    description = "The CloudFront aliases"
    type        = list(string)
    default = []
}

variable "acm_certificate_arn" {
    description = "The ARN of the ACM certificate"
    type        = string
    default = ""
}
