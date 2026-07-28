output "vpc_id" {
  value = aws_vpc.this.id
}

output "public_subnet_ids" {
  value = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  value = aws_subnet.private[*].id
}

output "subnet_ids" {
  value = concat(aws_subnet.public[*].id, aws_subnet.private[*].id)
}

output "nat_gateway_ip" {
  value = aws_eip.nat.public_ip
}

output "dns_zone_id" {
  value = aws_route53_zone.this.zone_id
}

output "dhcp_options_id" {
  value = aws_vpc_dhcp_options.this.id
}

output "default_security_group_id" {
  value = aws_security_group.default.id
}

output "app_security_group_id" {
  value = aws_security_group.app.id
}
