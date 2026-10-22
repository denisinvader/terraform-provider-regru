# Import ID format: zone,name,flag,tag,value.
# Value can safely contain commas.
terraform import regru_record_caa.this 'example.com,@,0,issuewild,ca.example.com; user_id=1234'
