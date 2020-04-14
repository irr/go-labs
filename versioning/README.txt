docker run --rm --net host --name nginx -v $PWD/nginx.conf:/etc/nginx/nginx.conf nginx

http localhost:3000 Accept:application/vnd.t.v1+json
http localhost:3000 Accept:application/vnd.t.v2+json