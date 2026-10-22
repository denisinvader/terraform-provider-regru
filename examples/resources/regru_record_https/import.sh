# Import ID format: zone,name,priority,value.
# Value can safely contain commas.
terraform import regru_record_https.api 'example.com,api,0,alpn=h2,h3 port=443'
