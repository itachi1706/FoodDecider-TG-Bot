FROM alpine:latest

ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM"

RUN apk add --no-cache libc6-compat gcompat

ADD /outfile/${TARGETPLATFORM}/FoodDecider-TG-Bot /FoodDecider-TG-Bot

WORKDIR /

RUN chmod +x FoodDecider-TG-Bot

CMD ["./FoodDecider-TG-Bot"]