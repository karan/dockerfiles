FROM node:12.18.1-alpine as build
RUN apk --no-cache add git python3 make g++
WORKDIR /app
# Commit from https://github.com/cloudrac3r/bibliogram/commits/master
RUN git clone https://github.com/cloudrac3r/bibliogram.git
RUN cd bibliogram && git reset --hard 5ed035a432d4965f63fdad7ea1e66ed056871304

WORKDIR /app/bibliogram
RUN npm install --no-optional

FROM node:12.18.1-alpine as app
WORKDIR /app/bibliogram
COPY --from=build /app /app
RUN ln -sf /opt/bibliogram/config.js /config.js
EXPOSE 10407
CMD ["npm", "start"]
