# Import ID format: zone,name,priority,value.
# The value can contain commas.
terraform import regru_record_https.api 'example.com,api,1,alpn=h2,h3 port=443'
