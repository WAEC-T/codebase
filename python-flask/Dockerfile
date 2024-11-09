FROM alpine:edge
# tesing dockerfile
WORKDIR /app

RUN apk add --no-cache python3

COPY requirements.txt requirements.txt

# create a vertial environment and install the requirements
RUN python3 -m venv venv
RUN source venv/bin/activate
RUN venv/bin/pip install --no-cache-dir -r requirements.txt

COPY . /app

EXPOSE 5000
ENV FLASK_APP=wsgi.py

CMD ["gunicorn", "--bind", "0.0.0.0:5000", "wsgi:app"]