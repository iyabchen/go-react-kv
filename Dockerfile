FROM golang:1.10-alpine AS go-builder
RUN apk --no-cache add git bzr mercurial
ENV D=/go/src/github.com/iyabchen/go-react-kv/server
# deps - using the closest thing to an official dependency tool: https://github.com/golang/dep
RUN go get -u github.com/golang/dep/...
COPY server/ $D/
RUN cd $D/ && dep ensure -v 
RUN cd $D/cmd && go build -o /tmp/srv

FROM node:9.6.1 as node-builder
RUN mkdir /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH
COPY client/package.json /usr/src/app/package.json
RUN npm install --silent
RUN npm install react-scripts@1.1.1 -g --silent
COPY client/ /usr/src/app/
RUN npm run build --production

# final stage
FROM alpine:3.8
WORKDIR /app
COPY --from=go-builder /tmp/srv /app/
COPY --from=node-builder /usr/src/app/build /app/client

# have to use CMD due to heroku limitation 
# heroku runs with non-privilege user and has no write permission on local storage
CMD ["/app/srv" ]