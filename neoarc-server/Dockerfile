FROM python:3.12-alpine AS builder

WORKDIR /app
COPY pyproject.toml .
RUN pip install --no-cache-dir .

FROM python:3.12-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY --from=builder /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY . .

RUN mkdir -p static/uploads

EXPOSE 59248

ENV NEOARC_DEBUG=False
ENV NEOARC_HOST=0.0.0.0
ENV NEOARC_PORT=59248

CMD ["python", "wsgi.py"]
