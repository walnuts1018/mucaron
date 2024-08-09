resource "aws_s3_bucket" "mucaron" {
  bucket = var.bucket_name
} 

data "aws_iam_policy_document" "mucaron" {
  statement {
principals {
      type        = "AWS"
      identifiers = ["*"]
    }
    actions = [
      "s3:GetObject",
    ]
    resources = [
      aws_s3_bucket.mucaron.arn,
      "${aws_s3_bucket.mucaron.arn}/*",
    ]
  }
}

resource "aws_s3_bucket_policy" "mucaron" {
  bucket = aws_s3_bucket.mucaron.id
  policy = data.aws_iam_policy_document.mucaron
}
