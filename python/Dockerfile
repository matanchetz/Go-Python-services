FROM python:3.8.0a2-alpine3.9

WORKDIR /usr/src/service/

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt


CMD [ "python", "./service.py" ]

