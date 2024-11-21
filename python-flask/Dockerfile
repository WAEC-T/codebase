FROM alpine:edge
WORKDIR /app

# install dependencies
RUN apk add --no-cache python3 py3-pip postgresql-dev gcc musl-dev python3-dev # psycopg2-binar py3-pip gcc musl-dev

# run requirements
COPY requirements.txt requirements.txt
RUN python3 -m venv venv
RUN ./venv/bin/pip install --upgrade pip
RUN source venv/bin/activate
RUN venv/bin/pip install --no-cache-dir -r requirements.txt

COPY . /app

EXPOSE 5000
ENV FLASK_APP=wsgi.py

CMD ["venv/bin/gunicorn", "--bind", "0.0.0.0:5000", "wsgi:app"]