FROM debian:bullseye-slim

ENV APP_ENV=production
ENV PORT=5000

RUN apt-get update && apt-get install -y curl build-essential libpq-dev libssl-dev
RUN curl -L https://github.com/postmodern/ruby-install/archive/master.tar.gz | tar xz
RUN cd ruby-install-master && make install 
RUN ruby-install --system ruby 3.3.0

WORKDIR /waect-ruby-sinatra

COPY Gemfile ./

RUN gem install bundler && bundle install

COPY ./src .

EXPOSE 5000

CMD ["bundle", "exec", "ruby", "main.rb"]
