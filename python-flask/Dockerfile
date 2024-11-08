FROM alpine:edge

# Install Python, pip, and PostgreSQL development tools
RUN apk add --no-cache python3 py3-pip postgresql-dev gcc musl-dev linux-headers

# Create a Python virtual environment
RUN python3 -m venv /opt/venv

# Set the virtual environment path
ENV PATH="/opt/venv/bin:$PATH"

# Set the working directory
WORKDIR /app

# Copy the application code and requirements
COPY . /app
COPY requirements.txt requirements.txt

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Expose the application port
EXPOSE 5000

# Set Flask environment variable
ENV FLASK_APP=wsgi.py

# Run the application using Gunicorn
CMD ["gunicorn", "--bind", "0.0.0.0:5000", "wsgi:app"]