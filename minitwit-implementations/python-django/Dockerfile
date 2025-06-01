FROM alpine:edge
WORKDIR /src

# Install dependencies
RUN apk add --no-cache python3 py3-pip postgresql-dev gcc musl-dev python3-dev

# Create and use virtual environment
RUN python3 -m venv /src/venv
ENV PATH="/src/venv/bin:$PATH"

# Upgrade pip and install dependencies
COPY requirements.txt requirements.txt
RUN /src/venv/bin/pip install --upgrade pip && /src/venv/bin/pip install --no-cache-dir -r requirements.txt

# Copy project files
COPY . /src

EXPOSE 5000

CMD ["/src/venv/bin/python", "src/manage.py", "runserver", "0.0.0.0:5000"]