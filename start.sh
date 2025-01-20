#!/bin/bash

chmod +x ./docker-utils/bash/init-network.sh
./docker-utils/bash/init-network.sh

chmod +x ./docker-utils/bash/add-client.sh
./docker-utils/bash/add-client.sh -n frontend1 -i 10.0.11.3 -p 3001:3000 -e NEXT_PUBLIC_API_URL=http://10.0.10.2:8080


# # # ./start_services.sh <mongo_container_name> <mongo_ip> <server_container_name> <server_ip>
# chmod +x ./docker-utils/bash/add-server.sh
# ./docker-utils/bash/add-server.sh -m mongo-server-1 -i 10.0.10.3 -s scrapper-server-1 -j 10.0.10.2 -p 8080:8080