
openssl genrsa -out spigell.key 4096
openssl req -config ./csr.cnf -new -key spigell.key -nodes -out spigell.csr
