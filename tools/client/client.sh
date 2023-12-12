#!/bin/bash

# Com as ocnfigurações atuais, o client faz uma requisição por segundo,
# Para aumentar o número de requisições por segundo, descomente as linhas 
# 14 e 32 e comente a linha 33, você terá que remover tanto a imagem como
# o container client-forum-api e subir a stack novamente para o rebuild.

#HOST='catalogo-api:8080'

HOST='localhost:8080'

while true
    do
	NUMB=`expr $RANDOM % 100 + 1`
	#TEMP=`expr 1 + $(awk -v seed="$RANDOM" 'BEGIN { srand(seed); printf("%.4f\n", rand()) }')`
        
	if [ $NUMB -le 25 ]; then
	    curl --silent --output /dev/null http://${HOST}/v1/addresses
	    curl --silent --output /dev/null http://${HOST}/addresses/2e61ddac-c3cc-46e9-ba88-0e86a790c924/cinemas
		elif [ $NUMB -ge 26 ] && [ $NUMB -le 35 ] ; then
			curl --silent --output /dev/null http://${HOST}/v1/cinemas/292cb98c-62ab-49ef-8e23-dc793a86061d
		elif [ $NUMB -ge 36 ] && [ $NUMB -le 45 ] ; then
			curl --silent --output /dev/null http://${HOST}/v1/cinemas/292cb98c-62ab-49ef-8e23-dc793a86061
		elif [ $NUMB -ge 46 ] && [ $NUMB -le 55 ] ; then
			curl --silent --output /dev/null http://${HOST}/v1/movies
		elif [ $NUMB -ge 56 ] && [ $NUMB -le 65 ] ; then
			curl --silent --output /dev/null http://${HOST}/v1/movies/206dad85-cbcd-4b71-8fda-efd6ca87ebc7
		elif [ $NUMB -ge 66 ] && [ $NUMB -le 75 ] ; then
			curl --silent --output /dev/null http://${HOST}/v1/movies/xxxx
		elif [ $NUMB -ge 76 ] && [ $NUMB -le 85 ] ; then
			curl --silent --output /dev/null http://${HOST}/v1/movies/206dad85-cbcd-4b71-8fda-efd6ca87ebc7
        elif [ $NUMB -ge 86 ] && [ $NUMB -le 95 ] ; then
			curl --silent --output /dev/null http://${HOST}/v1/movies/206dad85-cbcd-4b71-8fda-efd6ca87ebc7/posters/2175d4e2-4d9c-411d-a986-08dc8f4e6a51
        elif [ $NUMB -ge 96 ] && [ $NUMB -le 98 ] ; then
			curl --silent --output /dev/null http://${HOST}/v1/movies/206dad85-cbcd-4b71-8fda-efd6ca87ebc7/posters/2175d4e2-4d9c-411d-a986-08dc8f4e6a51
	else
	    curl --silent --output /dev/null http://${HOST}/v1/cinemas/292cb98c-62ab-49ef-8e23-dc793a86061f
        fi

	#sleep $TEMP
	sleep 0.75
done

