FROM golang:1.5-onbuild

COPY . /usr/src/app
CMD ["go-wrapper", "run"]
