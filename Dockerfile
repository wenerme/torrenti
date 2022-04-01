FROM wener/base

COPY bin/torrenti /opt/app/bin/torrenti
CMD [ "/opt/app/bin/torrenti", "serve" ]
