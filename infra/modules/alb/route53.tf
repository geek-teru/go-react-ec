resource "aws_route53_record" "ec" {
  zone_id = data.aws_route53_zone.public.zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name                   = aws_alb.public_alb.dns_name
    zone_id                = aws_alb.public_alb.zone_id
    evaluate_target_health = true
  }
}
