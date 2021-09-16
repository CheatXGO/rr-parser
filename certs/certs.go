package certs

const ( //openssl req -newkey rsa:2048 -sha256 -nodes -keyout cert.key -x509 -days 365 -out cert.pem -subj "/C=US/ST=New York/L=Brooklyn/O=Example Brooklyn Company/CN=YOURDOMAIN OR WHITE IP HERE"
	CERT = "C:/GitHub/rr-parser/certs/cert.pem"
	KEY  = "C:/GitHub/rr-parser/certs/cert.key"
)
