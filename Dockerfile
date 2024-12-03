FROM alpine:3.20

ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM" && \
    apk add --no-cache libc6-compat gcompat

ADD /outfile/${TARGETPLATFORM}/FoodDecider-TG-Bot /FoodDecider-TG-Bot

WORKDIR /

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

RUN chmod +x FoodDecider-TG-Bot

CMD ["./FoodDecider-TG-Bot"]