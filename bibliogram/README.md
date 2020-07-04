## Bibliogram

Docker-compose example:

```
bibliogram:
    image: gcr.io/krn-dev/bibliogram:sha-29d25ee
    container_name: bibliogram
    restart: always
    volumes:
      - /home/bibliogram/config.js:/config.js:ro
      - /home/bibliogram/db:/opt/bibliogram/db
    ports:
      - 10407:10407
```
