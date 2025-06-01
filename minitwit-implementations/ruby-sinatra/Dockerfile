FROM alpine:edge

ENV APP_ENV=production
ENV PORT=5000

RUN apk add --no-cache ruby ruby-dev build-base postgresql-dev && gem install bundler

WORKDIR /waect-ruby-sinatra

COPY Gemfile ./

RUN bundle install

COPY ./src .

EXPOSE 5000

CMD ["bundle", "exec", "ruby", "main.rb", "-o", "0.0.0.0", "-p", "5000"]
