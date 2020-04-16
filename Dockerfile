FROM gcr.io/distroless/base:latest

COPY bin/rds-ca-2019-root.pem /bin/rds-ca-2019-root.pem
COPY bin/orders /bin/orders

COPY config/tls/Certificates_PKCS7_v5.6_DoD.der.p7b /config/tls/Certificates_PKCS7_v5.6_DoD.der.p7b
COPY config/tls/dod-sw-ca-54.pem /config/tls/dod-sw-ca-54.pem

COPY swagger/* /swagger/

ENTRYPOINT ["/bin/orders"]

CMD ["serve", "--debug-logging"]

EXPOSE 8080 8443 9443
