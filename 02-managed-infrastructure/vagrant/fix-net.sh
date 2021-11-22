set -o nounset                              # Treat unset variables as an error

while sleep 1 ; do
  if [[ $(ip r get 8.8.8.8 | grep 10.0.2) ]]; then 
    echo 'Right gw set!'
    break
  else
    echo 'Another try!'
    ip l set down eth1
    sleep 1
    ip l set up eth1
  fi
done

