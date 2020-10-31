#!/bin/bash
blogUrl="https://joaonsantos.dev"
dateNow=$(command date)

if res=$(curl -f --url "$blogUrl"); then
	echo "$dateNow - success" >> /var/log/blog-alerts.log
else
	mail -s "[ALERT] Blog is down" -a "From: mail@joaonsantos.dev" joaopns05@gmail.com \
        <<< "Blog is down! joaonsantos.dev is not reachable."
fi

