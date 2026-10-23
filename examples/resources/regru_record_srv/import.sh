# Import ID format: zone,service,priority,weight,port,target.
terraform import regru_record_srv.sip 'example.com,_sip._udp,0,0,5060,sip-handler'
